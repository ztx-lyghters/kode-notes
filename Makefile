all: clean build run
all-attach: clean compile build run-attach

rebuild: clean build

build: compile
	cp ./schema/default.sql ./docker/kode-db/init.sql
	docker-compose build kode-db
	docker-compose build kode-notes
.PHONY: build

compile: tests
	# CGO needs to be disabled for Alpine containers
	# without libc-compat installed
	CGO_ENABLED=0 GODEBUG=netdns=go go build -v -o docker/kode-notes/kode-notes app/main.go

tests:
	go test -count=1 -race -v ./...

coverage:
	go test -count=1 -race -coverprofile=coverage.out -v ./...
	go tool cover -html=coverage.out
	rm coverage.out

run:
	docker-compose up kode-notes --detach

run-attach:
	docker-compose up kode-notes

stop:
	docker-compose stop

down:
	docker-compose down

attach:
	@printf "%s\n%s\n%s\n%s\n%s\n%s\n\n%s\n" \
	"docker-compose is a weird piece of junk..." \
	"^P^Q doesn't work, so you can't really detach" \
	"from the session after attaching." \
	"BUT you can use CTRL+\ (SIGQUIT) to accomplish that!" \
	"Anyway, just a friendly reminder..." \
	"<PRESS ENTER TO CONTINUE>"

	@read -r enter
	docker-compose up

genmocks:
	mockgen -source=repository/repository.go \
		-destination=mocks/repository.go -package=mocks

clean:
	docker-compose stop || true
	docker-compose rm -f || true
	docker-compose down || true
	docker-compose down -v || true
	rm docker/kode-notes/kode-notes || true
