APPNAME = structure_fi_coding_challenge
VERSION=`git log -n1 --format="%h"`
BUILD_TIMESTAMP=`date --rfc-3339=seconds`
TESTFLAGS=-v -cover -covermode=atomic -bench=.
TEST_COVERAGE_THRESHOLD=15.0

all: setup build

build:
	go build -tags netgo -ldflags "-w -s -X 'main.AppVersion=${VERSION}' -X 'main.BuildTimestamp=${BUILD_TIMESTAMP}'" -o ${APPNAME} .

build-linux:
	GOOS=linux GOARCH=amd64 go build -tags netgo -ldflags "-w -s -X 'main.AppVersion=${VERSION}' -X 'main.BuildTimestamp=${BUILD_TIMESTAMP}'" -o ${APPNAME}-linux-amd64 .
	shasum -a256 ${APPNAME}-linux-amd64

build-mac:
	GOOS=darwin GOARCH=amd64 go build -tags netgo -ldflags "-w -s -X 'main.AppVersion=${VERSION}' -X 'main.BuildTimestamp=${BUILD_TIMESTAMP}'" -o ${APPNAME}-darwin-amd64
	shasum -a256 ${APPNAME}-darwin-amd64

build-all: build-mac build-linux

setup:
	go mod download
	go install github.com/wadey/gocovmerge@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g api_doc.go

test-only:
	go test ${TESTFLAGS} github.com/ashwanthkumar/structure_fi_coding_challenge/${name}

test:
	go test ${TESTFLAGS} -coverprofile=main.txt github.com/ashwanthkumar/structure_fi_coding_challenge/
	go test ${TESTFLAGS} -coverprofile=store.txt github.com/ashwanthkumar/structure_fi_coding_challenge/store
	go test ${TESTFLAGS} -coverprofile=binance.txt github.com/ashwanthkumar/structure_fi_coding_challenge/binance

test-ci: test
	gocovmerge *.txt > coverage.txt
	@go tool cover -html=coverage.txt -o coverage.html
	@go tool cover -func=coverage.txt | grep "total:" | awk '{print $$3}' | sed -e 's/%//' > cov_total.out
	@bash -c 'COVERAGE=$$(cat cov_total.out);	\
	echo "Current Coverage % is $$COVERAGE, expected is ${TEST_COVERAGE_THRESHOLD}.";	\
	exit $$(echo $$COVERAGE"<${TEST_COVERAGE_THRESHOLD}" | bc -l)'
