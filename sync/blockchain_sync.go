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

func insertMemoCommentIntoMysql(txId *string, txHash *string, message *string, sender *string, blockTimestamp *uint32, blockHeight *uint32) error {
	sql_query := "INSERT INTO prefix_0x6d03 (txid,txhash,message,blocktimestamp,blockheight,sender) VALUES(?,?,?,?,?,?)"
	insert, err := db.Prepare(sql_query)
	if err != nil {
		fmt.Println(err)
	}
	defer insert.Close()
	_, err = insert.Exec(*txId, *txHash, *message, *blockTimestamp, *blockHeight, *sender)
	return err
}

func getUnconfirmed_E901(unconfirmedInDb []string) {
	uc_body, err := getBitDbData(UnconfirmedBitdbThecaQuery)
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
	query := fmt.Sprintf(ConfirmedBitdbThecaQuery, ScannerBlockHeight)
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
	response, err := getBitDbData(UnconfirmedMemoLikesQuery)
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
	query := fmt.Sprintf(MemoLikesQuery, ScannerBlockHeight)
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
	prefix := "6d03"
	query := fmt.Sprintf(MemoCommentsQuery, ScannerBlockHeight)
	response, _ := getBitDbData(query)

	json.Unmarshal(response, &bq)

	for i := range bq.Confirmed {
		exists := false
		row := bq.Confirmed[i]
		row.TxHash, _ = reverseHexStringBytes(row.TxHash)

		if row.BlockHeight > ScannerBlockHeight {
			ScannerBlockHeight = row.BlockHeight + 1
		}

		for i := range unconfirmedInDb {
			unconfirmedTxid := unconfirmedInDb[i]
			if unconfirmedTxid == row.TxId {
				exists = true
			}
		}
		if exists && !isUnconfirmedInDb(row.TxId) {
			err := updateMysql(prefix, row.TxId, row.BlockTimestamp, row.BlockHeight)
			if err != nil {
				fmt.Printf("CONFIRMED Memo 0x%s UPDATE FAILED ==> %s\n", prefix, err)
			} else {
				fmt.Printf("CONFIRMED Memo 0x%s UPDATE OK ==> %s answers %s\n", prefix, row.TxId, row.TxHash)
			}
		} else {
			err := insertMemoCommentIntoMysql(&row.TxId, &row.TxHash, &row.Message, &row.Sender, &row.BlockTimestamp, &row.BlockHeight)
			if err != nil {
				fmt.Printf("CONFIRMED Memo 0x%s INSERT FAILED/DUPLICATED ==> %s\n", prefix, err)
			} else {
				fmt.Printf("CONFIRMED Memo 0x%s INSERT OK ==> %s answers %s\n", prefix, row.TxId, row.TxHash)
			}
		}
	}

	for i := range bq.Unconfirmed {
		exists := false
		row := bq.Unconfirmed[i]
		row.TxHash, _ = reverseHexStringBytes(row.TxHash)
		for i := range unconfirmedInDb {
			unconfirmedTxid := unconfirmedInDb[i]
			if unconfirmedTxid == row.TxId {
				exists = true
			}
		}
		fmt.Println(exists)
		if !exists && !isUnconfirmedInDb(row.TxId) {
			err := insertMemoCommentIntoMysql(&row.TxId, &row.TxHash, &row.Message, &row.Sender, &row.BlockTimestamp, &row.BlockHeight)
			if err != nil {
				fmt.Printf("UNCONFIRMED Memo 0x%s INSERT FAILED/DUPLICATED ==> %s\n", prefix, err)
			} else {
				fmt.Printf("UNCONFIRMED Memo 0x%s INSERT OK ==> %s answers %s\n", prefix, row.TxId, row.TxHash)
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
	response, err := getBitDbData(BitdbBlockheightQuery)
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

	currentBlockheight := getBlockheight()
	fmt.Println("Currentblockheight: ", currentBlockheight)

	fmt.Println(MemoCommentsQuery)

	ScannerBlockHeight = 550255
	ScannerBlockHeight_E901 := ScannerBlockHeight
	ScannerBlockHeight_D604 := ScannerBlockHeight
	ScannerBlockHeight_D603 := ScannerBlockHeight

	db, err = sql.Open("mysql", "root:8drRNG8RWw9FjzeJuavbY6f9@tcp(192.168.12.1:3306)/theca")
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
		unconfirmedInDb_6D03, err := selectUnconfirmedMysql("0x6d03")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("THECA 0xE901 Confirmed ScannerHeight: > ", ScannerBlockHeight_E901)
		ScannerBlockHeight_E901 = getConfirmed_E901(ScannerBlockHeight_E901, unconfirmedInDb_E901)

		getUnconfirmed_E901(unconfirmedInDb_E901)

		fmt.Println("MEMO 0xD603 ScannerHeight: > ", ScannerBlockHeight_D603)
		getMemoComments(ScannerBlockHeight_D603, unconfirmedInDb_6D03)

		fmt.Println("MEMO Confirmed 0xD604 ScannerHeight: > ", ScannerBlockHeight_D604)
		ScannerBlockHeight_D604 = getMemoLikes(ScannerBlockHeight_D604, unconfirmedInDb_6D04)

		getUnconfirmedMemoLikes(unconfirmedInDb_6D04)

		time.Sleep(10 * time.Second)
	}

}
