
GOPATH:=$(shell go env GOPATH)
MODIFY=Mgithub.com/micro/go-micro/api/proto/api.proto=github.com/micro/go-micro/v2/api/proto
POSTGRESQL_URL=postgres://postgres:password1@localhost:5432/ocm_merchant?sslmode=disable

.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=${MODIFY}:. --go_out=${MODIFY}:. proto/merchant/merchant.proto
	go-bindata -pkg migrations -ignore bindata -prefix ./datastore/migrations/ -o ./datastore/migrations/bindata.go ./datastore/migrations
.PHONY: build
build: proto

	go build -o merchant-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t merchant-service:latest

bindata:
	go-bindata -pkg migrations -ignore bindata -prefix ./datastore/migrations/ -o ./datastore/migrations/bindata.go ./datastore/migrations

table:
	migrate create -ext sql -dir ./datastore/migrations -seq create_terminals

createdb:
	psql -h localhost -U postgres -w -c "ocm_merchant"

down:
	migrate -database ${POSTGRESQL_URL} -path datastore/migrations down

install:
	go get \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/jteeuwen/go-bindata/... \
		github.com/golang/mock/mockgen

merchantID:
	micro call service.merchant MerchantService.GetMerchantByID '{"id": 1}'
merchant:
	micro call service.merchant MerchantService.GetMerchants
create_merchant:
	micro call service.merchant MerchantService.CreateMerchant '{"merchant": {"number_of_product":343,"email" :"merchant1@gmail.com","role" : 2,"user_id": 234,"phone": {"type" : 1,"number": "04328844423"},"number_of_outlet": 3453,"business_name": "mercshant-store"},"password": "new passworad"}'
create_outlet:
	micro call service.merchant MerchantService.CreateMerchantOutlet '{"merchant_id": 1, "phone": {"type" : 2,"number": "04032884423"}, "latitude": 3435, "longitude": 5645, "position": 1, "city_id": 34, "country_id":3, "email": "outlet@gmail.com", "user_id": 3432, "address":"some address", "available":true}'

