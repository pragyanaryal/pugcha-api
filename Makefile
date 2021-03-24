check_install:
	which swagger || GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: check_install
	GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models

migrate:
	cd migrations; go-bindata -pkg migrations -ignore=migrate.go .

conf:
	cd config; go-bindata -pkg config -ignore=config.go -o ./asset.go .

build:
	go build .

all: migrate conf build

tls-local: migrate conf build
	sudo ./pugcha-backend --tls=true --local=true

tls-prod: migrate conf build
	sudo ./pugcha-backend --tls=true --env=prod

nohup: migrate conf build
	nohup sudo ./pugcha-backen --tls=true --env=prod