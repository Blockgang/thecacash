package main

const (
	// API Handles
	ThecaMainApiPath = "/api/tx/positions"
	CommentApiPath   = "/api/comments/{txid:[a-fA-F0-9]{64}}"
	LoginApiPath     = "/api/login"
	SignupApiPAth    = "/api/signup"
	TxApiPath        = "/api/tx/{txid:[a-fA-F0-9]{64}}"

	// Fileserver Handles
	FileServerHandlePath = "/var/www/"

	// Webserver
	MuxPort = ":8000"

	// MYSQL
	MysqlInsertUser = `INSERT INTO users (username,password,encrypted_pk) VALUES(?,?,?)`

	MysqlSelectUserLogin = `SELECT encrypted_pk
                            FROM users
                            WHERE username='%s' AND password='%s'`

	MysqlSelectThecaTxList = `SELECT txid,hash,type,title,blocktimestamp,blockheight,sender,likes,comments
                              FROM prefix_0xe901`

	MysqlSelectThecaTx = `SELECT txid,hash,type,title,blocktimestamp,blockheight,sender,UNIX_TIMESTAMP(timestamp)
                          FROM prefix_0xe901
                          WHERE txid='%s'`

	MysqlSelectComments = `SELECT txid,txhash,message,blocktimestamp,blockheight,sender,UNIX_TIMESTAMP(timestamp)
                          FROM prefix_0x6d03
                          WHERE txhash = '%s'`
)
