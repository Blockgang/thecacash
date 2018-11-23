package main

// Bitquery API response
type Bitquery struct {
	Unconfirmed []Row `json:"u"`
	Confirmed   []Row `json:"c"`
}

// Row in Bitquery API response
type Row struct {
	TxId           string `json:"txid"`
	Prefix         string `json:"prefix"`
	Link           string `json:"link"`
	Type           string `json:"type"`
	Title          string `json:"title"`
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
