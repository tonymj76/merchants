#!/bin/bash
readonly GOPATH=$(go env GOPATH)
readonly MODIFY=github.com/micro/go-micro/api/proto/api.proto=github.com/micro/go-micro/v2/api/proto
if [[ $# -eq 0 ]]; then
  echo "add options"
elif [[ "$1" == 'gen' ]]; then 
  echo "running..."
  $PROTOC -I . --go_out=plugins=micro:. proto/merchant/merchant.proto | echo
  echo "end"
elif [[ "$1" == 'proto' ]]; then #am using this one for merchant
  echo "running..."
  $PROTOC -I . --micro_out=. --go_out=. proto/merchant/merchant.proto | echo
  go-bindata -pkg migrations -ignore bindata -prefix ./datastore/migrations/ -o ./datastore/migrations/bindata.go ./datastore/migrations
  echo "end"
elif [[ "$1" == 'server' ]]; then
  echo "running server..."
  go run main.go | echo
  echo "end"
elif [[ "$1" == 'build' ]]; then
  echo "running server..."
  go build -o merchant-service *.go | echo
  echo "end"
elif [[ "$1" == 'swagger-ui' ]]; then
  echo "running server..."
  docker run --rm -p 8001:8080 swaggerapi/swagger-ui | echo
  echo "end"
elif [[ "$1" == 'swagger-editor' ]]; then
  echo "running server..."
  docker run --rm -p 8001:8080 swaggerapi/swagger-editor | echo
  echo "end"
# elif [[ "$2" == 'gen' ]]; then
#   if [[ "$2" == 'reset' ]]; then
#     echo "Running droping table and running migration ..."
#     buffalo-pop pop reset -e development
#     echo "Finished running migration"
#   elif [[ "$2" == 'migrate' ]]; then
#     echo "running migration ..."
#     buffalo-pop pop migrate up -e development
#     echo "Finished running migration"
#   elif [[ "$2" == 're-mi' ]]; then
#     buffalo-pop pop reset -e development
#     buffalo task db:seed
#   elif [[ "$2" == 'seed' ]]; then
#     buffalo task db:seed
#   else
#     buffalo-pop pop drop -e development
#     buffalo-pop pop create -e development
#     buffalo-pop pop migrate up -e development
#     buffalo task db:seed
#   fi
fi