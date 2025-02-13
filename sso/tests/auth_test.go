package tests

import (
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
var respLogin *ssov1.LoginResponse
var respRegistration *ssov1.RegisterResponse

func TestAuth_Login_Positive(t *testing.T) {
	ctx, st := suite.New(t)

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
