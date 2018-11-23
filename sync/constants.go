package main

const (
	// Blockheight to start sync
	ScannerBlockHeight = uint32(550255)

	// Prefixes
	ThecaPrefix             = "e901"
	MemoLikePrefix          = "6d04"
	MemoCommentPrefix       = "6d03"
	MemoProfileNamePrefix   = "6d01"
	MemoProfileTextPrefix   = "6d05"
	MemoProfilPicturePrefix = "6d0a"
	MemoFollowPrefix        = "6d06"
	MemoUnfollowPrefix      = "6d07"

	// BitdbBlockheightQuery : Query the Heighest Blocknumber
	BlockheightQuery = `{
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

	ThecaQuery = `{
		"v": 3,
		"q": {
			"find": {
				"out.h1": "e901",
				"$or": [{
    			"blk.i": {
    				"$gte" : %d
    			}
  			}, {
  			  "blk": null
  			}]
			},
			"limit":100000
		},
    "r": {
      "f": "[.[] | .tx.h as $tx | .in as $in | .blk as $blk | .out[] | select(.b0.op? and .b0.op == 106) | {link: .s2, type: .s3, title: .s4, txid: $tx, sender: $in[0].e.a, blockheight: (if $blk then $blk.i else null end), blocktimestamp: (if $blk then $blk.t else null end)}]"
    }
	}`

	MemoLikesQuery = `{
		"v": 3,
		"q": {
			"find": {
				"out.h1": "6d04",
				"$or": [{
					"blk.i": {
						"$gte" : %d
					}
				}, {
					"blk": null
				}]
			},
			"limit":100000
		},
		"r": {
			"f": "[.[] | .tx.h as $tx | .in as $in | .blk as $blk | .out[] | select(.b0.op? and .b0.op == 106) | {txhash: .h2, txid: $tx, sender: $in[0].e.a, blockheight: (if $blk then $blk.i else null end), blocktimestamp: (if $blk then $blk.t else null end)}]"
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
