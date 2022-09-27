BIN_FILE=adctl

.PHONY: test run build
all: check build

build: 
	go build -o ${BIN_FILE}


run: build
	./${BIN_FILE}

check:
	go mod tidy
