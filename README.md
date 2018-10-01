# theca.cash
Container infrastructure

## setup

### docker
install docker (on linux):
```
curl -fsSL get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker <username> 
```

### docker-compose
* https://docs.docker.com/compose/install/#install-compose

### theca.cash
get & build containers
```
git clone https://github.com/Blockgang/theca.cash.git
docker network create --subnet 192.168.11.0/24 --gateway 192.168.11.254 thecanet
docker-compose build
```
start in d1aemon mode
```
docker-compose up -d
```

stop dockers
```
docker-compose down
```

## how to access?

We can use the dns names (from inside the container)

* 192.168.11.2		db
* 192.168.11.3		memcache
* 192.168.11.4		sync
* 192.168.11.5		web

So we can use something like that in our binaries
```
db,err := sql.Open("mysql", root:8drRNG8RWw9FjzeJuavbY6f9@tcp(db:3306)/")

//Connect to our memcache instance
mc := memcache.New("memcache:11211")
```
