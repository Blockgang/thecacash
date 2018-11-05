package main

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/mux"

	"github.com/junhsieh/goexamples/fieldbinding/fieldbinding"

	_ "github.com/go-sql-driver/mysql"
)

type Tx struct {
	Txid           string
	Link           string
	DataType       string
	Title          string
	BlockTimestamp uint32
	BlockHeight    uint32
	Sender         string
	Likes          uint32
	Comments       uint32
	Score          float64
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
	router.HandleFunc("/api/tx/positions", getPositions).
		Methods("GET")
	router.HandleFunc("/api/login", postLogin).
		Methods("POST")
	router.HandleFunc("/api/signup", postSignup).
		Methods("POST")
	router.HandleFunc("/api/tx/{txid:[a-fA-F0-9]{64}}", getTransactionData).
		Methods("GET")

	// Static
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("/var/www/")))

	http.ListenAndServe(":8000", router)
	log.Println("Listening...")
}

type SignupPost struct {
	Username     string
	PasswordHash string
	EncryptedPk  string
}

type SignupResponse struct {
	Username    string
	EncryptedPk string
	Signup      bool
}

func postSignup(w http.ResponseWriter, r *http.Request) {
	var signupPost SignupPost

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&signupPost)
	log.Printf("Signup accessed: %s %s %s", signupPost.Username, signupPost.PasswordHash, signupPost.EncryptedPk)

	res, err := signup(signupPost.Username, signupPost.PasswordHash, signupPost.EncryptedPk)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func signup(userName string, passwordHash string, encryptedPk string) (SignupResponse, error) {
	var err error
	var signup bool

	encryptedKey, err := check_login(userName, passwordHash)
	if err != nil {
		fmt.Println(err)
	}
	if len(encryptedKey) > 0 {
		signup = false
	} else {
		signup = true
		// MYSQL Insert
		err = insertLoginIntoMysql(userName, passwordHash, encryptedPk)
	}

	res := SignupResponse{
		Username:    userName,
		EncryptedPk: encryptedPk,
		Signup:      signup,
	}
	return res, err
}

func insertLoginIntoMysql(userName string, passwordHash string, encryptedPk string) error {
	sql_query := "INSERT INTO users (username,password,encrypted_pk) VALUES(?,?,?)"
	insert, err := db.Prepare(sql_query)
	defer insert.Close()
	_, err = insert.Query(userName, passwordHash, encryptedPk)
	return err
}

type LoginPost struct {
	Username     string
	PasswordHash string
}

type LoginResponse struct {
	Username    string
	EncryptedPk string
	Login       bool
}

func postLogin(w http.ResponseWriter, r *http.Request) {
	var loginPost LoginPost
	var login bool

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&loginPost)
	log.Printf("Login accessed: %s %s", loginPost.Username, loginPost.PasswordHash)
	encryptedKey, err := check_login(loginPost.Username, loginPost.PasswordHash)
	if err != nil {
		log.Fatal(err)
	}
	if len(encryptedKey) > 0 {
		fmt.Println(encryptedKey, err)
		login = true
	} else {
		login = false
	}

	res := LoginResponse{
		Username:    loginPost.Username,
		EncryptedPk: encryptedKey,
		Login:       login,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func check_login(userName string, passwordHash string) (string, error) {
	var errCache error
	var err error
	var cache *memcache.Item
	var encryptedKey string

	sql_query := fmt.Sprintf("SELECT encrypted_pk FROM users WHERE username='%s' AND password='%s'", userName, passwordHash)
	cache_key := hasher(sql_query)
	cache, errCache = get_cache(cache_key)
	if errCache != nil {
		query, err := db.Query(sql_query)
		if err != nil {
			return encryptedKey, err
		}
		defer query.Close()

		for query.Next() {
			err = query.Scan(&encryptedKey)
		}
		cacheBytes := new(bytes.Buffer)
		json.NewEncoder(cacheBytes).Encode(encryptedKey)
		err = set_cache2(cache_key, cacheBytes.Bytes(), 5)
		if err != nil {
			fmt.Println("Set Cache Error:", err)
		}
	} else {
		json.Unmarshal(cache.Value, &encryptedKey)
	}
	return encryptedKey, err
}

func getTransactionData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txid := vars["txid"]
	fmt.Println("accessed getTransactionData", txid)
	txData, err := getTransactionDataFromBackend(txid)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txData)
}

func calculateScore(likes uint32, timestamp uint32) float64 {
	score := float64(0)
	gravity := float64(1.8)
	now := time.Now().Unix()
	hours := float64(now-int64(timestamp)) / 3600
	if likes > 0 {
		score = float64(likes-1) / math.Pow((hours+2), gravity)
	}
	return score
}

func getPositions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("accessed getPositions")
	txs, err := getPositionsFromBackend()
	fmt.Println(txs)
	for tx := range txs {
		txs[tx].Score = calculateScore(txs[tx].Likes, txs[tx].BlockTimestamp)
	}
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txs)
}

func getTransactionDataFromBackend(txid string) (Tx, error) {
	var tx Tx
	var errCache error
	var err error
	var cache *memcache.Item

	sql_query := fmt.Sprintf("SELECT txid,hash,type,title,blocktimestamp,blockheight,sender FROM prefix_0xe901 WHERE txid='%s'", txid)
	fmt.Println(sql_query)
	cache_key := hasher(sql_query)
	cache, errCache = get_cache(cache_key)
	if errCache != nil {
		query, err := db.Query(sql_query)
		if err != nil {
			return tx, err
		}
		defer query.Close()

		for query.Next() {
			var txid string
			var link string
			var dataType string
			var title string
			var blockTimestamp uint32
			var blockHeight uint32
			var sender string

			err = query.Scan(
				&txid,
				&link,
				&dataType,
				&title,
				&blockTimestamp,
				&blockHeight,
				&sender)

			tx = Tx{
				Txid:           txid,
				Link:           link,
				DataType:       dataType,
				Title:          title,
				BlockTimestamp: blockTimestamp,
				BlockHeight:    blockHeight,
				Sender:         sender}
		}
		cacheBytes := new(bytes.Buffer)
		json.NewEncoder(cacheBytes).Encode(tx)
		err = set_cache2(cache_key, cacheBytes.Bytes(), 5)
		if err != nil {
			fmt.Println("Set Cache Error:", err)
		}
	} else {
		json.Unmarshal(cache.Value, &tx)
	}
	return tx, err
}

func getPositionsFromBackend() ([]Tx, error) {
	var txs []Tx
	var errCache error
	var err error
	var cache *memcache.Item

	sql_query := "SELECT txid,hash,type,title,blocktimestamp,blockheight,sender,likes,comments FROM prefix_0xe901"
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
			var link string
			var dataType string
			var title string
			var blockTimestamp uint32
			var blockHeight uint32
			var sender string
			var likes uint32
			var comments uint32

			err = query.Scan(
				&txid,
				&link,
				&dataType,
				&title,
				&blockTimestamp,
				&blockHeight,
				&sender,
				&likes,
				&comments)

			txs = append(txs,
				Tx{
					Txid:           txid,
					Link:           link,
					DataType:       dataType,
					Title:          title,
					BlockTimestamp: blockTimestamp,
					BlockHeight:    blockHeight,
					Sender:         sender,
					Likes:          likes,
					Comments:       comments})
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

	sql_query := "SELECT txid,hash,type,title,blocktimestamp,blockheight,sender FROM prefix_0xe901"
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
