package main

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Initialize in init function
var DocsHost *string
var DocsPort *int

func init() {
	//	Load environment variable
	_ = godotenv.Load()

	envPort, err := strconv.Atoi(os.Getenv("DOCS_PORT"))
	if err != nil {
		envPort = 55555
	}

	envHost := os.Getenv("DOCS_HOST")
	if envHost == "" {
		envHost = "0.0.0.0"
	}

	DocsPort = flag.Int("docs port", envPort, "The port that used to run the docs")
	DocsHost = flag.String("docs host", envHost, "The host that used to run the docs")
	flag.Parse()

}

func main() {
	// Create Server
	mux := http.NewServeMux()

	// Serve public directory
	fileServerHandler := http.FileServer(http.Dir("."))
	mux.Handle("/", http.StripPrefix("/", fileServerHandler))

	address := *DocsHost + ":" + strconv.Itoa(*DocsPort)
	println("Server starting...", address)

	server := http.Server{
		Addr:    address,
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
