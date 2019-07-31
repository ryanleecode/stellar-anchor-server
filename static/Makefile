test-unit:
	go test ./internal/...

docker-test:
	docker-compose -f docker-compose.test.yml build && \
	docker-compose -f docker-compose.test.yml up -d && \
	docker-compose -f docker-compose.test.yml run app make test-unit && \
	docker-compose -f docker-compose.test.yml down
	docker-compose -f docker-compose.prod.yml build
