DIST?=dist

# Создаёт результирующий дистрибутив фронтенда
.PHONY: frontend
frontend:
	@pushd frontend && \
	test -d ${DIST} || mkdir ${DIST} && \
	yarn build:prod && \
	cp -R -f ${DIST}/frontend/* ../backend/assets && \
	popd

# Создаём результирующий дистрибутив только для Windows x64
.PHONY: build
build: frontend
	@$(MAKE) genfullbindata && \
	pushd backend && \
	test -d ${DIST} || mkdir ${DIST} && \
	set GOARCH=amd64&&set GOOS=windows &&\
	go build -o ${DIST}/server_win_x64.exe cmd/main.go && \
	popd && \
	$(MAKE) genfakebindata && \
	test -d ${DIST} || mkdir ${DIST} && \
	cp -R -f backend/${DIST}/* ./${DIST}/

# Создаёт результирующий дистрибутив под все возможноые версии ОС и архитекрур
.PHONY: all
all: frontend
	@$(MAKE) genfullbindata && \
	pushd backend && \
	test -d ${DIST} || mkdir ${DIST} && \
	set GOARCH=amd64&&set GOOS=windows && \
	go build -o ${DIST}/server_win_x64.exe cmd/main.go && \
	set GOARCH=amd64&&set GOOS=linux && \
	go build -o ${DIST}/server_linux_x64 cmd/main.go && \
	set GOARCH=arm&&set GOOS=linux && \
	go build -o ${DIST}/server_linux_arm cmd/main.go && \
	popd && \
	$(MAKE) genfakebindata && \
	test -d ${DIST} || mkdir ${DIST} && \
	cp -R -f backend/${DIST}/* ./${DIST}/

# Очищает результирующие дистрибутивы всех уровней
.PHONY: clean
clean:
	@rm -f -r ${DIST} && \
	rm -f -r backend/${DIST} && \
	rm -f -r frontend/${DIST}
	
# Генерация полной версии bindata.go со всеми реальными бинарниками
.PHONY: genfullbindata
genfullbindata:
	@pushd backend && \
	go generate cmd/main.go && \
	popd
	
# Генерация фейковой bindata.go для разработки
.PHONY: genfakebindata
genfakebindata:
	@pushd backend && \
	rm bindata/bindata.go && \
	go-bindata -debug -o bindata/bindata.go -pkg bindata -fs assets/... && \
	popd
	