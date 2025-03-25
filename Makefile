server: cmd/server/main.go
	go build -o bin/server cmd/server/main.go

migrate: cmd/migrate/main.go
	go build -o bin/migrate cmd/migrate/main.go

clear-checkouts: cmd/clear-checkouts/main.go
	go build -o bin/clear-checkouts cmd/clear-checkouts/main.go

mark-abandoned-bookings: cmd/mark-abandoned-bookings/main.go
	go build -o bin/mark-abandoned-bookings cmd/mark-abandoned-bookings/main.go

send-access-codes: cmd/send-access-codes/main.go
	go build -o bin/send-access-codes cmd/send-access-codes/main.go

