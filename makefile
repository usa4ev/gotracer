SRV_SRC=./cmd/main.go
SRV_BINARY_NAME=tracerServer
BIN_PATH=./bin
HTTP_SRV_ENDPOINT="0.0.0.0:8080"
MONGO_URI="mongodb://mongo:27017"
RESOURCE_PATH="./config/resources"

 # Build server
build-srv-linux:
	GOARCH=amd64 GOOS=linux go build -o $(BIN_PATH)/${SRV_BINARY_NAME}-linux $(SRV_SRC)

build-srv-darwin:
	GOARCH=amd64 GOOS=darwin go build -o $(BIN_PATH)/${SRV_BINARY_NAME}-darwin $(SRV_SRC)

build-srv-windows:
	GOARCH=amd64 GOOS=windows go build -o $(BIN_PATH)/${SRV_BINARY_NAME}-windows $(SRV_SRC)

# Run server
run-srv-linux: build-srv-linux
	$(BIN_PATH)/${SRV_BINARY_NAME}-linux -resource-path $(RESOURCE_PATH) -mongobd-uri ${MONGO_URI} -http-server-endpoint ${HTTP_SRV_ENDPOINT}

run-srv-darwin: build-srv-darwin
	$(BIN_PATH)/${SRV_BINARY_NAME}-darwin -resource-path $(RESOURCE_PATH) -mongobd-uri ${MONGO_URI} -http-server-endpoint ${HTTP_SRV_ENDPOINT}

run-srv-windows: build-srv-windows
	$(BIN_PATH)/${SRV_BINARY_NAME}-windows -resource-path $(RESOURCE_PATH) -mongobd-uri ${MONGO_URI} -http-server-endpoint ${HTTP_SRV_ENDPOINT}



