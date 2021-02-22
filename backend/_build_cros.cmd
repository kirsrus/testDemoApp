@echo off
:: Генерация bindata/bindata.go с полным наполнением контента
go generate -v cmd/main.go

mkdir dist

set GOARCH=amd64&&set GOOS=windows
go build -o dist/server_win_x64.exe cmd/main.go

set GOARCH=amd64&&set GOOS=linux
go build -o dist/server_linux_x64 cmd/main.go

set GOARCH=arm&&set GOOS=linux
go build -o dist/server_linux_arm cmd/main.go

::set GOARCH=amd64&&set GOOS=darwin
::go build -o server_darwin_x64 cmd/main.go

:: set GOARCH=amd64&&set GOOS=freebsd
:: go build -o server_freebsd_x64 cmd/main.go

:: set GOARCH=arm&&set GOOS=android
:: go build -o server_android_arm cmd/main.go

:: Удаляем наполненный bindata/bindata.go и заменяем его фейковым
:: для облегчения веса для работы в IDE
rm bindata/bindata.go
go-bindata -debug -o bindata/bindata.go -pkg bindata -fs assets/...
