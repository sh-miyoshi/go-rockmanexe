#!/bin/bash

SERVER_ADDR="http://localhost:8880"
URL="$SERVER_ADDR/api/v1"

function test_api() {
  url=$1
  method=$2

  if [ $# = 3 ]; then
    input=$3
    result=`curl --insecure -s -X $method -d "$input" $url | jq .`
    if [ $? != 0 ]; then
      echo "Failed to $method to $url"
      exit 1
    fi
    echo $result
  else
    result=`curl --insecure -s -X $method $url | jq .`
    if [ $? != 0 ]; then
      echo "Failed to $method to $url"
      exit 1
    fi
    echo $result
  fi
}

# Client Create
result=`test_api "$URL/client" POST`
echo "success to create a client"
client_id=`echo $result | jq -r .id`

# Client Get
test_api "$URL/client" GET
echo "success to get clients"

# Client Delete
test_api "$URL/client/$client_id" DELETE
echo "success to delete a client"

# Create clients for route test
result=`test_api "$URL/client" POST`
client1_id=`echo $result | jq -r .id`
result=`test_api "$URL/client" POST`
client2_id=`echo $result | jq -r .id`

# Route Create
test_api "$URL/route" POST  "{\"clients\": [\"$client1_id\", \"$client2_id\"]}"
echo "success to create a route"
route_id=`echo $result | jq -r .id`

# Route Get
test_api "$URL/route" GET
echo "success to get routes"

# Route Delete
test_api "$URL/route/$route_id" DELETE
echo "success to delete a route"

# Expect client empty after route delete
# TODO
