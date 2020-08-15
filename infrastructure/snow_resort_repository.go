package repository

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kotaroooo0/snowforecast-twitter-bot/domain"
)

type SnowResortRepositoryImpl struct {
	DB *sqlx.DB
}

func NewSnowResortRepositoryImpl(db *sqlx.DB) domain.SnowResortRepository {
	return &SnowResortRepositoryImpl{
		DB: db,
	}
}

func NewDBClient(dbConfig *DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbConfig.User, dbConfig.Password, dbConfig.Addr, dbConfig.Port, dbConfig.DB),
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewDBConfig(user, password, addr, port, db string) *DBConfig {
	return &DBConfig{
		User:     user,
		Password: password,
		Addr:     addr,
		Port:     port,
		DB:       db,
	}
}

type DBConfig struct {
	User     string
	Password string
	Addr     string
	Port     string
	DB       string
}

func (s SnowResortRepositoryImpl) FindAll() ([]*domain.SnowResort, error) {
	rows, err := s.DB.Queryx("select * from snow_resorts")
	if err != nil {
		return []*domain.SnowResort{}, err
	}

	var snowResorts []*domain.SnowResort
	for rows.Next() {
		var snowResort domain.SnowResort
		err := rows.StructScan(&snowResort)
		if err != nil {
			return []*domain.SnowResort{}, err
		}
		snowResorts = append(snowResorts, &snowResort)
	}
	return snowResorts, nil
}
