package main

type Theca struct {
	Txid           string
	Link           string
	DataType       string
	Title          string
	BlockTimestamp int64
	BlockHeight    uint32
	Sender         string
	Likes          uint32
	Comments       uint32
	Timestamp      int64
	Score          float64
}

type Comment struct {
	Txid           string
	TxHash         string
	Message        string
	BlockTimestamp int64
	BlockHeight    uint32
	Sender         string
	Timestamp      int64
	Likes          uint32
	Score          float64
}

type SignupPost struct {
	Username     string
	PasswordHash string
	EncryptedPk  string
}

type SignupResponse struct {
	Username    string
	EncryptedPk string
	Signup      bool
}

type LoginPost struct {
	Username     string
	PasswordHash string
}

type LoginResponse struct {
	Username    string
	EncryptedPk string
	Login       bool
}
