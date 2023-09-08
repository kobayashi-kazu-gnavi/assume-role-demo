build:
	GOOS=linux GOARCH=arm64 go1.20.7 build -o bootstrap
	zip bootstrap.zip bootstrap


