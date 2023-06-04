package datastore

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func CreateNewSearchRDB() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return rdb
}

func CreateNewSelectionRDB() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})
	return rdb
}

func CreateNewDBForSurvey() *sql.DB {
	connStr := "user=" + viper.GetString("postgre.user") +
		" password=" + viper.GetString("postgre.password") +
		" dbname=" + viper.GetString("postgre.dbname") +
		" sslmode=" + viper.GetString("postgre.sslmode")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}

func CreateNewDBForPreferences() *sql.DB {
	connStr := "user=" + viper.GetString("postgre.user") +
		" password=" + viper.GetString("postgre.password") +
		" dbname=" + viper.GetString("postgre.dbname1") +
		" sslmode=" + viper.GetString("postgre.sslmode")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}

func CreateNewDBForVehicles() *sql.DB {
	connStr := "user=" + viper.GetString("postgre.user") +
		" password=" + viper.GetString("postgre.password") +
		" dbname=" + viper.GetString("postgre.dbname2") +
		" sslmode=" + viper.GetString("postgre.sslmode")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}
