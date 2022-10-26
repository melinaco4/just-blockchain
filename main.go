package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Block struct {
	Position  int
	Data      SongCheckout
	TimeStamp string
	Hash      string
	PrevHash  string
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

func (b *Block) generateHash() {
	bytes, _ := json.Marshal(b.Data)

	data := string(b.Position) + b.TimeStamp + string(bytes) + b.PrevHash

	hash := sha256.New()
	hash.Write([]byte(data))

	b.Hash = hex.EncodeToString(hash.Sum(nil))

}

func CreateBlock(prevBlock *Block, checkoutitem SongCheckout) *Block {
	block := &Block{}
	block.Position = prevBlock.Position + 1
	block.TimeStamp = time.Now().String()
	block.Data = checkoutitem
	block.PrevHash = prevBlock.Hash
	block.generateHash()

	return block

}

var BlockChain *Blockchain

func (bc *Blockchain) AddBlock(data SongCheckout) {
	prevBlock := bc.blocks[len(bc.blocks)-1]

	block := CreateBlock(prevBlock, data)

	if validBlock(block, prevBlock) {
		bc.blocks = append(bc.blocks, block)
	}
}

func validBlock(block, prevBlock *Block) bool {

	if prevBlock.Hash != block.PrevHash {
		return false
	}

	if !block.validateHash(block.Hash) {
		return false
	}

	if prevBlock.Position+1 != block.Position {
		return false
	}

	return true
}

func (b *Block) validateHash(hash string) bool {
	b.generateHash()
	if b.Hash != hash {
		return false
	}

	return true
}

func writeBlock(w http.ResponseWriter, r *http.Request) {
	var checkoutitem SongCheckout

	if err := json.NewDecoder(r.Body).Decode(&checkoutitem); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not write block:%v", err)
		w.Write([]byte("could not write block"))
	}

	BlockChain.AddBlock(checkoutitem)

}

func newSong(w http.ResponseWriter, r *http.Request) {
	var song Song

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not create: %v", err)
		w.Write([]byte("could not create new song"))
		return
	}

	h := md5.New()
	io.WriteString(h, song.SKU+song.ReleaseDate)
	song.ID = fmt.Sprintf("%x", h.Sum(nil))

	resp, err := json.MarshalIndent(song, "", "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not marshal payload: %v", err)
		w.Write([]byte("could not save song daa"))
		return
	}
	w.WriteHeader((http.StatusOK))
	w.Write(resp)

}

func GenesisBlock() *Block {
	return CreateBlock(&Block{}, SongCheckout{IsGenesis: true})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{GenesisBlock()}}
}

func getBlockchain(w http.ResponseWriter, r *http.Request) {
	jbytes, err := json.MarshalIndent(BlockChain.blocks, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	io.WriteString(w, string(jbytes))
}

func main() {

	BlockChain = NewBlockchain()

	r := mux.NewRouter()
	r.HandleFunc("/", getBlockchain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newSong).Methods("POST")

	go func() {
		for _, block := range BlockChain.blocks {
			fmt.Printf("Prev. hash: %x\n", block.PrevHash)
			bytes, _ := json.MarshalIndent(block.Data, "", " ")
			fmt.Printf("Data:%v\n", string(bytes))
			fmt.Print("Hash:%x\n", block.Hash)
			fmt.Println()
		}
	}()

	log.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", r))
}
