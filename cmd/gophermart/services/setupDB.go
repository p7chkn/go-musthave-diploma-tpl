package services

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
)

func MustSetupDatabase(ctx context.Context, db *sql.DB, log *zap.SugaredLogger) {
	createUsersSQL := `CREATE TABLE IF NOT EXISTS users (
							   id SERIAL PRIMARY KEY,
							   login VARCHAR(50) NOT NULL UNIQUE,
	                      	   password TEXT NOT NULL,
							   first_name VARCHAR(50),
							   last_name VARCHAR(50),
							   balance INT DEFAULT 0,
							   spend INT DEFAULT 0
						   );`
	_, err := db.ExecContext(ctx, createUsersSQL)
	if err != nil {
		log.Fatal(err)
	}
	//	log.Infof("Create table res: %v err: %v", res, err)
	//	createOrdersSQL := `CREATE TABLE IF NOT EXISTS orders (
	//   id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY ,
	//   user_id uuid REFERENCES users(id) ON DELETE CASCADE ,
	//   number VARCHAR(50) NOT NULL UNIQUE,
	//   status VARCHAR(50),
	//   uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	//   accrual INT DEFAULT 0
	//);`
	//	res, err = db.ExecContext(ctx, createOrdersSQL)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Infof("Create table res: %v err: %v", res, err)
	//	createWithdrawalsSQL := `CREATE TABLE IF NOT EXISTS withdrawals (
	//    id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
	//    user_id uuid REFERENCES users(id) ON DELETE CASCADE ,
	//    order_number VARCHAR (50) NOT NULL UNIQUE ,
	//    status VARCHAR(50) DEFAULT 'NEW',
	//    processed_at TIMESTAMP,
	//    sum INT
	//);`
	//	res, err = db.ExecContext(ctx, createWithdrawalsSQL)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	createJobStoreSQL := `CREATE TABLE IF NOT EXISTS jobstore (
	//     id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
	//     type VARCHAR(50),
	//     next_time_execute TIMESTAMP,
	//     parameters json,
	//     count INT,
	//     executed BOOL DEFAULT FALSE
	//);`
	//	res, err = db.ExecContext(ctx, createJobStoreSQL)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
}
