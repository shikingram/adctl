BIN_FILE=adctl
RELEASE_LINUX_DIR=__LINUX_ADCTL_AMD64__
RELEASE_LINUX_FILE=linux_adctl_amd64.tar.gz
RELEASE_DARWIN_DIR=__DARWIN_ADCTL_AMD64__
RELEASE_DARWIN_FILE=darwin_adctl_amd64.tar.gz

.PHONY: test run build linux clean darwin
all: clean build

build: 
	go build -o ${BIN_FILE}


run: build
	./${BIN_FILE}

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BIN_FILE}
	mkdir ${RELEASE_LINUX_DIR}
	mv adctl ${RELEASE_LINUX_DIR}
	tar -zcvf ${RELEASE_LINUX_FILE} ${RELEASE_LINUX_DIR}

darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${BIN_FILE}
	mkdir ${RELEASE_DARWIN_DIR}
	mv adctl ${RELEASE_DARWIN_DIR}
	tar -zcvf ${RELEASE_DARWIN_FILE} ${RELEASE_DARWIN_DIR}

clean:
	rm -rf ${RELEASE_LINUX_FILE}
	rm -rf ${RELEASE_LINUX_DIR}
	rm -rf ${RELEASE_DARWIN_FILE}
	rm -rf ${RELEASE_DARWIN_DIR}
	rm -rf ${BIN_FILE}