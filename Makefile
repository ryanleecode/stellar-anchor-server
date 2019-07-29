start:
	go run ./cmd/main.go

test-unit:
	cd api-gateway && make test-unit
	cd authentication && make test-unit
	cd ethereum && make test-unit
	cd middleware && make test-unit
	cd static && make test-unit

cover:
	go test ./internal/... -coverprofile=coverage.out

test-coverage:
	go tool cover -html=coverage.out

coveralls:
	go test -v -covermode=count -coverprofile=coverage.out ./internal/...

test-e2e:
	cd ./test/simple-git-repo && git checkout master && git checkout --detach
	go test ./test
