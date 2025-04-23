package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"xlink/user_service/internal/utils"
)

type UserStorageRepositoryPostgres struct {
	pool        *pgxpool.Pool
	tokenLength int8
}

func NewUserStorageRepositoryPostgres(pool *pgxpool.Pool, tokenLength int8) *UserStorageRepositoryPostgres {
	return &UserStorageRepositoryPostgres{
		pool:        pool,
		tokenLength: tokenLength,
	}
}

func (u *UserStorageRepositoryPostgres) CheckToken(userId string, token string) (bool, error) {
	// get count of users with given id & token, that should be equal to 0 or 1
	query := `SELECT count(*) FROM user_service.users u WHERE u.id=$1 AND u.token=$2`

	var usersCount int
	err := u.pool.QueryRow(context.Background(), query, userId, token).Scan(&usersCount)
	if err != nil {
		return false, fmt.Errorf(
			"(postgres) error while checking token in postgres: user_id=%s, token=%s, err=%v",
			userId, token, err,
		)
	}
	return usersCount > 0, nil
}

func (u *UserStorageRepositoryPostgres) RefreshToken(userId string, token string) (string, error) {
	query := `UPDATE user_service.users SET token=$1 WHERE id=$2 AND token=$3`

	newToken := utils.GenerateToken(u.tokenLength)
	commandTag, err := u.pool.Exec(context.Background(), query, newToken, userId, token)
	if err != nil {
		return "", fmt.Errorf("(postgres) couldn't refresh token in postgres: %v", err)
	}
	if commandTag.RowsAffected() == 0 {
		return "", fmt.Errorf("(postgres) userId or/and token is invalid")
	}
	return newToken, nil
}

func (u *UserStorageRepositoryPostgres) CreateUser(telegramId *int64, isStaff *bool, isAdmin *bool) (string, string, error) {
	query := `INSERT INTO user_service.users (telegram_id, is_staff, is_admin, token)
				VALUES ($1, $2, $3, $4)
				RETURNING id`

	var id string
	token := utils.GenerateToken(u.tokenLength)
	err := u.pool.QueryRow(context.Background(), query, telegramId, isStaff, isAdmin, token).Scan(&id)
	if err != nil {
		return "", "", fmt.Errorf("(postgres) couldn't create user (tgId=%d, isStaff=%b, isAdmin=%b): %v",
			telegramId, isStaff, isAdmin, err)
	}
	return id, token, nil
}

func (u *UserStorageRepositoryPostgres) GetUser(userId string) (string, string, *int64, error) {
	query := `SELECT telegramId, isStaff, isAdmin FROM user_service.users WHERE id = $1`

	var telegramId sql.NullInt64
	var isStaff, isAdmin bool
	var role string

	err := u.pool.QueryRow(context.Background(), query, userId).Scan(&telegramId, &isStaff, &isAdmin)
	if err != nil {
		return "", "", nil, fmt.Errorf("(postgres) couldn't get user by id=%s, err=%v", userId, err)
	}

	role = utils.GetRoleByIsStaffIsAdmin(isStaff, isAdmin)

	return userId, role, &telegramId.Int64, nil
}

func (u *UserStorageRepositoryPostgres) GetUserIDByToken(token string) (string, bool, error) {
	query := `SELECT id, count(*) FROM user_service.users WHERE token = $1`

	var userId string

	err := u.pool.QueryRow(context.Background(), query, token).Scan(&userId)
	if err != nil {
		var formattedError error
		if errors.Is(err, pgx.ErrNoRows) {
			formattedError = errors.New("(postgres) token doesn't exist")
		} else {
			formattedError = fmt.Errorf("(postgres) couldn't get userId: %v", err)
		}

		return "", false, formattedError
	}

	return userId, true, nil
}

func (u *UserStorageRepositoryPostgres) GetUserIDByTgId(tgId int64) (string, bool, error) {
	query := `SELECT id, count(*) FROM user_service.users WHERE telegram_id = $1`

	var userId string

	err := u.pool.QueryRow(context.Background(), query, tgId).Scan(&userId)
	if err != nil {
		var formattedError error
		if errors.Is(err, pgx.ErrNoRows) {
			formattedError = fmt.Errorf("(postgres) telegram_id='%d' doesn't exist", tgId)
		} else {
			formattedError = fmt.Errorf("(postgres) couldn't get userId by tgId='%d': %v", tgId, err)
		}

		return "", false, formattedError
	}

	return userId, true, nil
}

func (u *UserStorageRepositoryPostgres) UpdateUser(userId string, telegramId *int64, isStaff *bool, isAdmin *bool) (bool, error) {

	// if we don't check whether user exists,
	// we won't know if rowsAffected=0 means 'doesn't exist' or 'wasn't updated'
	queryUserExists := `SELECT id FROM user_service.users WHERE id = $1`
	commandTag, err := u.pool.Exec(context.Background(), queryUserExists, userId)
	if err != nil {
		return false, fmt.Errorf("(postgres) couldn't check if user exists: %v", err)
	}
	if commandTag.RowsAffected() == 0 {
		return false, fmt.Errorf("(postgres) userId doesn't exist")
	}

	// query itself
	query := `UPDATE user_service.users WHERE id = $1`
	args := []any{userId}

	if telegramId != nil {
		args = append(args, telegramId)
		query += fmt.Sprintf(` SET telegram_id = $%d`, len(args))
	}

	if isStaff != nil {
		args = append(args, isStaff)
		query += fmt.Sprintf(` SET is_staff = $%d`, len(args))
	}

	if isAdmin != nil {
		args = append(args, isAdmin)
		query += fmt.Sprintf(` SET is_admin = $%d`, len(args))
	}

	commandTag, err = u.pool.Exec(context.Background(), query, args...)
	if err != nil {
		return false, fmt.Errorf("(postgres) couldn't update user in postgres: %v", err)
	}
	if commandTag.RowsAffected() == 0 {
		return false, fmt.Errorf("(postgres) no changes applied")
	}

	return true, nil
}

func (u *UserStorageRepositoryPostgres) DeleteUser(userId string) (bool, error) {
	query := `DELETE FROM user_service.users WHERE id = $1`
	commandTag, err := u.pool.Exec(context.Background(), query, userId)
	if err != nil {
		return false, fmt.Errorf("(postgres) couldn't delete user in postgres: %v", err)
	}
	if commandTag.RowsAffected() == 0 {
		return false, fmt.Errorf("(postgres) user_id='%s' doesn't exist", userId)
	}

	return true, nil
}

func (u *UserStorageRepositoryPostgres) GetRole(userId string) (string, bool, bool, error) {
	query := `SELECT is_staff, is_admin FROM user_service.users WHERE id = $1`

	var isStaff, isAdmin bool
	var role string

	err := u.pool.QueryRow(context.Background(), query, userId).Scan(&isStaff, &isAdmin)
	if err != nil {
		return "", false, false, fmt.Errorf("(postgres) couldn't get user role by id=%s, err=%v", userId, err)
	}

	role = utils.GetRoleByIsStaffIsAdmin(isStaff, isAdmin)

	return role, isStaff, isAdmin, nil
}
