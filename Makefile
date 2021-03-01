check_swagger:
	which swagger || (GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger: check_swagger
	rm ./docs/swagger.json
	GO111MODULE=off swagger generate spec -o ./docs/swagger.json --scan-models

server-swagger: check_swagger 
	swagger serve -F=swagger swagger.yaml

run:
	@export  BIND_ADDRESS=:9090; \
	export   DB_PBE_MYSQL_USERNAME=root; \
	export   DB_PBE_MYSQL_PASSWORD=mysecurepassword; \
	export   DB_PBE_MYSQL_HOST=127.0.0.1; \
	export   DB_PBE_MYSQL_SCHEMA=pbe; \
	go run main.go