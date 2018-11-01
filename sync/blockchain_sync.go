package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tidwall/gjson"
)

type Bitquery struct {
	Unconfirmed []Row `json:"u"`
	Confirmed   []Row `json:"c"`
}

type Row struct {
	TxId           string `json:"txid"`
	Prefix         string `json:"prefix"`
	TxHash         string `json:"txhash"`
	Sender         string `json:"sender"`
	Message        string `json:"message"`
	BlockHeight    uint32 `json:"blockheight"`
	BlockTimestamp uint32 `json:"blocktimestamp"`
}

type Query struct {
	Unconfirmed []Transaction `json:"u"`
	Confirmed   []Transaction `json:"c"`
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

const bitdbBlockheight_query = `{
  "v": 3,
  "q": {
    "db": ["c"],
    "find": { },
    "limit": 1
  },
  "r": {
    "f": "[.[] | .blk | { current_blockheight: .i} ]"
  }
}`

const c_query = `{
	"v": 3,
	"e": { "out.b1": "hex"  },
	"q": {
		"db": ["c"],
		"find": {
			"out.b1": "e901",
			"out.b0": {
				"op": 106
			},
			"blk.i": {
				"$gte" : %d
			}
		},
		"limit":100000,
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

const uc_query = `{
	"v": 3,
	"e": { "out.b1": "hex"  },
	"q": {
		"db": ["u"],
		"find": {
			"out.b1": "e901",
			"out.b0": {
				"op": 106
			}
		},
		"limit":100000,
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

const memoLikesQuery = `{
	"v": 3,
	"e": {
		"out.b1": "hex",
		"out.b2": "hex"
	},
	"q": {
		"db": ["c"],
		"find": {
			"out.b1": "6d04",
			"out.b0": {
				"op": 106
			},
			"blk.i": {
				"$gte" : %d
			}
		},
		"limit":100000,
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

const uc_memoLikesQuery = `{
	"v": 3,
	"e": {
		"out.b1": "hex",
		"out.b2": "hex"
	},
	"q": {
		"db": ["u"],
		"find": {
			"out.b1": "6d04",
			"out.b0": {
				"op": 106
			}
		},
		"limit":100000,
		"project": {
			"out.b1": 1,
			"out.b2": 1,
			"tx.h": 1,
			"in.e.a":1,
			"_id": 0
		}
	}
}`

const memoCommentsQuery = `{
  "v": 3,
  "q": {
    "find": {
			"out.h1": "6d03",
			"$or": [{
  			"blk.i": {
  				"$gte" : %d
  			}
			}, {
			  "blk": null
			}]
		},
    "limit": 10
  },
  "r": {
    "f": "[.[] | .tx.h as $tx | .in as $in | .blk as $blk | .out[] | select(.b0.op? and .b0.op == 106) | {txhash: .h2, message: .s3, txid: $tx, sender: $in[0].e.a, blockheight: (if $blk then $blk.i else null end), blocktimestamp: (if $blk then $blk.t else null end)}]"
  }
}`

var db *sql.DB
var q Query
var bq Bitquery
var ScannerBlockHeight uint32
var LastScannerBlockHeight uint32

func selectUnconfirmedMysql(prefix string) ([]string, error) {
	var uc_txs []string
	query := "SELECT txid FROM prefix_%s WHERE blockheight = 0"
	query = fmt.Sprintf(query, prefix)
	uc_query, err := db.Query(query)
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

func updateMysql(prefix string, TxId string, blocktimestamp uint32, blockheight uint32) error {
	sql_update := "UPDATE prefix_0x%s SET blockheight=?,blocktimestamp=? where txid=?"
	sql_update = fmt.Sprintf(sql_update, prefix)
	update, err := db.Prepare(sql_update)
	defer update.Close()
	_, err = update.Exec(blockheight, blocktimestamp, TxId)
	return err
}

func insertIntoMysql(TxId string, hash string, data_type string, title string, blocktimestamp uint32, blockheight uint32, sender string) error {
	sql_query := "INSERT INTO prefix_0xe901 (txid,hash,type,title,blocktimestamp,blockheight,sender) VALUES(?,?,?,?,?,?,?)"
	insert, err := db.Prepare(sql_query)
	defer insert.Close()
	_, err = insert.Exec(TxId, hash, data_type, title, blocktimestamp, blockheight, sender)
	return err
}

func insertMemoLikeIntoMysql(TxId string, txHash string, Sender string, BlockTimestamp uint32, BlockHeight uint32) error {
	sql_query := "INSERT INTO prefix_0x6d04 (txid,txhash,blocktimestamp,blockheight,sender) VALUES(?,?,?,?,?)"
	insert, err := db.Prepare(sql_query)
	defer insert.Close()
	_, err = insert.Exec(TxId, txHash, BlockTimestamp, BlockHeight, Sender)
	return err
}

func getUnconfirmed_E901(unconfirmedInDb []string) {
	uc_body, err := getBitDbData(uc_query)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(uc_body, &q)

	for i := range q.Unconfirmed {
		Sender := q.Unconfirmed[i].In[0].E.A
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
			if !exists && !isUnconfirmedInDb(TxId) {
				err := insertIntoMysql(TxId, Hash, Datatype, Title, 0, 0, Sender)
				if err != nil {
					fmt.Println("UNCONFIRMED Theca 0x9e01 INSERT FAILED/DUPLICATED ==> ", err)
				} else {
					fmt.Println("UNCONFIRMED Theca 0x9e01 INSERT OK ==> ", TxId)
				}
			}
		}
	}
}

func getConfirmed_E901(ScannerBlockHeight uint32, unconfirmedInDb []string) uint32 {
	query := fmt.Sprintf(c_query, ScannerBlockHeight)
	response, err := getBitDbData(query)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(response, &q)
	if err != nil {
		fmt.Println(err)
	}

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

			if exists && !isUnconfirmedInDb(TxId) {
				err := updateMysql(Prefix, TxId, BlockTimestamp, BlockHeight)
				if err != nil {
					fmt.Println("CONFIRMED Theca 0xe901 UPDATE FAILED ==> ", err)
				} else {
					fmt.Println("CONFIRMED Theca 0xe901 UPDATE OK ==> ", TxId)
				}
			} else {
				err := insertIntoMysql(TxId, Hash, Datatype, Title, BlockTimestamp, BlockHeight, Sender)
				if err != nil {
					fmt.Println("CONFIRMED Theca 0xe901 INSERT FAILED/DUPLICATED ==> ", err)
				} else {
					fmt.Println("CONFIRMED Theca 0xe901 INSERT OK ==>", TxId)
				}
			}
		}
	}
	return ScannerBlockHeight
}

func getUnconfirmedMemoLikes(unconfirmedInDb []string) {
	response, err := getBitDbData(uc_memoLikesQuery)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(response, &q)

	for i := range q.Unconfirmed {
		Sender := q.Unconfirmed[i].In[0].E.A
		TxId := q.Unconfirmed[i].Tx.H
		txOuts := q.Unconfirmed[i].Out
		var Prefix string
		var Hash string
		for a := range txOuts {
			if txOuts[a].B1 == "6d04" {
				Prefix = txOuts[a].B1
				Hash, _ = reverseHexStringBytes(txOuts[a].B2)
			}
		}

		if len(Prefix) != 0 && len(Hash) > 20 {
			exists := false
			for i := range unconfirmedInDb {
				uc_txid := unconfirmedInDb[i]
				if uc_txid == TxId {
					exists = true
				}
			}
			if !exists && !isUnconfirmedInDb(TxId) {
				err := insertMemoLikeIntoMysql(TxId, Hash, Sender, 0, 0)
				if err != nil {
					fmt.Println("UNCONFIRMED Memo 0x6d04 INSERT FAILED ==>", err)
				} else {
					fmt.Println("UNCONFIRMED Memo 0x6d04 INSERT OK ==> ", TxId, " likes ", Hash)
				}
			}
		}
	}
}

func getMemoLikes(ScannerBlockHeight uint32, unconfirmedInDb []string) uint32 {
	query := fmt.Sprintf(memoLikesQuery, ScannerBlockHeight)
	response, err := getBitDbData(query)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(response, &q)

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
		for a := range txOuts {
			if txOuts[a].B1 == "6d04" {
				Prefix = txOuts[a].B1
				Hash, _ = reverseHexStringBytes(txOuts[a].B2)
			}
		}

		if len(Prefix) != 0 && len(Hash) > 20 {
			exists := false
			for i := range unconfirmedInDb {
				uc_txid := unconfirmedInDb[i]
				if uc_txid == TxId {
					exists = true
				}
			}
			if exists && !isUnconfirmedInDb(TxId) {
				err := updateMysql(Prefix, TxId, BlockTimestamp, BlockHeight)
				if err != nil {
					fmt.Println("CONFIRMED Memo 0x6d04 UPDATE FAILED ==> ", err)
				} else {
					fmt.Println("CONFIRMED Memo 0x6d04 UPDATE OK ==> ", TxId)
				}
			} else {
				err := insertMemoLikeIntoMysql(TxId, Hash, Sender, BlockTimestamp, BlockHeight)
				if err != nil {
					fmt.Println("CONFIRMED Memo 0x6d04 INSERT FAILED/DUPLICATED ==> ", err)
				} else {
					fmt.Println("CONFIRMED Memo 0x6d04 INSERT OK ==> ", TxId, " likes ", Hash)
				}
			}
		}
	}
	return ScannerBlockHeight
}

func getMemoComments(ScannerBlockHeight uint32, unconfirmedInDb []string) uint32 {
	query := fmt.Sprintf(memoCommentsQuery, ScannerBlockHeight)
	fmt.Println(query)
	response, err := getBitDbData(query)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(response, &bq)
	fmt.Println(string(response))

	var BlockHeight uint32

	for i := range bq.Confirmed {
		Sender := bq.Confirmed[i].Sender
		TxId := bq.Confirmed[i].TxId
		TxHash, _ := reverseHexStringBytes(bq.Confirmed[i].TxHash)
		Message := bq.Confirmed[i].Message
		BlockTimestamp := bq.Confirmed[i].BlockTimestamp
		BlockHeight = bq.Confirmed[i].BlockHeight
		if BlockHeight > ScannerBlockHeight {
			ScannerBlockHeight = BlockHeight + 1
		}
		fmt.Println(TxId, TxHash, BlockTimestamp, Sender, Message)
	}
	// 	if len(Prefix) != 0 && len(Hash) > 20 {
	// 		exists := false
	// 		for i := range unconfirmedInDb {
	// 			uc_txid := unconfirmedInDb[i]
	// 			if uc_txid == TxId {
	// 				exists = true
	// 			}
	// 		}
	// 		if exists && !isUnconfirmedInDb(TxId) {
	// 			err := updateMysql(Prefix, TxId, BlockTimestamp, BlockHeight)
	// 			if err != nil {
	// 				fmt.Println("CONFIRMED Memo 0x6d04 UPDATE FAILED ==> ", err)
	// 			} else {
	// 				fmt.Println("CONFIRMED Memo 0x6d04 UPDATE OK ==> ", TxId)
	// 			}
	// 		} else {
	// 			err := insertMemoLikeIntoMysql(TxId, Hash, Sender, BlockTimestamp, BlockHeight)
	// 			if err != nil {
	// 				fmt.Println("CONFIRMED Memo 0x6d04 INSERT FAILED/DUPLICATED ==> ", err)
	// 			} else {
	// 				fmt.Println("CONFIRMED Memo 0x6d04 INSERT OK ==> ", TxId, " likes ", Hash)
	// 			}
	// 		}
	// 	}
	// }
	return ScannerBlockHeight
}

func getBitDbData(query string) ([]byte, error) {
	b64_query := base64.StdEncoding.EncodeToString([]byte(query))
	url := "https://bitdb.network/q/" + b64_query
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("key", "qz6qzfpttw44eqzqz8t2k26qxswhff79ng40pp2m44")
	res, _ := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	return body, err
}

var ucInDb []string

func isUnconfirmedInDb(txid string) bool {
	exists := false
	for i := range ucInDb {
		uc_txid := ucInDb[i]
		if uc_txid == txid {
			exists = true
		}
	}
	if exists == false {
		ucInDb = append(ucInDb, txid)
	}
	return exists
}

func reverseHexStringBytes(hexString string) (string, error) {
	hexBytes, err := hex.DecodeString(hexString)
	runes := []byte(hexBytes)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	reversedHexString := hex.EncodeToString(runes)
	return string(reversedHexString), err
}

func getBlockheight() uint64 {
	response, err := getBitDbData(bitdbBlockheight_query)
	if err != nil {
		log.Fatal(err)
	}
	var currentBlockheight uint64

	json := gjson.Get(string(response), "c.#.current_blockheight")
	for _, name := range json.Array() {
		currentBlockheight = name.Uint()
	}
	return currentBlockheight
}

func main() {
	var err error

	var test []string
	getMemoComments(554000, test)
	time.Sleep(10 * time.Second)

	currentBlockheight := getBlockheight()
	fmt.Println("Currentblockheight: ", currentBlockheight)

	ScannerBlockHeight = 550255
	ScannerBlockHeight_E901 := ScannerBlockHeight
	ScannerBlockHeight_D604 := ScannerBlockHeight

	db, err = sql.Open("mysql", "theca:theca123!@tcp(127.0.0.1:3306)/theca")
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(30)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Start ScannerBlockHeight: >", ScannerBlockHeight)

	loop := true
	for loop {
		unconfirmedInDb_E901, err := selectUnconfirmedMysql("0xe901")
		if err != nil {
			fmt.Println(err)
		}
		unconfirmedInDb_6D04, err := selectUnconfirmedMysql("0x6d04")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("THECA 0xE901 Confirmed ScannerHeight: >", ScannerBlockHeight_E901)
		ScannerBlockHeight_E901 = getConfirmed_E901(ScannerBlockHeight_E901, unconfirmedInDb_E901)

		getUnconfirmed_E901(unconfirmedInDb_E901)

		fmt.Println("MEMO Confirmed 0xD604 ScannerHeight: >", ScannerBlockHeight_D604)
		ScannerBlockHeight_D604 = getMemoLikes(ScannerBlockHeight_D604, unconfirmedInDb_6D04)

		getUnconfirmedMemoLikes(unconfirmedInDb_6D04)

		time.Sleep(10 * time.Second)
	}

}
