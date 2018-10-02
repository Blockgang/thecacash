package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	//Response
	router.HandleFunc("/meminfo", Getmeminfo).Methods("GET")
	router.HandleFunc("/getmysql", Getmysql).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))

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
