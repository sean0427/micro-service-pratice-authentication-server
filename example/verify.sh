#!/bin/bash

read -p "access_token: " access_token

curl -X POST http://localhost:8080/verify -H "Authorization: Bearer ${access_token}" -v
