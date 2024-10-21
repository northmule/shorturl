cd ..
go test ./... -coverprofile test_cover.out && go tool cover -func test_cover.out