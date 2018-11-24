package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tidwall/gjson"
)

var db *sql.DB
var bq Bitquery

func selectUnconfirmedMysql(prefix string) ([]string, error) {
	var uc_txs []string
	query := "SELECT txid FROM prefix_0x%s WHERE blockheight = 0"
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

func insertIntoMysql(TxId string, link string, data_type string, title string, blocktimestamp uint32, blockheight uint32, sender string) error {
	sql_query := "INSERT INTO prefix_0xe901 (txid,hash,type,title,blocktimestamp,blockheight,sender) VALUES(?,?,?,?,?,?,?)"
	insert, err := db.Prepare(sql_query)
	defer insert.Close()
	_, err = insert.Exec(TxId, link, data_type, title, blocktimestamp, blockheight, sender)
	return err
}

func insertMemoLikeIntoMysql(TxId *string, txHash *string, Sender *string, BlockTimestamp *uint32, BlockHeight *uint32) error {
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

func getConfirmed_E901(ScannerBlockHeight uint32, unconfirmedInDb []string) uint32 {
	query := fmt.Sprintf(ThecaQuery, ScannerBlockHeight)
	response, err := getBitDbData(query)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(response, &bq)
	if err != nil {
		fmt.Println(err)
	}

	for i := range bq.Confirmed {
		row := bq.Confirmed[i]

		if row.BlockHeight > ScannerBlockHeight {
			ScannerBlockHeight = row.BlockHeight + 1
		}

		if isUnconfirmedInDb(row.TxId) {
			err := updateMysql(ThecaPrefix, row.TxId, row.BlockTimestamp, row.BlockHeight)
			if err != nil {
				fmt.Printf("CONFIRMED %s UPDATE FAILED ==> %s\ndata:%v\n", ThecaPrefix, err, row)
			} else {
				fmt.Printf("CONFIRMED %s UPDATE OK ==> %s\n", ThecaPrefix, row.TxId)
			}
		} else {
			err := insertIntoMysql(row.TxId, row.Link, row.Type, row.Title, row.BlockTimestamp, row.BlockHeight, row.Sender)
			if err != nil {
				fmt.Printf("CONFIRMED %s INSERT FAILED/DUPLICATED ==> %s\ndata:%v\n", ThecaPrefix, err, row)
			} else {
				fmt.Printf("CONFIRMED %s INSERT OK ==> %s\n", ThecaPrefix, row.TxId)
			}
		}
	}

	for i := range bq.Unconfirmed {
		row := bq.Unconfirmed[i]

		if row.BlockHeight > ScannerBlockHeight {
			ScannerBlockHeight = row.BlockHeight + 1
		}

		if !isUnconfirmedInDb(row.TxId) {
			err := insertIntoMysql(row.TxId, row.Link, row.Type, row.Title, row.BlockTimestamp, row.BlockHeight, row.Sender)
			if err != nil {
				fmt.Printf("UNCONFIRMED %s INSERT FAILED/DUPLICATED ==> %s\ndata:%v\n", ThecaPrefix, err, row)
			} else {
				fmt.Printf("UNCONFIRMED %s INSERT OK ==> %s\n", ThecaPrefix, row.TxId)
			}
		}
	}
	return ScannerBlockHeight
}

func getMemoLikes(ScannerBlockHeight uint32, unconfirmedInDb []string) uint32 {
	query := fmt.Sprintf(MemoLikesQuery, ScannerBlockHeight)
	response, err := getBitDbData(query)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(response, &bq)

	for i := range bq.Confirmed {
		row := bq.Confirmed[i]
		row.TxHash, _ = reverseHexStringBytes(row.TxHash)

		if row.BlockHeight > ScannerBlockHeight {
			ScannerBlockHeight = row.BlockHeight + 1
		}

		if isUnconfirmedInDb(row.TxId) {
			err := updateMysql(MemoLikePrefix, row.TxId, row.BlockTimestamp, row.BlockHeight)
			if err != nil {
				fmt.Printf("CONFIRMED %s UPDATE FAILED ==> %s\ndata:%v\n", MemoLikePrefix, err, row)
			} else {
				fmt.Printf("CONFIRMED %s UPDATE OK ==> %s\n", MemoLikePrefix, row.TxId)
			}
		} else {
			err := insertMemoLikeIntoMysql(&row.TxId, &row.TxHash, &row.Sender, &row.BlockTimestamp, &row.BlockHeight)
			if err != nil {
				fmt.Printf("CONFIRMED %s INSERT FAILED/DUPLICATED ==> %s\ndata:%v\n", MemoLikePrefix, err, row)
			} else {
				fmt.Printf("CONFIRMED %s INSERT OK ==> %s\n", MemoLikePrefix, row.TxId)
			}
		}
	}

	for i := range bq.Unconfirmed {
		row := bq.Unconfirmed[i]

		if !isUnconfirmedInDb(row.TxId) {
			err := insertMemoLikeIntoMysql(&row.TxId, &row.TxHash, &row.Sender, &row.BlockTimestamp, &row.BlockHeight)
			if err != nil {
				fmt.Printf("UNCONFIRMED %s INSERT FAILED/DUPLICATED ==> %s\ndata:%v\n", MemoLikePrefix, err, row)
			} else {
				fmt.Printf("UNCONFIRMED %s INSERT OK ==> %s\n", MemoLikePrefix, row.TxId)
			}
		}
	}

	return ScannerBlockHeight
}

func getMemoComments(ScannerBlockHeight uint32, unconfirmedInDb []string) uint32 {
	query := fmt.Sprintf(MemoCommentsQuery, ScannerBlockHeight)
	response, _ := getBitDbData(query)

	json.Unmarshal(response, &bq)

	for i := range bq.Confirmed {
		row := bq.Confirmed[i]
		row.TxHash, _ = reverseHexStringBytes(row.TxHash)

		if row.BlockHeight > ScannerBlockHeight {
			ScannerBlockHeight = row.BlockHeight + 1
		}

		if isUnconfirmedInDb(row.TxId) {
			err := updateMysql(MemoCommentPrefix, row.TxId, row.BlockTimestamp, row.BlockHeight)
			if err != nil {
				fmt.Printf("CONFIRMED %s UPDATE FAILED ==> %s\ndata:%v\n", MemoCommentPrefix, err, row)
			} else {
				fmt.Printf("CONFIRMED %s UPDATE OK ==> %s answers %s\n", MemoCommentPrefix, row.TxId, row.TxHash)
			}
		} else {
			err := insertMemoCommentIntoMysql(&row.TxId, &row.TxHash, &row.Message, &row.Sender, &row.BlockTimestamp, &row.BlockHeight)
			if err != nil {
				fmt.Printf("CONFIRMED %s INSERT FAILED/DUPLICATED ==> %s\ndata:%v\n", MemoCommentPrefix, err, row)
			} else {
				fmt.Printf("CONFIRMED %s INSERT OK ==> %s answers %s\n", MemoCommentPrefix, row.TxId, row.TxHash)
			}
		}
	}

	for i := range bq.Unconfirmed {
		row := bq.Unconfirmed[i]
		row.TxHash, _ = reverseHexStringBytes(row.TxHash)

		if !isUnconfirmedInDb(row.TxId) {
			err := insertMemoCommentIntoMysql(&row.TxId, &row.TxHash, &row.Message, &row.Sender, &row.BlockTimestamp, &row.BlockHeight)
			if err != nil {
				fmt.Printf("UNCONFIRMED %s INSERT FAILED/DUPLICATED ==> %s\ndata:%v\n", MemoCommentPrefix, err, row)
			} else {
				fmt.Printf("UNCONFIRMED %s INSERT OK ==> %s answers %s\n", MemoCommentPrefix, row.TxId, row.TxHash)
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
	response, err := getBitDbData(BlockheightQuery)
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
	fmt.Printf("Currentblockheight: %d\n", currentBlockheight)

	ScannerBlockHeight_E901 := ScannerBlockHeight
	ScannerBlockHeight_D604 := ScannerBlockHeight
	ScannerBlockHeight_D603 := ScannerBlockHeight

	db, err = sql.Open("mysql", "root:8drRNG8RWw9FjzeJuavbY6f9@tcp("+os.Getenv("DATABASE_HOSTNAME")+":"+os.Getenv("DATABASE_PORT")+")/"+os.Getenv("DATABASE_NAME"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Printf("Start ScannerBlockHeight: > %d\n", ScannerBlockHeight)

	loop := true
	for loop {
		unconfirmedInDb_E901, err := selectUnconfirmedMysql(ThecaPrefix)
		if err != nil {
			fmt.Println(err)
		}
		unconfirmedInDb_6D04, err := selectUnconfirmedMysql(MemoLikePrefix)
		if err != nil {
			fmt.Println(err)
		}
		unconfirmedInDb_6D03, err := selectUnconfirmedMysql(MemoCommentPrefix)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("THECA %s Confirmed ScannerHeight: > %d\n", ThecaPrefix, ScannerBlockHeight_E901)
		ScannerBlockHeight_E901 = getConfirmed_E901(ScannerBlockHeight_E901, unconfirmedInDb_E901)

		fmt.Printf("MEMO %s ScannerHeight: > %d\n", MemoCommentPrefix, ScannerBlockHeight_D603)
		ScannerBlockHeight_D603 = getMemoComments(ScannerBlockHeight_D603, unconfirmedInDb_6D03)

		fmt.Printf("MEMO %s ScannerHeight: > %d\n", MemoLikePrefix, ScannerBlockHeight_D604)
		ScannerBlockHeight_D604 = getMemoLikes(ScannerBlockHeight_D604, unconfirmedInDb_6D04)

		time.Sleep(10 * time.Second)
	}

}
