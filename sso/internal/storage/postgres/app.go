package postgres

import (
	"context"
	"database/sql"
	"errors"
	"sso/internal/errorsPackage"
	"sso/internal/models"
)

func (s *Storage) GetAppById(ctx context.Context, id int) (app models.App, err error) {
	query := `SELECT id, name, secret FROM apps WHERE id = $1`
	err = s.db.QueryRowContext(ctx, query, id).Scan(&app.Id, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return app, errorsPackage.ErrAppNotFound
		}
		return app, err
	}
	return app, nil
}
