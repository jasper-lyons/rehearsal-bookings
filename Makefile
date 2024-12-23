server: cmd/server/main.go
	go build -o bin/server cmd/server/main.go

migrate: cmd/migrate/main.go
	go build -o bin/migrate cmd/migrate/main.go
