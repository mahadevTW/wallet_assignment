#Command to generate html coverage report 
go test ./... -coverpkg=<package_path> -coverprofile=cover.out&&go tool cover -html=cover.out -o coverage.html
eg .go test ./... -coverpkg=wallet/app/handler -coverprofile=cover.out&&go tool cover -html=cover.out -o coverage.html