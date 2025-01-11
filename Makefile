all: migrate

run:
	@if [ -z "$(CONFIG_FLAG)" ]; then \
		echo 'run without CONFIG_FLAG flag'; \
		go run ./cmd/GrpcSSO/main.go; \
	else \
		echo 'run with CONFIG_FLAG flag'; \
		go run ./cmd/GrpcSSO/main.go --config=$(CONFIG_FLAG); \
	fi

handle-migrate:
	go run ./cmd/migrator/main.go --config=./config/local.yaml --migrations=./migrations

migrate:
	# postgres://$(CRUD_USER):$(CRUD_PASSWORD)@localhost:5432/$(DBNAME)?sslmode=disable
	@migrate -database "$(CONNECTION_STRING)" -path migrations up

rollback:
	@migrate -database "$(CONNECTION_STRING)" -path migrations down
