package pgsql

import (
	"expence_management/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PqsqlRepo struct {
	conn *gorm.DB
}

func NewPgsqlRepo(dsn string) (*PqsqlRepo, error) {
	conn, err := connect(dsn)
	if err != nil {
		return nil, err
	}

	conn.AutoMigrate(
		domain.User{},
		domain.Income{},
		domain.Expense{},
		domain.Category{},
		domain.UpcomingExpense{},
	)

	return &PqsqlRepo{
		conn: conn,
	}, nil
}

// dsn "postgres://username:password@localhost:5432/mydb"
func connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, err
}
