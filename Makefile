all: migrate

migrate:
	# postgres://$(CRUD_USER):$(CRUD_PASSWORD)@localhost:5432/$(DBNAME)?sslmode=disable
	migrate -database "$(CONNECTION_STRING)" -path migrations up

rollback:
	migrate -database "$(CONNECTION_STRING)" -path migrations down
