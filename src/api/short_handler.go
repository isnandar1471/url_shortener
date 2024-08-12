package api

import (
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/isnandar1471/url_shortener/src/database"
	"github.com/isnandar1471/url_shortener/src/structs"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Short struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Code           string `json:"code"`
	DestinationUrl string `json:"destination_url"`
	UserId         int    `json:"user_id"`
	CreatedAt      int    `json:"created_at"`
	ClickCount     int    `json:"click_count"`
}

func HandleGetShorts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authorization, isAuthProvide := r.Header["Authorization"]
	if !isAuthProvide {
		w.WriteHeader(http.StatusForbidden)
		bytes, _ := json.Marshal(structs.DefaultResponse{
			Message: "Unauthorized",
		})
		_, _ = w.Write(bytes)
		return
	}
	token := strings.Split(authorization[0], " ")[1]

	var jwtClaim jwt.RegisteredClaims
	jwtToken, _ := jwt.ParseWithClaims(token, &jwtClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if !jwtToken.Valid {
		w.WriteHeader(http.StatusForbidden)
		bytes, _ := json.Marshal(structs.DefaultResponse{
			Message: "Invalid authorization",
		})
		_, _ = w.Write(bytes)
		return
	}

	username := jwtClaim.Subject

	conn := database.MakeConnection()
	defer conn.Close(context.Background())

	rows, _ := conn.Query(context.Background(), "SELECT shorts.id, shorts.name, shorts.code, shorts.destination_url, shorts.user_id, shorts.created_at, shorts.click_count FROM shorts JOIN users on users.id = shorts.user_id WHERE users.username=$1", username)

	shorts := []Short{}
	for rows.Next() {
		short := Short{}
		_ = rows.Scan(&short.Id, &short.Name, &short.Code, &short.DestinationUrl, &short.UserId, &short.CreatedAt, &short.ClickCount)

		shorts = append(shorts, short)
	}

	bytes, _ := json.Marshal(shorts)
	_, _ = w.Write(bytes)
}

func HandlePostShort(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authorization, isAuthProvided := r.Header["Authorization"]
	if !isAuthProvided {
		w.WriteHeader(http.StatusForbidden)
		bytes, _ := json.Marshal(structs.DefaultResponse{
			Message: "Unauthorized",
		})
		_, _ = w.Write(bytes)
		return
	}
	token := strings.Split(authorization[0], " ")[1]

	var jwtClaim jwt.RegisteredClaims
	jwtToken, _ := jwt.ParseWithClaims(token, &jwtClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if !jwtToken.Valid {
		w.WriteHeader(http.StatusForbidden)
		bytes, _ := json.Marshal(structs.DefaultResponse{
			Message: "Invalid authorization",
		})
		_, _ = w.Write(bytes)
		return
	}

	username := jwtClaim.Subject

	short := Short{}
	bytes, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(bytes, &short)

	if short.Name == "" || short.Code == "" || short.DestinationUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		r, _ := json.Marshal(structs.DefaultResponse{
			Message: "data is invalid",
		})
		_, _ = w.Write(r)
		return
	}

	conn := database.MakeConnection()
	defer conn.Close(context.Background())
	trx, err1 := conn.Begin(context.Background())
	_, err2 := trx.Exec(context.Background(), "INSERT INTO shorts (name, code, destination_url, user_id, created_at) VALUES ($1, $2, $3, (SELECT id FROM users WHERE username=$4 LIMIT 1), $5)", short.Name, short.Code, short.DestinationUrl, username, time.Now().Unix())
	_, err3 := trx.Exec(context.Background(), "UPDATE users SET shorts_count = shorts_count + 1 WHERE username=$1", username)

	if err1 != nil || err2 != nil || err3 != nil {
		println(err1)
		println(err2)
		println(err3)
		trx.Rollback(context.Background())
	}
	trx.Commit(context.Background())

	w.WriteHeader(http.StatusCreated)
	bytes, _ = json.Marshal(structs.DefaultResponse{
		Message: "Success",
	})
	_, _ = w.Write(bytes)
}

func HandlePatchShortByCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	shortCode := r.PathValue("short_code")

	authorization, isAuthProvided := r.Header["Authorization"]
	if !isAuthProvided {
		w.WriteHeader(http.StatusForbidden)
		bytes, _ := json.Marshal(structs.DefaultResponse{
			Message: "Unauthorized",
		})
		_, _ = w.Write(bytes)
		return
	}
	token := strings.Split(authorization[0], " ")[1]

	var jwtClaim jwt.RegisteredClaims
	jwtToken, _ := jwt.ParseWithClaims(token, &jwtClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if !jwtToken.Valid {
		w.WriteHeader(http.StatusForbidden)
		bytes, _ := json.Marshal(structs.DefaultResponse{
			Message: "Invalid authorization",
		})
		_, _ = w.Write(bytes)
		return
	}

	username := jwtClaim.Subject

	short := Short{}
	bytes, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(bytes, &short)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		bytes, _ = json.Marshal(structs.DefaultResponse{
			Message: "FAILED: " + err.Error(),
		})
		_, _ = w.Write(bytes)
		return
	}

	query := ""
	paramNum := 1
	params := []any{}
	if short.Name != "" {
		query += "name=$" + strconv.Itoa(paramNum)
		paramNum++
		params = append(params, short.Name)
	}
	if short.Code != "" {
		if query != "" {
			query += ", "
		}
		query += "code=$" + strconv.Itoa(paramNum)
		paramNum++
		params = append(params, short.Code)
	}
	if short.DestinationUrl != "" {
		if query != "" {
			query += ", "
		}
		query += "destination_url=$" + strconv.Itoa(paramNum)
		paramNum++
		params = append(params, short.DestinationUrl)
	}
	if query != "" {
		query = "UPDATE shorts SET " + query + " WHERE code=$" + strconv.Itoa(paramNum)
		paramNum++
		query = query + " AND user_id=(SELECT id FROM users WHERE username=$" + strconv.Itoa(paramNum) + " LIMIT 1)"

		params = append(params, shortCode, username)

		conn := database.MakeConnection()
		defer conn.Close(context.Background())
		_, _ = conn.Exec(context.Background(), query, params...)
	}

	bytes, _ = json.Marshal(structs.DefaultResponse{
		Message: "Success",
	})
	_, _ = w.Write(bytes)
}

func HandleDeleteShortByCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	shortCode := r.PathValue("short_code")

	authorization, isAuthProvided := r.Header["Authorization"]
	if !isAuthProvided {
		w.WriteHeader(http.StatusForbidden)
		bytes, _ := json.Marshal(structs.DefaultResponse{
			Message: "Unauthorized",
		})
		_, _ = w.Write(bytes)
		return
	}
	token := strings.Split(authorization[0], " ")[1]

	var jwtClaim jwt.RegisteredClaims
	jwtToken, _ := jwt.ParseWithClaims(token, &jwtClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if !jwtToken.Valid {
		w.WriteHeader(http.StatusForbidden)
		bytes, _ := json.Marshal(structs.DefaultResponse{
			Message: "Invalid authorization",
		})
		_, _ = w.Write(bytes)
		return
	}

	username := jwtClaim.Subject

	conn := database.MakeConnection()
	defer conn.Close(context.Background())
	trx, err1 := conn.Begin(context.Background())
	_, err2 := trx.Exec(context.Background(), "DELETE FROM shorts WHERE code=$1 AND user_id=(SELECT id FROM users WHERE username=$2 LIMIT 1);", shortCode, username)
	_, err3 := trx.Exec(context.Background(), "UPDATE users SET shorts_count = shorts_count - 1 WHERE username=$1;", username)

	if err1 != nil || err2 != nil || err3 != nil {
		println(err1)
		println(err2)
		println(err3)
		trx.Rollback(context.Background())
	}
	trx.Commit(context.Background())

	bytes, _ := json.Marshal(structs.DefaultResponse{
		Message: "Success",
	})
	_, _ = w.Write(bytes)
}
