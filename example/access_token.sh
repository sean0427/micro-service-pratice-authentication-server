#!/bin/bash

read -p "user: " user
read -p "password: " password

curl -X POST http://localhost:8080/access_token -H 'Content-Type: application/json' -d "{\"name\":\"${user}\",\"password\":\"${password}\"}" -v
