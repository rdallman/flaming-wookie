#!/bin/bash

curl -X POST -d "{\"title\":\"This is my quiz\",\"cid\":11,\"info\": {\"questions\":[{\"text\":\"Are there many like it?\",\"answers\":[\"Yay\",\"Nay\"],\"correct\":0}],\"grades\":{\"9025038951111\":0}}}" -b cookies.txt http://localhost:8080/dashboard/quiz
curl -X POST -d "{\"name\":\"myclass\", \"students\": {\"123456789\":\"Reed\",\"12345\":\"You\"}}" -b cookies.txt http://localhost:8080/dashboard/classes
