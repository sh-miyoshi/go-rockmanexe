#!/bin/bash

SERVER_ADDR="http://localhost:8880"
URL="$SERVER_ADDR/api/v1"

function test_api() {
  url=$1
  method=$2

  if [ $# = 3 ]; then
    input=$3
    result=`curl --insecure -s -X $method -d "$input" $url`
    # result=`curl --insecure -s -X $method -d "$input" $url | jq .`
    if [ $? != 0 ]; then
      echo "Failed to $method to $url"
      exit 1
    fi
    echo $result
  else
    result=`curl --insecure -s -X $method $url`
    # result=`curl --insecure -s -X $method $url | jq .`
    if [ $? != 0 ]; then
      echo "Failed to $method to $url"
      exit 1
    fi
    echo $result
  fi
}

# Client Create
test_api "$URL/client" POST

# Client Get
test_api "$URL/client" GET

# Client Delete
#test_api "$URL/client/$client_id" DELETE

# Create clients for route test
#test_api "$URL/client" POST
#test_api "$URL/client" POST

# Route Create
#test_api "$URL/route" POST  "{\"clients\": [\"$client1_id\", \"$client2_id\"]}"

# Route Get
test_api "$URL/route" GET

# Route Delete
#test_api "$URL/route/$route_id" DELETE

# Expect client empty after route delete
#clients == []
