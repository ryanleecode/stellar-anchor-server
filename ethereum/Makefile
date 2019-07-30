test-unit:
	go test ./internal/...

test-integration:
	@if [ ! -f .env.test ]; then \
		echo DB_HOST=$$DB_HOST >> .env.test; \
		echo DB_USER=$$DB_USER >> .env.test; \
		echo DB_NAME=$$DB_NAME >> .env.test; \
		echo DB_PORT=$$DB_PORT >> .env.test; \
		echo DB_SSL_MODE=$$DB_SSL_MODE >> .env.test; \
		echo ETH_RPC_ENDPOINT=$$ETH_RPC_ENDPOINT >> .env.test; \
	fi
	cat .env.test
	go test ./test/...