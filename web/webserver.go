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
	"os"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/mux"
	"github.com/pmylund/sortutil"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var mc *memcache.Client

func main() {
	//MEMCACHED
	mc = memcache.New(os.Getenv("MEMCACHE_HOSTNAME") + ":" + os.Getenv("MEMCACHE_PORT"))

	//MYSQL
	var err error
	db, err = sql.Open("mysql", "root:8drRNG8RWw9FjzeJuavbY6f9@tcp("+os.Getenv("DATABASE_HOSTNAME")+":"+os.Getenv("DATABASE_PORT")+")/"+os.Getenv("DATABASE_NAME"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//Router
	router := mux.NewRouter()

	//Response
	router.HandleFunc(ThecaMainApiPath, getPositions).
		Methods("GET")
	router.HandleFunc(CommentApiPath, getComments).
		Methods("GET")
	router.HandleFunc(LoginApiPath, postLogin).
		Methods("POST")
	router.HandleFunc(SignupApiPAth, postSignup).
		Methods("POST")
	router.HandleFunc(TxApiPath, getTransactionData).
		Methods("GET")

	// Static
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(FileServerHandlePath)))

	http.ListenAndServe(MuxPort, router)
	log.Println("Listening...")
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
	sql_query := MysqlInsertUser
	insert, err := db.Prepare(sql_query)
	defer insert.Close()
	_, err = insert.Query(userName, passwordHash, encryptedPk)
	return err
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

	sql_query := fmt.Sprintf(MysqlSelectUserLogin, userName, passwordHash)
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
		err = set_cache(cache_key, cacheBytes.Bytes(), 5)
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

func calculateScore(likes uint32, timestamp int64) float64 {
	score := float64(0)
	gravity := float64(1.8)
	now := time.Now().Unix()
	hours := float64(now-int64(timestamp)) / 3600
	if likes > 0 {
		score = float64(likes-1) / math.Pow((hours+2), gravity)
	}
	return score
}

func getComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txid := vars["txid"]
	fmt.Println("accessed getComments", txid)
	comments, err := getCommentsFromBackend(txid)
	if err != nil {
		log.Fatal(err)
	}
	for comment := range comments {
		if comments[comment].BlockTimestamp == 0 {
			comments[comment].BlockTimestamp = comments[comment].Timestamp
		}
		comments[comment].Score = calculateScore(comments[comment].Likes, comments[comment].BlockTimestamp)
	}
	sortutil.DescByField(comments, "Score")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func getPositions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("accessed getPositions")
	txs, err := getPositionsFromBackend()
	fmt.Println(txs)
	for tx := range txs {
		if txs[tx].BlockTimestamp == 0 {
			txs[tx].BlockTimestamp = txs[tx].Timestamp
		}
		txs[tx].Score = calculateScore(txs[tx].Likes, txs[tx].BlockTimestamp)
	}
	sortutil.DescByField(txs, "Score")
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txs)
}

func getTransactionDataFromBackend(txid string) (Theca, error) {
	var tx Theca
	var errCache error
	var err error
	var cache *memcache.Item

	sql_query := fmt.Sprintf(MysqlSelectThecaTx, txid)
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
			var blockTimestamp int64
			var blockHeight uint32
			var sender string
			var timestamp int64

			err = query.Scan(
				&txid,
				&link,
				&dataType,
				&title,
				&blockTimestamp,
				&blockHeight,
				&sender,
				&timestamp)

			tx = Theca{
				Txid:           txid,
				Link:           link,
				DataType:       dataType,
				Title:          title,
				BlockTimestamp: blockTimestamp,
				BlockHeight:    blockHeight,
				Sender:         sender,
				Timestamp:      timestamp}
		}
		cacheBytes := new(bytes.Buffer)
		json.NewEncoder(cacheBytes).Encode(tx)
		err = set_cache(cache_key, cacheBytes.Bytes(), 5)
		if err != nil {
			fmt.Println("Set Cache Error:", err)
		}
	} else {
		json.Unmarshal(cache.Value, &tx)
	}
	return tx, err
}

func getCommentsFromBackend(txid string) ([]Comment, error) {
	var txs []Comment
	var errCache error
	var err error
	var cache *memcache.Item
	sql_query := fmt.Sprintf(MysqlSelectComments, txid)
	fmt.Println(sql_query)
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
			var txhash string
			var message string
			var blockTimestamp int64
			var blockHeight uint32
			var sender string
			var timestamp int64

			err = query.Scan(
				&txid,
				&txhash,
				&message,
				&blockTimestamp,
				&blockHeight,
				&sender,
				&timestamp)

			txs = append(txs,
				Comment{
					Txid:           txid,
					TxHash:         txhash,
					Message:        message,
					BlockTimestamp: blockTimestamp,
					BlockHeight:    blockHeight,
					Sender:         sender,
					Timestamp:      timestamp})
		}
		cacheBytes := new(bytes.Buffer)
		json.NewEncoder(cacheBytes).Encode(txs)
		err = set_cache(cache_key, cacheBytes.Bytes(), 5)
		if err != nil {
			fmt.Println("Set Cache Error:", err)
		}
	} else {
		json.Unmarshal(cache.Value, &txs)
	}
	return txs, err
}

func getPositionsFromBackend() ([]Theca, error) {
	var txs []Theca
	var errCache error
	var err error
	var cache *memcache.Item

	sql_query := MysqlSelectThecaTxList
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
			var blockTimestamp int64
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
				Theca{
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
		err = set_cache(cache_key, cacheBytes.Bytes(), 5)
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

func set_cache(key string, value []byte, expiretime int32) error {
	fmt.Println("set key:", key)
	err := mc.Set(&memcache.Item{Key: key, Value: value, Expiration: expiretime})
	return err
}
