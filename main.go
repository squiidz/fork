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
	devEnv  = "dev"
	prod    = "cloud"
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
	if env == devEnv {
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

	r.HandleFunc("/", HomeHandler)
	r.Handle("/gen-link", allowCORS(http.HandlerFunc(srvr.GenerateLink))).Methods("POST", "OPTIONS")
	r.Handle("/update-link", allowCORS(http.HandlerFunc(srvr.UpdateLink))).Methods("POST", "OPTIONS")
	r.Handle("/info-link/{id}", allowCORS(http.HandlerFunc(srvr.InfoLink))).Methods("GET")
	r.Handle("/{id}", allowCORS(http.HandlerFunc(srvr.FowardLinkHandler))).Methods("GET", "OPTIONS")

	if *env == devEnv {
		r.PathPrefix("/").Handler(http.HandlerFunc(StaticHandler))
	}

	log.Printf("Running on %s with URL %s:%s\n", *env, *url, *intPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *intPort), clw.Merge(r)))
}
