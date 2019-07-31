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
	go test ./test/...

docker-test:
	docker-compose -f docker-compose.test.yml build && \
  docker-compose -f docker-compose.test.yml up -d && \
  docker-compose -f docker-compose.test.yml run app make test-unit && \
  docker-compose -f docker-compose.test.yml run -e DB_HOST=ethereum-server-db -e DB_USER=postgres -e DB_NAME=postgres -e DB_PORT=5432 -e DB_SSL_MODE=disable -e ETH_RPC_ENDPOINT=node1:8501 app make test-integration && \
	docker-compose -f docker-compose.test.yml down
