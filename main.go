package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-zoo/claw"
	mw "github.com/go-zoo/claw/middleware"
	"github.com/gorilla/mux"
	"github.com/namsral/flag"
)

var (
	env     = flag.String("env", "dev", "--env env name")
	url     = flag.String("url", "http://localhost", "--url domain name")
	intPort = flag.String("intPort", "8080", "--intPort 8080")
)

// Server contains the base URL for creating the short link urls
// the store is a layer on top of firestore to store links and the
// global counter.
type Server struct {
	URL   string
	Port  string
	Store *Store
}

// NewServer initialize a new server instance and create the store
func NewServer(env, url, port string) *Server {
	if env == "dev" {
		url = fmt.Sprintf("%s:%s", url, port)
	}
	s, err := NewStore()
	if err != nil {
		log.Fatal(err)
	}
	return &Server{
		URL:   url,
		Port:  port,
		Store: s,
	}
}

func main() {
	flag.Parse()

	r := mux.NewRouter()
	clw := claw.New(mw.NewLogger(os.Stdout, "[FORK]", 2))

	srvr := NewServer(*env, *url, *intPort)

	r.HandleFunc("/gen-link", http.HandlerFunc(srvr.GenerateLink)).Methods("POST", "OPTIONS")
	r.HandleFunc("/update-link", http.HandlerFunc(srvr.UpdateLink)).Methods("POST", "OPTIONS")
	r.HandleFunc("/info-link/{id}", http.HandlerFunc(srvr.InfoLink)).Methods("GET")
	r.HandleFunc("/{id}", http.HandlerFunc(srvr.FowardLinkHandler)).Methods("GET", "OPTIONS")

	r.HandleFunc("/", HomeHandler)
	if *env == "dev" {
		r.PathPrefix("/").Handler(http.HandlerFunc(StaticHandler))
	}

	log.Printf("Running on %s with URL %s:%s\n", *env, *url, *intPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *intPort), clw.Merge(allowCORS(r))))
}
