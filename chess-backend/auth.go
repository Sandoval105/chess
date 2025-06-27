package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "erro ao ler dados ", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, "erro ao gerar hash", http.StatusInternalServerError)
		return
	}

	_, err = DB.Exec("INSERT INTO users (nickname,password_hash) VALUES (?,?)", u.Nickname, hash)

	if err != nil {
		http.Error(w, "erro ao enviar ao banco", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("cadastro concluido com sucesso!"))

}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	var u User

	err := json.NewDecoder(r.Body).Decode(&u)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, "erro ao ler dados", http.StatusBadRequest)
		return
	}

	var hash string
	err = DB.QueryRow("SELECT password_hash FROM users WHERE nickname = ?", u.Nickname).Scan(&hash)

	if err != nil {
		http.Error(w, "erro ao buscar no banco", http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(u.Password)) != nil {
		http.Error(w, "senha incorreta", http.StatusUnauthorized)
		return
	}

	w.Write([]byte("login com sucesso! "))

}
