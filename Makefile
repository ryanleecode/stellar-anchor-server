start:
	go run ./cmd/main.go

test-unit:
	cd api-gateway && make test-unit
	cd authentication && make test-unit
	cd ethereum && make test-unit
	cd middleware && make test-unit
	cd static && make test-unit

test-integration:
	cd ethereum && make test-integration

cover:
	go test ./internal/... -coverprofile=coverage.out

test-coverage:
	go tool cover -html=coverage.out

coveralls:
	go test -v -covermode=count -coverprofile=coverage.out ./internal/...

