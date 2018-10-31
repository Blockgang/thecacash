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
)

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

var db *sql.DB
var q Query
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
			if !exists {
				err := insertIntoMysql(TxId, Hash, Datatype, Title, 0, 0, Sender)
				if err != nil {
					fmt.Println("INSERT FAILED (unconfirmed): error or duplicated db entry")
				} else {
					fmt.Println("INSERT OK (unconfirmed)==> ", TxId, Prefix, Hash, Datatype, Title)
				}
			}
		}
	}
}

func getConfirmed_E901(ScannerBlockHeight uint32, unconfirmedInDb []string) uint32 {
	query := fmt.Sprintf(c_query, ScannerBlockHeight)
	body, err := getBitDbData(query)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(body, &q)
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

			if exists {
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
	body, err := getBitDbData(uc_memoLikesQuery)
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(body, &q)

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
	body, err := getBitDbData(query)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(body, &q)
	fmt.Println(err)

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
			if exists {
				err := updateMysql(Prefix, TxId, BlockTimestamp, BlockHeight)
				if err != nil {
					fmt.Println("CONFIRMED UPDATE FAILED ==> ", err)
				} else {
					fmt.Println("CONFIRMED UPDATE OK ==> ", TxId)
				}
			} else {
				err := insertMemoLikeIntoMysql(TxId, Hash, Sender, BlockTimestamp, BlockHeight)
				if err != nil {
					fmt.Println("CONFIRMED INSERT FAILED/DUPLICATED ==> ", err)
				} else {
					fmt.Println("CONFIRMED INSERT OK ==> ", TxId, " likes ", Hash)
				}
			}
		}
	}
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
	// fmt.Println(url)
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

func main() {
	ScannerBlockHeight = 550255
	ScannerBlockHeight_E901 := ScannerBlockHeight
	// ScannerBlockHeight_D604 := ScannerBlockHeight

	var err error
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

		fmt.Println("E901 Confirmed ScannerHeight: >", ScannerBlockHeight_E901)
		ScannerBlockHeight_E901 = getConfirmed_E901(ScannerBlockHeight_E901, unconfirmedInDb_E901)

		getUnconfirmed_E901(unconfirmedInDb_E901)

		// fmt.Println("MEMO D604 ScannerHeight: >", ScannerBlockHeight_D604)
		// ScannerBlockHeight_D604 = getMemoLikes(ScannerBlockHeight_D604, unconfirmedInDb_6D04)

		getUnconfirmedMemoLikes(unconfirmedInDb_6D04)

		time.Sleep(10 * time.Second)
	}

}
