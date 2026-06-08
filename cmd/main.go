package main

import (
	"CloudHub/cmd/api"
	"CloudHub/config"
	db2 "CloudHub/db"
	"CloudHub/internal/queue"
	"database/sql"
	"log"
)

func main() {
	db, _ := db2.NewPostgreSqlStorage(config.Envs.ConnString)

	rdb := queue.NewRedisClient()

	app := api.NewApplication(config.Envs.Port, db, rdb)

	initStorage(db)

	err := app.Start(app.Mount())
	if err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("DB connected Successfully ...!!")
}
