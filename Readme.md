# A Golang server skeleton 🔨.
A simple skeleton for a server written in Go. The project structure tries to follow reccomendations and standards explained in: https://github.com/golang-standards/project-layout.

The project uses Go Modules (https://github.com/golang/go/wiki/Modules) for dependencies.

To run the server: `go run cmd/main.go`

## Localhost, dev, prod, etc.
Create `.yml` files on the `/configs` directory.
You can duplicate the existing `local.yml`.

Run the server with `-config={path_to_config.yml}` argument.
During development, you can do `go run cmd/main.go -config=./config/local.yml`.

## Database migrations
`migrate create -ext sql -dir internal/db/migrations -seq create_<tablename>_table`

_https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md_

## Key generation
Generate RSA keys for signing and verifying JWT tokens:
 * Generate private key: `openssl genrsa -out private.pem 2048`
 * Extract public key: `openssl rsa -in private.pem -pubout > public.pem`

_https://developers.yubico.com/PIV/Guides/Generating_keys_using_OpenSSL.html_

### TODO's
 * Handle database connection loss gracefully.
 * Have package for JWT token authorization and authentication.
   - Support access and refresh tokens. Access token should be stateless, but refresh tokens should be stored in some database as "valid token". Whenever a refresh token operation is submitted its checked against this storage. In case of DoS attack, all refresh token operations will fail.
 * Dockerize the project with simple deployment stages.
