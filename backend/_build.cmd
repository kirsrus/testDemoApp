@echo off
:: Генерация bindata/bindata.go с полным наполнением контента
go generate -v cmd/main.go

mkdir dist
go build -o dist/server_win_x64.exe cmd/main.go

:: Удаляем наполненный bindata/bindata.go и заменяем его фейковым
:: для облегчения веса для работы в IDE
rm bindata/bindata.go
go-bindata -debug -o bindata/bindata.go -pkg bindata -fs assets/...
