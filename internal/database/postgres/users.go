package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/models"
)

func (db *PostgreDataBase) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	fmt.Println(user)
	sqlCreateUser := `INSERT INTO users (login, password, first_name, last_name) VALUES ($1, crypt($2, gen_salt('bf', 8)), $3, $4)`
	_, err := db.conn.ExecContext(ctx, sqlCreateUser, user.Login, user.Password, user.FirstName, user.LastName)
	if err != nil {
		return nil, err
	}
	resultUser, err := db.getUserByLogin(ctx, user.Login)
	if err != nil {
		return nil, err
	}
	return resultUser, err
}

func (db *PostgreDataBase) CheckPassword(ctx context.Context, user models.User) (models.User, error) {
	resultUser := models.User{}
	fmt.Println(user)

	var pass string
	q := db.conn.QueryRowContext(ctx, `SELECT password FROM users WHERE login = $1`, user.Login)
	err := q.Scan(&pass)
	fmt.Println(pass)

	q2 := db.conn.QueryRowContext(ctx, `SELECT crypt($2, password) FROM users WHERE login = $1`, user.Login, user.Password)

	err = q2.Scan(&pass)
	fmt.Println(pass)

	sqlCheckUserPassword := `SELECT id FROM users WHERE login = $1 AND password = crypt($2, password) FETCH FIRST ROW ONLY;`
	query := db.conn.QueryRowContext(ctx, sqlCheckUserPassword, user.Login, user.Password)
	err = query.Scan(&resultUser.ID)
	if err != nil {
		return resultUser, errors.New("wrong login or password")
	}
	if resultUser.ID == "" {
		return resultUser, errors.New("wrong login or password")
	}
	return resultUser, nil
}

func (db *PostgreDataBase) getUserByLogin(ctx context.Context, login string) (*models.User, error) {
	resultUser := &models.User{}
	sqlGetUser := `SELECT id, login, first_name, last_name, balance, spend FROM users WHERE login = $1`
	query := db.conn.QueryRowContext(ctx, sqlGetUser, login)
	err := query.Scan(&resultUser.ID, &resultUser.Login, &resultUser.FirstName, &resultUser.LastName,
		&resultUser.Balance, &resultUser.Spent)
	if err != nil {
		return resultUser, err
	}
	return resultUser, nil
}

func (db *PostgreDataBase) getUser(ctx context.Context, id string) (*models.User, error) {
	resultUser := &models.User{}
	sqlGetUser := `SELECT id, login, first_name, last_name, balance, spend FROM users WHERE id = $1`
	query := db.conn.QueryRowContext(ctx, sqlGetUser, id)
	err := query.Scan(&resultUser.ID, &resultUser.Login, &resultUser.FirstName, &resultUser.LastName,
		&resultUser.Balance, &resultUser.Spent)
	if err != nil {
		return resultUser, err
	}
	return resultUser, nil
}

func (db *PostgreDataBase) GetBalance(ctx context.Context, userID string) (models.UserBalance, error) {
	var result models.UserBalance
	sqlGetBalance := `SELECT balance, spend FROM users WHERE id = $1`
	query := db.conn.QueryRowContext(ctx, sqlGetBalance, userID)
	err := query.Scan(&result.Balance, &result.Spent)
	if err != nil {
		return result, err
	}
	return result, nil
}
