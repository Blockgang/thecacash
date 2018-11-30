package main

const (
	// API Version
	ApiVersion = "v1"

	// API Handles
	ThecaMainApiPath = "/api/" + ApiVersion + "/theca/all"
	TxApiPath        = "/api/" + ApiVersion + "/theca/{txid:[a-fA-F0-9]{64}}"
	CommentApiPath   = "/api/" + ApiVersion + "/comments/{txid:[a-fA-F0-9]{64}}"
	LoginApiPath     = "/api/" + ApiVersion + "/login"
	SignupApiPAth    = "/api/" + ApiVersion + "/signup"

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
