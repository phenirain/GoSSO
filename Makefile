.PHONY: swagger_gen

swagger_gen:
	swag init -g cmd/sso/main.go -o docs --parseDependency --parseInternal --dir ./internal