BINARY_NAME=tex_money

linux:
	GOARCH=amd64 GOOS=window go build -ldflags -H=windowsgui -o ./bin/${BINARY_NAME}.exe

window:
	go generate
	go build -ldflags -H=windowsgui -o ./bin/${BINARY_NAME}.exe
window2:
	go generate
	go build -o ./bin/${BINARY_NAME}.exe

window386:
	386.bat

run:
	go run .

clean:
	go clean