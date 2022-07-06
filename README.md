# deck-manager

## Solution description
While designing this service, I tried to keep structure as simple as possible to keep code and structure easy-to-understand.

For example, I think following hexagonal architecture would make this code way easier to extend (e.g. adding GRPc or replacing HTTP caller with something else),
but in case of this problem I believe that is overengineering.

To store data I decided to use Redis as the key-value storage. Relational databases surely could be used too. Chosen architecture allows to easily switch from one storage to another. To do that it's necessary to implement class which satisfies [Repository interface](/internal/repository.go) and inject it into the handler in [main.go](/cmd/main.go)

## Tests
### Unit tests
To run unit tests run following command from the root of this repository
```shell
go test ./app/...
```

### Integration tests
To run integration tests run following command from the root of this repository (requires running Redis instance)
```shell
go test ./repositories/...
```

### Black-box tests
To run black-box tests of the _running_ up run following commands from the root of this repository
```shell
python3 test/test_app.py
```

## Building & running
### Configuration
Application's configuration is located at [.env](/.env) file.

### Running with docker-compose
To run the whole app with docker-compose, run the following command from the root of this repo:
```shell
docker-compose up
```
This will create two separate docker containers - one with the developed app and another with the redis instance.
Ports of those containers will be connected to the host ports, which means you can safely access the app from your host computer.

### Running without docker
`Requires: go-1.18, redis`

Simply start redis instance, build the app, and run the app.
```shell
redis-server &
go build -o ./cmd/main.go
./main.go 
```
