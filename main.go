package main

import (
	"flag"
	"github.com/isnandar1471/url_shortener/src/api"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

// Initialize in init function
var Port *int
var Host *string

func init() {
	//	Load environment variable
	_ = godotenv.Load()

	envPort, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		envPort = 5555
	}

	envHost := os.Getenv("HOST")
	if envHost == "" {
		envHost = "0.0.0.0"
	}

	envJwtKey := os.Getenv("JWT_KEY")
	if envJwtKey == "" {
		envJwtKey = randSeq(10)
	}

	Port = flag.Int("port", envPort, "The port that used to run the app")
	Host = flag.String("host", envHost, "The host that used to run the app")
	JwtKey := flag.String("jwt-key", envJwtKey, "The JWT Key that used to run the app")
	flag.Parse()

	_ = os.Setenv("JWT_KEY", *JwtKey)
}

func main() {
	mux := http.NewServeMux()

	//region Add api handler in this section
	mux.HandleFunc("POST /api/register", enableCors(api.HandlePostRegister))
	mux.HandleFunc("POST /api/login", enableCors(api.HandlePostLogin))
	mux.HandleFunc("POST /api/check-user", enableCors(api.HandleGetCheckUserExist))
	mux.HandleFunc("GET /api/shorts", enableCors(api.HandleGetShorts))
	mux.HandleFunc("POST /api/short", enableCors(api.HandlePostShort))
	mux.HandleFunc("PATCH /api/short/{short_code}", enableCors(api.HandlePatchShortByCode))
	mux.HandleFunc("DELETE /api/short/{short_code}", enableCors(api.HandleDeleteShortByCode))
	mux.HandleFunc("GET /api/short_clicks/{short_code}", enableCors(api.HandleGetShortClickByCode))
	mux.HandleFunc("GET /{short_code}", enableCors(api.HandleGetGo))

	// Sepertinya ini untuk mengatasi cors secara global. Perlu mencoba untuk yang secara spesifik
	mux.HandleFunc("OPTIONS /", enableCors(func(w http.ResponseWriter, r *http.Request) {}))

	//endregion

	address := *Host + ":" + strconv.Itoa(*Port)
	println("Server starting...", address)

	server := http.Server{
		Addr:    address,
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func enableCors(handler http.HandlerFunc) http.HandlerFunc {
	println("enableCors otw")
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		println("enableCors jalan")
		handler(w, r)
	}
}
