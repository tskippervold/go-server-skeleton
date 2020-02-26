# A Golang server skeleton 🔨.
A simple skeleton for a server written in Go. The project structure tries to follow reccomendations and standards explained in: https://github.com/golang-standards/project-layout.

The project uses Go Modules (https://github.com/golang/go/wiki/Modules) for dependencies.

To run the server: `go run cmd/main.go`

## Localhost, dev, prod, etc.
Create `.yml` files on the `/configs` directory.
You can duplicate the existing `local.yml`.

Run the server with `-config={filename}` argument.
During development, you can do `go run cmd/main.go -config=local.yml`.

### TODO's
 * Handle database connection loss gracefully.
 * Support multiple databases.
 * Have package for JWT token authorization and authentication.
 * Dockerize the project with simple deployment stages.
 * (Nice to have) Add request timing middleware logging duration of request/responses.
 