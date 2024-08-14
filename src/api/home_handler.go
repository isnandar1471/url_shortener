package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/isnandar1471/url_shortener/src/database"
	"github.com/isnandar1471/url_shortener/src/structs"
	"net/http"
	"time"
)

func HandleGetHome(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("halo home"))
}

func HandleGetGo(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("short_code")

	ip := r.RemoteAddr
	if forwardedFor, forwardedForExist := r.Header["X-Forwarded-For"]; forwardedForExist {
		ip = forwardedFor[0]
	}
	if realIP, realIPExist := r.Header["X-Real-IP"]; realIPExist {
		ip = realIP[0]
	}

	conn := database.MakeConnection()
	defer conn.Close(context.Background())

	row := conn.QueryRow(context.Background(), `SELECT id, destination_url FROM shorts WHERE code=$1`, shortCode)
	var shortId int
	var destinationUrl string
	err := row.Scan(&shortId, &destinationUrl)

	if err != nil {
		w.WriteHeader(http.StatusGone)
		bytes, _ := json.Marshal(structs.DefaultResponse{
			Message: "Page that you are looking for doesnt exist:" + err.Error(),
		})
		_, _ = w.Write(bytes)
		return
	}

	trx, err1 := conn.Begin(context.Background())
	_, err2 := trx.Exec(context.Background(), "UPDATE shorts SET click_count = click_count + 1 WHERE code=$1", shortCode)
	_, err3 := trx.Exec(context.Background(), "INSERT INTO short_clicks (id, short_id, clicked_at, ip_address, user_agent) VALUES ((SELECT MAX(id) FROM short_clicks) + 1, $1, $2, $3, $4)", shortId, time.Now().Unix(), ip, r.UserAgent())

	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Println(err1)
		fmt.Println(err2)
		fmt.Println(err3)
		trx.Rollback(context.Background())
	}
	trx.Commit(context.Background())

	http.Redirect(w, r, destinationUrl, http.StatusSeeOther)

}
