package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("harryp_istiv")

func verificarToken(r string) (string, error) {
	tokenHeader := r

	if tokenHeader == "" {
		return "", errors.New("token ausente")
	}

	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("algoritmo inválido")
		}
		return jwtSecret, nil

	})

	if err != nil || !token.Valid {
		return "", errors.New("token invalido")

	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("erro nas claims")
	}

	nickname, ok := claims["nickname"].(string)
	if !ok {
		return "", errors.New("nickname ausente")
	}

	return nickname, nil

}

func gerarToken(nickname string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nickname": nickname,
		"exp":      time.Now().Add((time.Hour * 24)).Unix(),
	})

	return token.SignedString(jwtSecret)

}

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

	// Aqui geramos um "token" simples
	token, err := gerarToken(u.Nickname)
	if err != nil {
		http.Error(w, "erro ao gerar token ", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token":    token,
		"nickname": u.Nickname,
	})

}
