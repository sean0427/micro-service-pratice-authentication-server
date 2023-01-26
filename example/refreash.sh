#!/bin/bash

read -p "access_token: " access_token

curl -X POST http://localhost:8080/refresh -H "Authorization: Bearer ${access_token}" -v
