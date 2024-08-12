package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/isnandar1471/url_shortener/src/database"
	"github.com/isnandar1471/url_shortener/src/structs"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type RegisterBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginBody struct {
	Username string
	Password string
}

func HandlePostRegister(w http.ResponseWriter, r *http.Request) {
	var b RegisterBody
	bytes, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(bytes, &b)

	if b.Username == "" || b.Email == "" || b.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		bytes, _ := json.Marshal(structs.DefaultResponse{
			Message: "data is invalid",
		})
		_, _ = w.Write(bytes)
		return
	}

	conn := database.MakeConnection()

	//Check User Availability
	rows := conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE username=$1 AND email=$2", b.Username, b.Email)
	var totalSameUser int
	_ = rows.Scan(&totalSameUser)

	if totalSameUser != 0 {
		w.WriteHeader(http.StatusBadRequest)
		resBytes, _ := json.Marshal(structs.DefaultResponse{
			Message: "Those username and password has used by other user",
		})
		_, _ = w.Write(resBytes)
		return
	}

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(b.Password), bcrypt.DefaultCost)
	fmt.Print("encryptedPassword", string(encryptedPassword))
	createdAt := time.Now()
	_, err := conn.Exec(context.Background(), "INSERT INTO users(username, email, password_hash, created_at) VALUES ($1, $2, $3, $4)", b.Username, b.Email, string(encryptedPassword), createdAt.Unix())
	_ = conn.Close(context.Background())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bytes, _ := json.Marshal(structs.DefaultResponse{
			Message: err.Error(),
		})
		_, _ = w.Write(bytes)
		return
	}

	w.WriteHeader(http.StatusCreated)
	bytes, _ = json.Marshal(structs.DefaultResponse{
		Message: "Success",
	})
	_, _ = w.Write(bytes)
}

func HandleGetCheckUserExist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bytes, _ := io.ReadAll(r.Body)
	identifier := string(bytes)

	conn := database.MakeConnection()
	row := conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE username=$1 OR email=$1", identifier)
	_ = conn.Close(context.Background())

	var totalFoundUser int
	_ = row.Scan(&totalFoundUser)

	if totalFoundUser > 0 {
		w.WriteHeader(http.StatusBadRequest)
		bytes, _ := json.Marshal(structs.DefaultResponse{
			Message: "Those username and password has used by other user",
		})
		_, _ = w.Write([]byte(bytes))
		return
	}

	bytes, _ = json.Marshal(structs.DefaultResponse{
		Message: "Those username and password is avail",
	})
	_, _ = w.Write([]byte(bytes))
}

func HandlePostLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var B LoginBody
	b, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(b, &B)

	if B.Username == "" || B.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		bytes, _ := json.Marshal(structs.DefaultResponse{
			Message: "data is invalid",
		})
		_, _ = w.Write(bytes)
		return
	}

	conn := database.MakeConnection()
	defer conn.Close(context.Background())
	row := conn.QueryRow(context.Background(), "SELECT username, password_hash FROM users WHERE username=$1 OR email=$1 LIMIT 1", B.Username)

	var username, hashedPassword string

	_ = row.Scan(&username, &hashedPassword)
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(B.Password))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		bytes, _ := json.Marshal(structs.DefaultResponse{
			Message: "password not match: " + err.Error(),
		})
		_, _ = w.Write(bytes)
		return
	}

	payloadMap := jwt.RegisteredClaims{
		Subject:   strings.Trim(username, " "),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payloadMap)

	token, _ := jwtToken.SignedString([]byte(os.Getenv("JWT_KEY")))
	bytes, _ := json.Marshal(structs.DefaultResponse{
		Message: token,
	})
	_, _ = w.Write(bytes)
}
