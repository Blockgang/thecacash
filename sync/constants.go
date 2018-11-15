package constants

const (
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

	BitdbBlockheightQuery = `{
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

	ConfirmedBitdbThecaQuery = `{
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

	UnconfirmedBitdbThecaQuery = `{
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

	MemoLikesQuery = `{
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

	UnconfirmedMemoLikesQuery = `{
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

	MemoCommentsQuery = `{
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
      "limit": 100000
    },
    "r": {
      "f": "[.[] | .tx.h as $tx | .in as $in | .blk as $blk | .out[] | select(.b0.op? and .b0.op == 106) | {txhash: .h2, message: .s3, txid: $tx, sender: $in[0].e.a, blockheight: (if $blk then $blk.i else null end), blocktimestamp: (if $blk then $blk.t else null end)}]"
    }
  }`
)
