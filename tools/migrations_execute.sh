cd $SHORTURL_PROJECT_PATH
./cmd/goose/goose  -dir db/migrations postgres "postgres://postgres:123@localhost:5456/shorturl?sslmode=disable" up