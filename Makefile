generate-ogen:
	go run github.com/ogen-go/ogen/cmd/ogen@latest --target pkg/rest -package gsn --clean api/rest/collector.swagger.yml

generate-sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate -f configs/sqlc.yaml