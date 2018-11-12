# THECA.CASH
Media plattform based on Bitcoin Cash (OP_RETURN)

## Setup

### theca.cash
get & build containers
```
git clone https://github.com/Blockgang/theca.cash.git
docker network create --subnet 192.168.12.0/24 --gateway 192.168.12.254 thecanet
docker-compose build
```
start in daemon mode
```
docker-compose up -d
```
start in verbose mode
```
docker-compose up
```

stop dockers
```
docker-compose down
```


## API Access
### GET ###
Get Transaction Infos
```
http://192.168.12.5:8000/api/tx/{txid}
http://192.168.12.5:8000/api/tx/569be470b326e50afbbc739531ea428b5c6977fd900091e3a8faeaf90b85140b
```
Get All Transactions (inkl. like,comment counter + score)
```
http://192.168.12.5:8000/api/tx/positions
```

### POST ###
#### SIGNUP ####
POST-Request:
```
curl -X POST -i 'http://192.168.12.5:8000/api/signup' --data '{"Username":"testuser8","PasswordHash":"105d5b6c13df8c30686b0d75b89d98ada04dc32421fd97acfb77bc81e43f6075","EncryptedPk":"this is the excrypted privatekey"}'
```
Possible Responses:
```
OK:
{"Username":"**username**","EncryptedPk":"**enc_key**","Signup":true}
Failed:
{"Username":"**username**","EncryptedPk":"**enc_key**","Signup":false}
```
#### LOGIN ####
POST-Request:
```
curl -X POST -H 'Content-Type: application/json' -i 'http://192.168.12.5:8000/api/login' --data '{"Username":"testuser8","PasswordHash":"105d5b6c13df8c30686b0d75b89d98ada04dc32421fd97acfb77bc81e43f6075"}'

OK:
{"Username":"**username**","EncryptedPk":"**enc_key**","Login":true}
Failed:
{"Username":"**username**","EncryptedPk":"","Login":false}
```

## Dependencies
#### sync
##### dependencies (for build)
```
 go get -u github.com/go-sql-driver/mysql
 go get github.com/tidwall/gjson
```

#### web
##### dependencies (for build)
```
go get github.com/bradfitz/gomemcache/memcache
go get github.com/gorilla/mux
go get github.com/pmylund/sortutil
go get github.com/junhsieh/goexamples/fieldbinding/fieldbinding
```

## Links?
* https://blockgang.github.io/chaintube
* https://instant.io/
* https://icons8.com/icon/set/error/all
* https://github.com/unwriter/datacash
* https://docs.bitdb.network/docs/quickstart
* https://github.com/webtorrent/webtorrent

# Prefix

Prefix: 0xe901 (Main)
```

# OP_RETURN 0xe901 <magnet/ipfs-hash> <data-type> <titel>

OP_RETURN (PD1)0xe901  (PD2)magnet:?xt=urn:btih:678d1a0744863813bd11e12c473e0a2ab3d07f27 (PD3)mp4 (PD4)DAS IST DER TITEL DES VIEOS

```

Prefix: 0xe902 (Description)
```
# OP_RETURN 0xe902 <hash>|<chunk-nr>|<data>

OP_RETURN (PD1)0xe902 (PD2)magnet:?xt=urn:btih:678d1a0744863813bd11e12c473e0a2ab3d07f27 (PD3)0 (PD4)Erster Teil der Beschreibung

OP_RETURN (PD1)0xe902 (PD2)magnet:?xt=urn:btih:678d1a0744863813bd11e12c473e0a2ab3d07f27 (PD3)1 (PD4)und das ist der 2. Teil der Beschreibung
...
```

Prefix: 0x6d04 (MEMO Like + Tip)
Prefix: 0x6d03 (MEMO Reply)


MEMO-Example:
```
Action 	Prefix 	Values 	Status 	Example
Set name 	0x6d01 	name(217) 	Implemented 	
Post memo 	0x6d02 	message(217) 	Implemented 	
Reply to memo 	0x6d03 	txhash(30), message(184) 	Implemented 	
Like / tip memo 	0x6d04 	txhash(30) 	Implemented 	
Set profile text 	0x6d05 	message(217) 	Implemented 	
Follow user 	0x6d06 	address(35) 	Implemented 	
Unfollow user 	0x6d07 	address(35) 	Implemented 	
Set profile picture 	0x6d0a 	url(217) 	Implemented 	
Repost memo 	0x6d0b 	txhash(30), message(184) 	Planned 	-
Post topic message 	0x6d0c 	topic_name(variable), message(214 - topic length) 	Implemented 	
Topic follow 	0x6d0d 	topic_name(variable) 	Implemented 	
Topic unfollow 	0x6d0e 	topic_name(variable) 	Implemented 	
Create poll 	0x6d10 	poll_type(1), option_count(1), question(209) 	Implemented 	
Add poll option 	0x6d13 	poll_txhash(30), option(184) 	Implemented 	
Poll vote 	0x6d14 	poll_txhash(30), comment(184) 	Implemented 	
Send money 	0x6d24 	message(217) 	Planned
```

## TODO

* webtorrent performance
* webtorrent bug fix ( videos werden nicht immer angezeigt)
