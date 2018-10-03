package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

type Tx struct {
	Txid           string
	Prefix         string
	Link           string
	DataType       string
	Title          string
	BlockTimestamp uint32
	BlockHeight    uint32
}

var db *sql.DB

func main() {
	//MYSQL
	var err error
	db, err = sql.Open("mysql", "root:8drRNG8RWw9FjzeJuavbY6f9@tcp(192.168.12.2:3306)/theca")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//Router
	router := mux.NewRouter()
	//Response
	router.HandleFunc("/tx/{txid:[a-fA-F0-9]{64}}", TransactionHandler).Methods("GET")
	router.HandleFunc("/api/positions", getPositions).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/server/static")))
	http.ListenAndServe(":8000", router)
	log.Println("Listening...")
}

func getPositions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("accessed getPositions")
	txs, err := selectFromMysql()
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(txs)
}

func selectFromMysql() ([]Tx, error) {
	var txs []Tx
	sql_query := "SELECT txid,prefix,hash,type,title,blocktimestamp,blockheight FROM opreturn"
	query, err := db.Query(sql_query)
	if err != nil {
		return txs, err
	}
	defer query.Close()

	for query.Next() {
		var txid string
		var prefix string
		var link string
		var dataType string
		var title string
		var blockTimestamp uint32
		var blockHeight uint32

		err = query.Scan(
			&txid,
			&prefix,
			&link,
			&dataType,
			&title,
			&blockTimestamp,
			&blockHeight)

		txs = append(txs,
			Tx{
				Txid:           txid,
				Prefix:         prefix,
				Link:           link,
				DataType:       dataType,
				Title:          title,
				BlockTimestamp: blockTimestamp,
				BlockHeight:    blockHeight})
	}
	return txs, err
}

func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txid := vars["txid"]
	cache_val := ""
	fmt.Println("accessed TransactionHandler")

	cache, err := get_cache(txid)
	if err != nil {
		cache_val = string(txid)
		if set_cache(txid, cache_val, 10) != nil {
			fmt.Println("Set Cache Error:", err)
		}
	} else {
		cache_val = string(cache.Value)
	}
	fmt.Println("Cache Value:", cache_val)
	json.NewEncoder(w).Encode(cache_val)
}

func hasher(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func get_cache(key string) (*memcache.Item, error) {
	mc := memcache.New("192.168.12.3:11211")
	val, err := mc.Get(key)
	return val, err
}

func set_cache(key string, value string, expiretime int32) error {
	fmt.Println("set key:", key)
	mc := memcache.New("192.168.12.3:11211")
	err := mc.Set(&memcache.Item{Key: key, Value: []byte(value), Expiration: expiretime})
	return err
}
