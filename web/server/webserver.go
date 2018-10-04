package main

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/mux"
	"github.com/junhsieh/goexamples/fieldbinding/fieldbinding"

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
	Sender         string
}

var db *sql.DB
var mc *memcache.Client

func main() {
	//MEMCACHED
	mc = memcache.New("192.168.12.3:11211")

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
	router.HandleFunc("/tx/{txid:[a-fA-F0-9]{64}}", TransactionHandler).
		Methods("GET")
	router.HandleFunc("/api/positions", getPositions).
		Methods("GET")
	router.HandleFunc("/api/login", getlogin).
		Methods("POST")

	// Static
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/web")))

	http.ListenAndServe(":8000", router)
	log.Println("Listening...")
}

type LoginPost struct {
	Username     string
	PasswordHash string
}

func getlogin(w http.ResponseWriter, r *http.Request) {
	var loginPost LoginPost
	login := false

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&loginPost)
	log.Printf("Login accessed: %s %s", loginPost.Username, loginPost.PasswordHash)
	//if
	login = true

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(login)
}

func getPositions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("accessed getPositions")
	txs, err := getPositionsFromBackend()
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txs)
}

func getPositionsFromBackend() ([]Tx, error) {
	var txs []Tx
	var errCache error
	var err error
	var cache *memcache.Item

	sql_query := "SELECT txid,prefix,hash,type,title,blocktimestamp,blockheight,sender FROM prefix_0xe901"
	cache_key := hasher(sql_query)
	cache, errCache = get_cache(cache_key)
	if errCache != nil {
		query, err := db.Query(sql_query)
		if err != nil {
			return nil, err
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
			var sender string

			err = query.Scan(
				&txid,
				&prefix,
				&link,
				&dataType,
				&title,
				&blockTimestamp,
				&blockHeight,
				&sender)

			txs = append(txs,
				Tx{
					Txid:           txid,
					Prefix:         prefix,
					Link:           link,
					DataType:       dataType,
					Title:          title,
					BlockTimestamp: blockTimestamp,
					BlockHeight:    blockHeight,
					Sender:         sender})
		}
		cacheBytes := new(bytes.Buffer)
		json.NewEncoder(cacheBytes).Encode(txs)
		err = set_cache2(cache_key, cacheBytes.Bytes(), 5)
		if err != nil {
			fmt.Println("Set Cache Error:", err)
		}
	} else {
		json.Unmarshal(cache.Value, &txs)
	}

	return txs, err
}

func selectFromMysql2() ([]Tx, error) {
	var txs []Tx
	var errCache error
	var err error
	var cache *memcache.Item

	sql_query := "SELECT txid,prefix,hash,type,title,blocktimestamp,blockheight,sender FROM prefix_0xe901"
	cache_key := hasher(sql_query)
	cache, errCache = get_cache(cache_key)
	if errCache != nil {
		query, err := db.Query(sql_query)
		if err != nil {
			return nil, err
		}
		defer query.Close()

		var fArr []string
		fb := fieldbinding.NewFieldBinding()
		fArr, err = query.Columns()
		if err != nil {
			return nil, err
		}
		fb.PutFields(fArr)

		outArr := []interface{}{}

		for query.Next() {
			err := query.Scan(fb.GetFieldPtrArr()...)
			if err != nil {
				return nil, err
			}
			outArr = append(outArr, fb.GetFieldArr())
		}

		cacheBytes := new(bytes.Buffer)
		json.NewEncoder(cacheBytes).Encode(outArr)
		err = set_cache2(cache_key, cacheBytes.Bytes(), 5)
		if err != nil {
			fmt.Println("Set Cache Error:", err)
		}
	} else {
		json.Unmarshal(cache.Value, &txs)
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cache_val)
}

func hasher(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func get_cache(key string) (*memcache.Item, error) {
	val, err := mc.Get(key)
	return val, err
}

func set_cache(key string, value string, expiretime int32) error {
	fmt.Println("set key:", key)
	err := mc.Set(&memcache.Item{Key: key, Value: []byte(value), Expiration: expiretime})
	return err
}

func set_cache2(key string, value []byte, expiretime int32) error {
	fmt.Println("set key:", key)
	err := mc.Set(&memcache.Item{Key: key, Value: value, Expiration: expiretime})
	return err
}
