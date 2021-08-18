BINNAME = graphql-postman

clean:
	rm -rf bin

build:
	GO111MODULE=on CGO_ENABLED=0 go build -a -tags netgo -ldflags="-s -w -extldflags '-static'" -o bin/${BINNAME} ./src

run:
	./bin/${BINNAME} -endpoint ${ENDPOINT}

full: clean build run

test: clean build
	./bin/${BINNAME} -endpoint "https://swapi-graphql.netlify.app/.netlify/functions/index"

image:
	docker build -t ${BINNAME}:latest .

docker:
	docker run -it --rm -v ${PWD}:/go/src/result ${BINNAME}:latest

docker-test:
	docker run -it --rm -v ${PWD}:/go/src/result ${BINNAME}:latest /go/src/bin/${BINNAME} -endpoint "https://swapi-graphql.netlify.app/.netlify/functions/index" -output "/go/src/result/api.postman_collection.json"

fmt:
	gofmt -w .

.PHONY: clean build run full test image docker docker-test fmt
