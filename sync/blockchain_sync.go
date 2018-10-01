package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Query struct {
	Unconfirmed []Transaction `json:"unconfirmed"`
	Confirmed   []Transaction `json:"confirmed"`
}

type Transaction struct {
	Tx  Id    `json:"tx"`
	Out []Sub `json:"out"`
	Blk Info  `json:"blk"`
}

type Sub struct {
	B1 string `json:"b1"`
	S2 string `json:"s2"`
	S3 string `json:"s3"`
	S4 string `json:"s4"`
}

type Info struct {
	T uint32 `json:"t"`
	I uint32 `json:"i"`
}

type Id struct {
	H string `json:"h"`
}

func insertIntoMysql(TxId string, prefix string, hash string, data_type string, title string, blocktimestamp uint32, blockheight uint32) bool {
	fmt.Println("==> ", blockheight, blocktimestamp, TxId, prefix, hash, data_type, title)
	//Mysql
	// db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/theca")
	db, err := sql.Open("mysql", "root:8drRNG8RWw9FjzeJuavbY6f9@tcp(172.18.0.5:3306)/theca")
	if err != nil {
		return false
	}
	defer db.Close()

	sql_query := "INSERT INTO opreturn VALUES(?,?,?,?,?,?,?)"
	insert, err := db.Prepare(sql_query)
	defer insert.Close()

	_, err = insert.Query(TxId, prefix, hash, data_type, title, blocktimestamp, blockheight)
	if err != nil {
		return false
	}
	return true
}

func main() {
	var q Query
	var ScannerBlockHeight uint32
	var LastScannerBlockHeight uint32
	ScannerBlockHeight = 550000
	LastScannerBlockHeight = 0

	fmt.Println("Start ScannerBlockHeight: >", ScannerBlockHeight)

	loop := true
	for loop {
		var query = `{
		  "v": 2,
		  "e": { "out.b1": "hex"  },
		  "q": {
		    "find": {
		      "out.b1": "e901",
		      "out.b0": {
		        "op": 106
		      },
		      "blk.i": {
		        "$gte" : %d
		      }

		    },
		    "project": {
		      "out.b1": 1,
		      "out.s2": 1,
		      "out.s3": 1,
		      "out.s4": 1,
		      "tx.h": 1,
					"blk.t": 1,
					"blk.i": 1,
		      "_id": 0
		    }
		  }
		}`

		query = fmt.Sprintf(query, ScannerBlockHeight)

		if LastScannerBlockHeight != ScannerBlockHeight {
			fmt.Println("ScannerHeight: >", ScannerBlockHeight)
			LastScannerBlockHeight = ScannerBlockHeight
		}

		//url encoded query : blocksize greater than 550'000
		b64_query := base64.StdEncoding.EncodeToString([]byte(query))
		url := "https://bitdb.network/q/" + b64_query
		client := &http.Client{}
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("key", "qz6qzfpttw44eqzqz8t2k26qxswhff79ng40pp2m44")

		res, _ := client.Do(req)

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}

		json.Unmarshal(body, &q)

		for i := range q.Confirmed {
			TxId := q.Confirmed[i].Tx.H
			txOuts := q.Confirmed[i].Out
			BlockTimestamp := q.Confirmed[i].Blk.T
			BlockHeight := q.Confirmed[i].Blk.I
			var Prefix string
			var Hash string
			var Datatype string
			var Title string
			for a := range txOuts {
				if txOuts[a].B1 == "e901" {
					Prefix = txOuts[a].B1
					Hash = txOuts[a].S2
					Datatype = txOuts[a].S3
					Title = txOuts[a].S4
				}
			}
			if BlockHeight > ScannerBlockHeight {
				ScannerBlockHeight = BlockHeight + 1
			}

			if len(Prefix) != 0 && len(Hash) > 20 && len(Datatype) > 2 {
				insert := insertIntoMysql(TxId, Prefix, Hash, Datatype, Title, BlockTimestamp, BlockHeight)
				if insert != true {
					fmt.Println("Insert failed ! (error or duplicated db entry)")
				} else {
					fmt.Println("Insert OK")
				}
			}
		}
		time.Sleep(5 * time.Second)
	}

}
