package postgresql

import (
	"biletter/internal/domain/entity"
	"context"
)

func (s *storage) GetUserByEmail(ctx context.Context, email string) (data entity.AuthUser, err error) {
	query := `
		SELECT 
			user_id, email, is_active, first_name, surname, last_logged_in, password_hash, password_plain
		FROM "users"
		WHERE email = $1 AND is_active = true;
	`
	err = s.client.QueryRow(ctx, query, email).Scan(
		&data.ID,
		&data.Email,
		&data.IsActive,
		&data.FirstName,
		&data.Surname,
		&data.LastLoggedIn,
		&data.PasswordHash,
		&data.PasswordPlain,
	)

	if err != nil {
		s.logger.Errorf("Пользователя с email=%s не существует: ERROR=%v", email, err)
		return data, err
	}

	return data, nil
}
