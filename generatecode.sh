#!/bin/bash

#Creation of Server-stub
#swagger-codegen generate -i ./server/api/registry.yaml -l  -o ./server/
java -jar ~/workspace/swagger/swagger-codegen/modules/swagger-codegen-cli/target/swagger-codegen-cli.jar generate \
   -i ./API/registry.yaml \
   -l nodejs-server \
   -o ./server/

#Documentation
#swagger-codegen generate -i ./server/api/registry.yaml -l dynamic-html -o .

#java -jar ~/workspace/swagger/swagger-codegen/modules/swagger-codegen-cli/target/swagger-codegen-cli.jar generate \
#   -i ./API/registry.yaml \
#   -l dynamic-html \
#   -o .

