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
