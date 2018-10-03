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

var db *sql.DB

func selectUnconfiremedMysql() ([]string, error) {
	var uc_txs []string
	sql_query_uc := "SELECT txid FROM opreturn WHERE blockheight = 0"
	uc_query, err := db.Query(sql_query_uc)
	if err != nil {
		return uc_txs, err
	}
	defer uc_query.Close()

	for uc_query.Next() {
		var uc_txid string
		err = uc_query.Scan(&uc_txid)
		uc_txs = append(uc_txs, uc_txid)
	}
	return uc_txs, err
}
func updateMysql(TxId string, blocktimestamp uint32, blockheight uint32) error {
	sql_update := "UPDATE opreturn SET blockheight=?,blocktimestamp=? where txid=?"
	update, err := db.Prepare(sql_update)
	defer update.Close()

	_, err = update.Exec(blockheight, blocktimestamp, TxId)
	return err
}

func insertIntoMysql(TxId string, prefix string, hash string, data_type string, title string, blocktimestamp uint32, blockheight uint32) error {
	sql_query := "INSERT INTO opreturn VALUES(?,?,?,?,?,?,?)"
	insert, err := db.Prepare(sql_query)
	defer insert.Close()

	_, err = insert.Query(TxId, prefix, hash, data_type, title, blocktimestamp, blockheight)
	return err
}

func main() {
	var err error
	db, err = sql.Open("mysql", "root:8drRNG8RWw9FjzeJuavbY6f9@tcp(192.168.12.2:3306)/theca")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var q Query
	var ScannerBlockHeight uint32
	var LastScannerBlockHeight uint32
	ScannerBlockHeight = 550255
	LastScannerBlockHeight = 0

	fmt.Println("Start ScannerBlockHeight: >", ScannerBlockHeight)

	loop := true
	for loop {
		//list of unconfirmed tx in db
		unconfirmedInDb, err := selectUnconfiremedMysql()

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

		var BlockHeight uint32

		for i := range q.Confirmed {
			TxId := q.Confirmed[i].Tx.H
			txOuts := q.Confirmed[i].Out
			BlockTimestamp := q.Confirmed[i].Blk.T
			BlockHeight = q.Confirmed[i].Blk.I
			if BlockHeight > ScannerBlockHeight {
				ScannerBlockHeight = BlockHeight + 1
			}
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

			if len(Prefix) != 0 && len(Hash) > 20 && len(Datatype) > 2 {
				exists := false
				for i := range unconfirmedInDb {
					uc_txid := unconfirmedInDb[i]
					if uc_txid == TxId {
						exists = true
					}
				}

				if exists {
					err := updateMysql(TxId, BlockTimestamp, BlockHeight)
					if err != nil {
						fmt.Println("UPDATE FAILED (confirmed) error")
					} else {
						fmt.Println("UPDATE OK (confirmed)==> ", TxId, Prefix, Hash, Datatype, Title, BlockTimestamp, BlockHeight)
					}
				} else {
					err := insertIntoMysql(TxId, Prefix, Hash, Datatype, Title, BlockTimestamp, BlockHeight)
					if err != nil {
						fmt.Println("INSERT DUP / FAILED (confirmed) error or duplicated db entry")
					} else {
						fmt.Println("INSERT OK (confirmed)==> ", TxId, Prefix, Hash, Datatype, Title, BlockTimestamp, BlockHeight)
					}
				}
			}
		}

		// unconfirmed transactions
		var uc_query = `{
			"v": 2,
			"e": { "out.b1": "hex"  },
			"q": {
				"find": {
					"out.b1": "e901",
					"out.b0": {
						"op": 106
					}
				},
				"project": {
					"out.b1": 1,
					"out.s2": 1,
					"out.s3": 1,
					"out.s4": 1,
					"tx.h": 1,
					"_id": 0
				}
			}
		}`

		//url encoded query : blocksize greater than 550'000
		b64_uc_query := base64.StdEncoding.EncodeToString([]byte(uc_query))
		uc_url := "https://bitdb.network/q/" + b64_uc_query
		uc_client := &http.Client{}
		uc_req, _ := http.NewRequest("GET", uc_url, nil)
		uc_req.Header.Set("key", "qz6qzfpttw44eqzqz8t2k26qxswhff79ng40pp2m44")

		uc_res, _ := uc_client.Do(uc_req)

		uc_body, err := ioutil.ReadAll(uc_res.Body)

		if err != nil {
			log.Fatalln(err)
		}

		json.Unmarshal(uc_body, &q)
		for i := range q.Unconfirmed {
			TxId := q.Unconfirmed[i].Tx.H
			txOuts := q.Unconfirmed[i].Out
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

			if len(Prefix) != 0 && len(Hash) > 20 && len(Datatype) > 2 {
				exists := false
				for i := range unconfirmedInDb {
					uc_txid := unconfirmedInDb[i]
					if uc_txid == TxId {
						exists = true
					}
				}
				if !exists {
					err := insertIntoMysql(TxId, Prefix, Hash, Datatype, Title, 0, 0)
					if err != nil {
						fmt.Println("INSERT FAILED (unconfirmed): error or duplicated db entry")
					} else {
						fmt.Println("INSERT OK (unconfirmed)==> ", TxId, Prefix, Hash, Datatype, Title)
					}
				}
			}
		}

		time.Sleep(10 * time.Second)
	}

}
