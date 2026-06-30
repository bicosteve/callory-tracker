run:
	nodemon \
		--watch cmd/web \
		--ext go, json, html, tmpl \
		--ignore '**/*_test.go' \
		--exec "go run ./cmd/web" \
		--signal SIGTERM

test:
	go test -v ./...
