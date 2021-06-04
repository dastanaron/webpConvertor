.PHONY: all clean build

APP_NAME=webpConvertor

all: prod

prod: clean compile

## default run
run:
	cd src; go run . $(ARGS)

## check race condition
race:
	cd src; go run -race .

## default build
build:	
	cd src; go build -o ../build/${APP_NAME} .

## production build (strip the debugging information)
compile:
	cd src;	GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o ../build/${APP_NAME}_x86 .
	cd src;	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ../build/${APP_NAME} .
	cd src;	GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o ../build/${APP_NAME}_x86.exe .
	cd src;	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ../build/${APP_NAME}.exe .
	cd src; env GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ../build/${APP_NAME}_mac_os .

## clear cache and remove builds
clean:
	go clean
	rm -f build/${APP_NAME}
	rm -f build/${APP_NAME}_x86
	rm -f build/${APP_NAME}
	rm -f build/${APP_NAME}_x86.exe
	rm -f build/${APP_NAME}
