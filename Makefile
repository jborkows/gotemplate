.PHONY: run tests tests-json failed-tests create_test_project

run:
	@air -c ./config/air.toml 
build:
	@go build -o bin/gotemplate cmd/main.go
tests: 
	@echo "Running tests..."
	@go test ./... -v -race -shuffle=on 
tests-json: 
	@echo "Running tests..."
	@go test ./... -v -race -shuffle=on -json 
failed-tests: 
	@echo "Running tests..."
	@go test ./... -v -race -shuffle=on -json | jq '.|select(.Action=="fail" and .Test!=null)'

