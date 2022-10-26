package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Block struct {
	Position     int
	Data         SongCheckout
	TimeStamp    string
	Hash         string
	PreviousHash string
}

type SongCheckout struct {
	SongID       string `json:"song_id"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis    bool   `jsons:"is_genesis"`
}

type Song struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	ReleaseDate string `json:"release_date"`
	SKU         string `json:"sku"`
}

type Blockchain struct {
	blocks []*Block
}

var Blockchain *Blockchain

func main() {
	r := mux.NewRouter()
	r.Handle("/", getBlockchain).Methods("GET")
	r.Handle("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newSong).Methods("POST")

	log.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", r))
}
