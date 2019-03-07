DATABASE_URL=oauth2_development

test:
	go test -race ./...

run_client:
	go run main.go test_client

run_development:
	go run main.go server

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html coverage.out

createdb:
	createdb $(DATABASE_URL)

dropdb:
	dropdb $(DATABASE_URL)
