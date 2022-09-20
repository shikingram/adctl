BIN_FILE=acc

.PHONY: test run build
all: check build

build: 
	@go build -o "${BIN_FILE}" 

run:
	./"${BIN_FILE}"

check:
	@go mod tidy
