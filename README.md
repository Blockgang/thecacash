# theca.cash
Container infrastructure

## setup

install docker (on linux):
```
curl -fsSL get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker <username> 
```

install docker-compose
* https://docs.docker.com/compose/install/#install-compose

create data volumes
```
docker volume create thecaweb
docker volume create thecadb
```

get & build containers
```
docker-compose build
```
start in daemon mode
```
docker-compose up -d
```

stop dockers
```
docker-compose down
```

## how to access?

We can use the dns names (from inside the container)

* db
* memcache
* web
* sync 

So we can use something like that in our binaries
```
db,err := sql.Open("mysql", root:8drRNG8RWw9FjzeJuavbY6f9@tcp(db:3306)/")

//Connect to our memcache instance
mc := memcache.New("memcache:11211")
```
