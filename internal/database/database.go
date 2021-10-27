package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/models"
	"log"

	"github.com/p7chkn/go-musthave-diploma-tpl/internal/app/handlers"
)

type PostgreDataBase struct {
	conn *sql.DB
}

func NewDatabaseRepository(db *sql.DB) handlers.RepositoryInterface {
	return handlers.RepositoryInterface(NewDatabase(db))
}

func NewDatabase(db *sql.DB) *PostgreDataBase {
	result := &PostgreDataBase{
		conn: db,
	}
	return result
}

func (db *PostgreDataBase) Ping(ctx context.Context) error {

	err := db.conn.PingContext(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (db *PostgreDataBase) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	sqlCreateUser := `INSERT INTO users (login, password, first_name, last_name) VALUES ($1, crypt($2, gen_salt('bf', 8)), $3, $4)`
	_, err := db.conn.ExecContext(ctx, sqlCreateUser, user.Login, user.Password, user.FirstName, user.LastName)
	if err != nil {
		return nil, err
	}
	resultUser, err := db.getUser(ctx, user.Login)
	if err != nil {
		return nil, err
	}
	return resultUser, err
}

func (db *PostgreDataBase) CheckPassword(ctx context.Context, user models.User) (models.User, error) {
	resultUser := models.User{}
	sqlCheckUserPassword := `SELECT id FROM users WHERE login = lower($1) AND password = crypt($2, password) FETCH FIRST ROW ONLY;`
	query := db.conn.QueryRowContext(ctx, sqlCheckUserPassword, user.Login, user.Password)
	err := query.Scan(&resultUser.Id)
	if err != nil {
		return resultUser, err
	}
	if resultUser.Id == "" {
		return resultUser, errors.New("wrong login or password")
	}
	return resultUser, nil
}

func (db *PostgreDataBase) getUser(ctx context.Context, login string) (*models.User, error) {
	resultUser := &models.User{}
	sqlGetUser := `SELECT id, login, first_name, last_name, balance, spent FROM users WHERE login = $1`
	query := db.conn.QueryRowContext(ctx, sqlGetUser, login)
	err := query.Scan(&resultUser.Id, &resultUser.Login, &resultUser.FirstName, &resultUser.LastName,
		&resultUser.Balance, &resultUser.Spent)
	if err != nil {
		return resultUser, err
	}
	return resultUser, nil
}
