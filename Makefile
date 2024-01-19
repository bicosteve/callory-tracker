hello:
	echo "Hello test"


run:
	nodemon --exec go run ./cmd/web/*.go --signal SIGTERM 


test: 
	go test -v ./...