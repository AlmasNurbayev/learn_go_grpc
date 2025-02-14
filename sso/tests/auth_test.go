package tests

import (
	"fmt"
	"sso/tests/suite"
	"testing"
	"time"

	ssov1 "github.com/AlmasNurbayev/learn_go_grpc_protos/generated/sso"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppId = 0
	AppId      = 100
	appSecret  = "test_secret"

	passDefaultLen = 8
)

var loginTime time.Time

func TestAuth_RegistrationLogin_Positive(t *testing.T) {
	ctx, st := suite.New(t)
	var respLogin *ssov1.LoginResponse
	var respRegistration *ssov1.RegisterResponse

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, false, false, passDefaultLen)

	countryCode := gofakeit.Number(1, 9)                  // Код страны (например, США = 1)
	areaCode := gofakeit.Number(100, 999)                 // Код региона
	subscriberNumber := gofakeit.Number(1000000, 9999999) // 7-значный номер
	phone := fmt.Sprintf("+%d%d%d", countryCode, areaCode, subscriberNumber)

	t.Run("Registration email", func(t *testing.T) {
		var err error
		respRegistration, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
			Email:    email,
			Password: password,
		})

		require.NoError(t, err)
		assert.NotEmpty(t, respRegistration.GetUserId())
	})

	t.Run("Registration phone", func(t *testing.T) {
		respRegistrationPhone, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
			Phone:    phone,
			Password: password,
		})
		// t.Log(respRegistrationPhone)
		// t.Logf("Phone: %s", phone)
		// t.Logf("%s", err.Error())
		require.NoError(t, err)

		assert.NotEmpty(t, respRegistrationPhone.GetUserId())
	})

	t.Run("Login", func(t *testing.T) {
		loginTime = time.Now()
		var err error
		respLogin, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
			Login:    email,
			Password: password,
			Type:     "email",
			AppId:    AppId,
		})

		require.NoError(t, err)
		assert.NotEmpty(t, respLogin.GetToken())
	})

	t.Run("Token parse", func(t *testing.T) {
		tokenParsed, err := jwt.Parse(respLogin.GetToken(), func(token *jwt.Token) (interface{}, error) {
			return []byte(appSecret), nil
		})
		require.NoError(t, err)

		claims, ok := tokenParsed.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, respRegistration.GetUserId(), int64(claims["id"].(float64)))
		assert.Equal(t, email, claims["email"].(string))
		assert.Equal(t, AppId, int(claims["app_id"].(float64)))
		assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), int64(claims["exp"].(float64)), 2)
	})
}

func TestAuth_RegistrationLogin_Negative(t *testing.T) {
	ctx, st := suite.New(t)
	var respRegistration *ssov1.RegisterResponse

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, false, false, passDefaultLen)

	t.Run("Registration", func(t *testing.T) {
		var err error
		respRegistration, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
			Email:    email,
			Password: password,
		})

		require.NoError(t, err)
		assert.NotEmpty(t, respRegistration.GetUserId())
	})

	t.Run("Registration empty email and phone", func(t *testing.T) {
		var err error
		respRegistration, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
			Phone:    "",
			Email:    "",
			Password: password,
		})
		require.Error(t, err)
	})

	t.Run("Registration duplicate", func(t *testing.T) {
		var err error
		respRegistration, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
			Email:    email,
			Password: password,
		})

		require.Error(t, err)
		assert.ErrorContains(t, err, "already exists")
	})

	t.Run("Login password wrong", func(t *testing.T) {
		loginTime = time.Now()
		_, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
			Login:    email,
			Password: password + " ",
			Type:     "email",
			AppId:    AppId,
		})
		require.Error(t, err)
	})

	t.Run("Login password empty", func(t *testing.T) {
		_, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
			Login:    email,
			Password: "",
			Type:     "email",
			AppId:    AppId,
		})
		require.Error(t, err)
	})

	t.Run("Login wrong", func(t *testing.T) {
		_, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
			Login:    "test",
			Password: password,
			Type:     "email",
			AppId:    AppId,
		})
		require.Error(t, err)
	})

	t.Run("Login empty", func(t *testing.T) {
		_, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
			Login:    "",
			Password: password,
			Type:     "email",
			AppId:    AppId,
		})
		require.Error(t, err)
	})

	t.Run("AppId wrong", func(t *testing.T) {
		_, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
			Login:    email,
			Password: password,
			Type:     "email",
			AppId:    gofakeit.GlobalFaker.Int32(),
		})
		require.Error(t, err)
	})

}
