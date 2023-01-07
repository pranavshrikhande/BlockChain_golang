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
	Pos       int
	Data      BookCheckout
	TimeStamp string
	Hash      string
	PrevHash  string
}

type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	PublishData string `json:"publish_date"`
	ISBN        string `json:"isbn"`
}

type BookCheckout struct {
	BookID       string `json:"book_id"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis    bool   `json:"is_genesis"`
}

type Blockchain struct {
	blocks []*Block
}

var BlockChain *Blockchain

// this function creates hash
func (b *Block) generateHash() {
	bytes, _ := json.Marshal(b.Data)

	//elements to consider to create hash
	data := string(b.Pos) + b.TimeStamp + string(bytes) + b.PrevHash

	hash := sha256.New()

	hash.Write([]byte(data))

	b.Hash = hex.EncodeToString(hash.Sum(nil))
}

func CreateBlock(prevBlock *Block, checkoutitem BookCheckout) *Block {

	//1 Create block, that needs to be returned

	block := &Block{}

	block.Pos = prevBlock.Pos + 1
	block.TimeStamp = time.Now().String()
	block.PrevHash = prevBlock.Hash
	block.generateHash()

	return block
}

// as we are using mux, we can take req and res
// if we use Gin and Fiber, we dont write this. we just pass the context
func newBook(w http.ResponseWriter, r *http.Request) {
	var book Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not create:%v", err)
		w.Write([]byte("could not create a new book"))
		return
	}
	//to create new ID
	h := md5.New()
	io.WriteString(h, book.ISBN+book.PublishData)
	book.ID = fmt.Sprintf("%x", h.Sum(nil))

	//now we send back in JSON to user
	resp, err := json.MarshalIndent(book, "", " ")

	//while converting book to JSON, there could be error so status 500
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not marshal payload:%v", err)
		w.Write([]byte("could not save book data"))
	}

	//status 200
	w.WriteHeader(http.StatusOK)
	//We create Book with ID, with this ID we use the next route writeBlock
	w.Write(resp)

}

func (bc *Blockchain) AddBlock(data BookCheckout) {

	// to add block to chain we take previous block

	prevBlock := bc.blocks[len(bc.blocks)-1]

	//createBlock function would take the Bookcheckout data and create block
	block := CreateBlock(prevBlock, data)

	//on creating block, check its validity

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

	//to check the discrepancy of position of block
	if prevBlock.Pos+1 != block.Pos {
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
	//to create a block, we first need Bookcheckout data in request

	var checkoutitem BookCheckout

	if err := json.NewDecoder(r.Body).Decode(&checkoutitem); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not write block:%v", err)
		w.Write([]byte("could not write a block"))
	}

	//now the deocoded bookChekout data needs to be sent as a block to Blockchain

	//Here we are not calling function AddBlock, we are calling struct method AddBlock.
	/*
		it looks like
		func(bc *Blockhain) AddBlock(data BookCheckout){


	*/

	BlockChain.AddBlock(checkoutitem)
	resp, err := json.MarshalIndent(checkoutitem, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not marshal payload: %v", err)
		w.Write([]byte("could not write block"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

// what this function does is  create a blockchain with Genesisblock, it creates a block with blockcheckout information
func GenesisBlock() *Block {
	return CreateBlock(&Block{}, BookCheckout{IsGenesis: true})
}

// this what creates blockchain, doesnt take anything but returns blockchain
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{GenesisBlock()}}
}

// It is sending your entire blockchain i.e "BlockChain" and all its blocks using json
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

	//as soon as the program begins, you need to first instantiate a blockchain
	BlockChain = NewBlockchain()

	r := mux.NewRouter()

	r.HandleFunc("/", getBlockchain).Methods("GET")

	//when created a book, also have to pass info such as name,id etc
	r.HandleFunc("/", writeBlock).Methods("POST")

	//this route we hit in the beginning to create a new book
	r.HandleFunc("/new", newBook).Methods("POST")

	//as soon as program starts, below function needs to be executed
	//this will work as go routine
	go func() {
		for _, block := range BlockChain.blocks {
			fmt.Printf("Prev .hash:%x\n", block.PrevHash)
			bytes, _ := json.MarshalIndent(block.Data, "", " ")
			fmt.Printf("Data:%v\n", string(bytes))
			fmt.Printf("Hash:%x\n", block.Hash)
			fmt.Println()
		}
	}()

	log.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", r))
}
