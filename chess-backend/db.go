package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("mysql", "root:harryp_istiv@tcp(127.0.0.1:3306)/chess_game")

	if err != nil {
		log.Fatal("erro ao abrir conexão com o banco ", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("erro ao se concetar com o banco ", err)
	}
	log.Println("banco de dados conectado com sucesso! ")
}
