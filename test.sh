#!/bin/bash

curl -X POST -d "{\"title\":\"This is my quiz\",\"cid\":11,\"info\": {\"questions\":[{\"text\":\"Are there many like it?\",\"answers\":[\"Yay\",\"Nay\"],\"correct\":0}],\"grades\":{\"9025038951111\":0}}}" -b cookies.txt http://localhost:8080/dashboard/classes/11/quiz
curl -X POST -d "{\"name\":\"myclass\", \"students\": {\"123456789\":\"Reed\",\"12345\":\"You\"}}" -b cookies.txt http://localhost:8080/dashboard/classes
#curl -X POST -d "{\"name\":\"Maghen\", \"sid\":\"1\"}" -b cookies.txt http://localhost:8080/dashboard/classes/11/students
curl -X PUT -d "{\"state\":0}" http://localhost:8080/quiz/46/state
curl -X PUT -d "{\"id\":\"123456\", \"answer\":0}" http://localhost:8080/quiz/46/answer
