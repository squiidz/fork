package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func HomeHandler(rw http.ResponseWriter, req *http.Request) {
	if *env == "dev" {
		http.ServeFile(rw, req, "fork-ui/dist/index.html")
		return
	}
	http.Redirect(rw, req, "www.fork.pw", http.StatusMovedPermanently)
}

func StaticHandler(rw http.ResponseWriter, req *http.Request) {
	http.StripPrefix("/", http.FileServer(http.Dir("fork-ui/dist"))).ServeHTTP(rw, req)
}

func (s *Server) FowardLinkHandler(rw http.ResponseWriter, req *http.Request) {
	v := mux.Vars(req)["id"]
	if v == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	l, err := s.Store.GetURL(context.Background(), v)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if wbURL, ok := IsUP(l.URL); !ok {
		l.URL = wbURL
	}
	go s.Store.updateLastViewed(context.TODO(), v)
	http.Redirect(rw, req, l.URL, http.StatusFound)
}

func (s *Server) GenerateLink(rw http.ResponseWriter, req *http.Request) {
	data := make(map[string]string)
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	pURL := prefixHTTP(data["url"])

	count, err := s.Store.getCount(context.Background())
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	nl, err := NewLink(pURL, count)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := s.Store.AddURL(context.Background(), nl); err != nil {
		if err != linkAlreadyExistsError {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	_, err = s.Store.incrementCounter(context.Background())
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
	res := map[string]string{
		"genURL": fmt.Sprintf("%s/%s", s.URL, nl.Short),
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(res); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
}

func (s *Server) UpdateLink(rw http.ResponseWriter, req *http.Request) {
	data := make(map[string]string)
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	id := strings.Split(data["url"], "/")[3]
	nURL := prefixHTTP(data["new"])
	if err := s.Store.UpdateURL(context.Background(), id, nURL); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	l, err := s.Store.GetURLInfo(context.Background(), id)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(l); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
}

func (s *Server) InfoLink(rw http.ResponseWriter, req *http.Request) {
	v := mux.Vars(req)["id"]
	if v == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	l, err := s.Store.GetURLInfo(context.Background(), v)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(l); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
	}
}
