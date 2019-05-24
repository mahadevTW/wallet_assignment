### Command to generate html coverage report 
```
go test ./... -coverpkg=<package_path> -coverprofile=cover.out&&go tool cover -html=cover.out -o coverage.html
```
### Example
```
go test ./... -coverpkg=wallet/app/handler -coverprofile=cover.out&&go tool cover -html=cover.out -o coverage.html
```
###Run Locally?
Install glide from local by putting into src of go path.
[https://github.com/Masterminds/glide](install)

```
glide install
```
export envoronment variables for database connectivity service.env file in parent directory;
Project entirly runs on docker stack, so go parent directory and bring up all dependent containers by doing

```
docker-compose build --no-cache
docker compose up
```
