#!/bin/bash

SERVER_ADDR="http://localhost:3000"
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

# Client auth
test_api "$URL/client/auth" POST "{\"client_id\": \"tester1\", \"client_key\": \"testtest\"}"
status=`echo $result | jq .status`
if [ "$status" != "null" ]; then
  echo "failed to authenticate a client"
  exit 1
fi

echo "success to authenticate a client"
session_id=`echo $result | jq -r .session_id`

# Session get
test_api "$URL/session/$session_id" GET
status=`echo $result | jq .status`
if [ "$status" != "null" ]; then
  echo "failed to get a session"
  exit 1
fi

echo "success to get a session"
