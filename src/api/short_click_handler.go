package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/isnandar1471/url_shortener/src/database"
	"github.com/isnandar1471/url_shortener/src/structs"
	"net/http"
	"os"
	"strings"
)

type ShortClick struct {
	Id        int    `json:"id"`
	ShortId   int    `json:"short_id"`
	ClickedAt int    `json:"clicked_at"`
	IpAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
}

func HandleGetShortClickByCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	shortCode := r.PathValue("short_code")

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
	fmt.Println("short_code", shortCode)
	fmt.Println("username", username)
	rows, _ := conn.Query(context.Background(), `
		SELECT short_clicks.id, 
		       short_clicks.short_id, 
		       short_clicks.clicked_at, 
		       short_clicks.ip_address, 
		       short_clicks.user_agent FROM short_clicks 
		           LEFT JOIN shorts ON shorts.id = short_clicks.short_id 
		           LEFT JOIN users ON shorts.user_id = users.id 
		                               WHERE shorts.code=$1 
		                                 AND users.username=$2
		`, shortCode, username)
	//WHERE short_clicks.id IN (SELECT id FROM shorts WHERE shorts.codo=$1)
	shortClicks := []ShortClick{}
	for rows.Next() {
		fmt.Println("ADA DATA")
		shortClick := ShortClick{}
		_ = rows.Scan(&shortClick.Id, &shortClick.ShortId, &shortClick.ClickedAt, &shortClick.IpAddress, &shortClick.UserAgent)
		fmt.Println(shortClick.Id, shortClick.ShortId, shortClick.ClickedAt, shortClick.IpAddress, shortClick.UserAgent)
		shortClicks = append(shortClicks, shortClick)
	}

	bytes, _ := json.Marshal(shortClicks)
	_, _ = w.Write(bytes)
}
