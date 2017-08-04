#!/bin/bash

echo

curl -X POST --header 'Content-Type: application/json' --header 'Accept: application/json' -d '{
   "key": "string",
   "value" : "value1"
}' 'http://'$1':'$2'/r/put'

echo

curl -X POST --header 'Content-Type: application/json' --header 'Accept: application/json' -d '{
  "key": "string"
}' 'http://'$1':'$2'/r/get'

echo

curl -X POST --header 'Content-Type: application/json' --header 'Accept: application/json' -d '{
  "key": "string"
}' 'http://'$1':'$2'/r/delete'

echo 
