package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/mux"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func main() {
	router := mux.NewRouter()
	//Response
	router.HandleFunc("/tx/{txid:[a-fA-F0-9]{64}}", TransactionHandler).Methods("GET")
	router.HandleFunc("/template", templatehandler).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/server/static")))
	http.ListenAndServe(":8000", router)
	log.Println("Listening...")
}

func templatehandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("accessed TemplateHandler")
	data := TodoPageData{
		PageTitle: "My list",
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}
	t, _ := template.ParseFiles("./web/server/templates/example.html")
	t.Execute(w, data)
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
