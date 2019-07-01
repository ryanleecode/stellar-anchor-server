start:
	go run ./cmd/main.go

test-unit:
	go test ./internal/...

cover:
	go test ./internal/... -coverprofile=coverage.out

test-coverage:
	go tool cover -html=coverage.out

coveralls:
	go test -v -covermode=count -coverprofile=coverage.out ./internal/...

test-e2e:
	cd ./test/simple-git-repo && git checkout master && git checkout --detach
	go test ./test

generate-docs:
	swagger generate spec -b ./cmd -o swagger.json