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
	Tx  Id       `json:"tx"`
	Out []OutSub `json:"out"`
	In  []InSub  `json:"in"`
	Blk Info     `json:"blk"`
}

type OutSub struct {
	B1 string `json:"b1"`
	B2 string `json:"b2"`
	S2 string `json:"s2"`
	S3 string `json:"s3"`
	S4 string `json:"s4"`
}

type InSub struct {
	E Sender `json:"e"`
}

type Sender struct {
	A string `json:"a"`
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
	sql_query_uc := "SELECT txid FROM prefix_0xe901 WHERE blockheight = 0"
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
	sql_update := "UPDATE prefix_0xe901 SET blockheight=?,blocktimestamp=? where txid=?"
	update, err := db.Prepare(sql_update)
	defer update.Close()

	_, err = update.Exec(blockheight, blocktimestamp, TxId)
	return err
}

func insertIntoMysql(TxId string, prefix string, hash string, data_type string, title string, blocktimestamp uint32, blockheight uint32, sender string) error {
	sql_query := "INSERT INTO prefix_0xe901 VALUES(?,?,?,?,?,?,?,?)"
	insert, err := db.Prepare(sql_query)
	defer insert.Close()
	_, err = insert.Query(TxId, prefix, hash, data_type, title, blocktimestamp, blockheight, sender)
	return err
}

func insertMemoLikeIntoMysql(TxId string, txHash string, Sender string, BlockTimestamp uint32, BlockHeight uint32) error {
	sql_query := "INSERT INTO prefix_0x6d04 VALUES(?,?,?,?,?)"
	fmt.Println(TxId, txHash, BlockTimestamp, BlockHeight, Sender)
	time.Sleep(10 * time.Millisecond) //todo: ansonsten too many connections
	insert, err := db.Prepare(sql_query)
	if err != nil {
		fmt.Println(err)
	}
	_, err = insert.Query(TxId, txHash, BlockTimestamp, BlockHeight, Sender)
	insert.Close()
	return err
}

func getMemoLikes(ScannerBlockHeight uint32) {
	var q Query
	var memoLikesQuery = `{
		"v": 2,
		"e": {
			"out.b1": "hex",
			"out.b2": "hex"
		},
		"q": {
			"find": {
				"out.b1": "6d04",
				"out.b0": {
					"op": 106
				},
				"blk.i": {
					"$gte" : %d
				}

			},
			"limit":5000,
			"project": {
				"out.b1": 1,
				"out.b2": 1,
				"tx.h": 1,
				"blk.t": 1,
				"blk.i": 1,
				"in.e.a":1,
				"_id": 0
			}
		}
	}`

	memoLikesQuery = fmt.Sprintf(memoLikesQuery, ScannerBlockHeight)

	b64_query := base64.StdEncoding.EncodeToString([]byte(memoLikesQuery))
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
		Sender := q.Confirmed[i].In[0].E.A
		TxId := q.Confirmed[i].Tx.H
		txOuts := q.Confirmed[i].Out
		BlockTimestamp := q.Confirmed[i].Blk.T
		BlockHeight = q.Confirmed[i].Blk.I
		// if BlockHeight > ScannerBlockHeight {
		// 	ScannerBlockHeight = BlockHeight + 1
		// }
		var Prefix string
		var Hash string
		for a := range txOuts {
			if txOuts[a].B1 == "6d04" {
				Prefix = txOuts[a].B1
				Hash = txOuts[a].B2
			}
		}

		if len(Prefix) != 0 && len(Hash) > 20 {
			exists := false
			//todo: unconfirmed memo likes
			// for i := range unconfirmedInDb {
			// 	uc_txid := unconfirmedInDb[i]
			// 	if uc_txid == TxId {
			// 		exists = true
			// 	}
			// }

			if exists {
				err := updateMysql(TxId, BlockTimestamp, BlockHeight)
				if err != nil {
					fmt.Println("UPDATE FAILED (confirmed) error")
				} else {
					fmt.Println("UPDATE OK (confirmed)==> ", TxId, Hash, Sender, BlockTimestamp, BlockHeight)
				}
			} else {
				err := insertMemoLikeIntoMysql(TxId, Hash, Sender, BlockTimestamp, BlockHeight)
				if err != nil {
					fmt.Println("INSERT DUP / FAILED (confirmed) error or duplicated db entry")
				} else {
					fmt.Println("INSERT OK (confirmed)==> ", TxId, Hash, Sender, BlockTimestamp, BlockHeight)
				}
			}
		}
	}
}

func main() {
	var err error
	db, err = sql.Open("mysql", "root:8drRNG8RWw9FjzeJuavbY6f9@tcp(192.168.12.2:3306)/theca")
	//db.SetMaxOpenConns(10000)
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
					"in.e.a":1,
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
			Sender := q.Confirmed[i].In[0].E.A
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
					err := insertIntoMysql(TxId, Prefix, Hash, Datatype, Title, BlockTimestamp, BlockHeight, Sender)
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
					"in.e.a":1,
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
			Sender := q.Confirmed[i].In[0].E.A
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
					err := insertIntoMysql(TxId, Prefix, Hash, Datatype, Title, 0, 0, Sender)
					if err != nil {
						fmt.Println("INSERT FAILED (unconfirmed): error or duplicated db entry")
					} else {
						fmt.Println("INSERT OK (unconfirmed)==> ", TxId, Prefix, Hash, Datatype, Title)
					}
				}
			}
		}
		// get memo likes
		//getMemoLikes(ScannerBlockHeight)
		time.Sleep(10 * time.Second)
	}

}
