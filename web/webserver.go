package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	//Response
	router.HandleFunc("/", HomeHandler).Methods("GET")
	router.HandleFunc("/tx/{txid:[a-fA-F0-9]{64}}", TransactionHandler).Methods("GET")
	http.ListenAndServe(":8000", router)

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("accessed HomeHandler")
}

func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("accessed TransactionHandler")
}

func hasher(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func get_cache(key string) (*memcache.Item, error) {
	mc := memcache.New("127.0.0.1:11211")
	val, err := mc.Get(key)
	return val, err
}

func set_cache(key string, value string, expiretime int32) error {
	fmt.Println("set:", key)
	mc := memcache.New("127.0.0.1:11211")
	err := mc.Set(&memcache.Item{Key: key, Value: []byte(value), Expiration: expiretime})
	return err
}
