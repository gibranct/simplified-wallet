package db

import (
	"log"
	"sync"
	"time"

	"github.com.br/gibranct/simplified-wallet/internal/config"
	"github.com/jmoiron/sqlx"
)

var once sync.Once

func NewPostgresDB() *sqlx.DB {
	pgConfig := config.GetPostgresConfig()

	var db *sqlx.DB

	once.Do(func() {
		var err error
		db, err = sqlx.Connect("postgres", pgConfig.GetPostgresURL())
		if err != nil {
			log.Fatalln("Failed connect to db, err:", err)
		}
		db.SetMaxOpenConns(20)
		db.SetMaxIdleConns(5)
		db.SetConnMaxIdleTime(time.Minute * 30)
		err = db.Ping()
		if err != nil {
			log.Fatalln("Failed to ping db, err:", err)
		}
		log.Println("Connected to Postgres DB")
	})

	return db
}
