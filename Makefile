server: cmd/server/main.go
	go build -o bin/server cmd/server/main.go

migrate: cmd/migrate/main.go
	go build -o bin/migrate cmd/migrate/main.go

clear-checkouts: cmd/clear-checkouts/main.go
	go build -o bin/clear-checkouts cmd/clear-checkouts/main.go
