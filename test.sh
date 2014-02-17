#!/bin/bash

curl -X GET -b cookies.txt http://localhost:8080/quiz
echo 
curl -X GET -b cookies.txt http://localhost:8080/classes
echo 
curl -X POST -d "{\"name\":\"myclass\", \"students\": {\"123456789\":\"Reed\",\"12345\":\"You\"}}" -b cookies.txt http://localhost:8080/classes
echo 
curl -X GET -b cookies.txt http://localhost:8080/classes/11/quiz
echo 
curl -X POST -d "{\"title\":\"This is my quiz\",\"cid\":11,\"info\": {\"questions\":[{\"text\":\"Are there many like it?\",\"answers\":[\"Yay\",\"Nay\"],\"correct\":0}],\"grades\":{\"9025038951111\":0}}}" -b cookies.txt http://localhost:8080/classes/11/quiz
echo 
curl -X POST -d "{\"name\":\"Maghen\", \"sid\":\"1\"}" -b cookies.txt http://localhost:8080/classes/11/student
echo 
curl -X GET -b cookies.txt http://localhost:8080/quiz/75
echo 
curl -X PUT -d "{\"state\":0}" -b cookies.txt http://localhost:8080/quiz/71/state
echo 
curl -X PUT -d "{\"id\":\"12345\", \"answer\":0}" http://localhost:8080/quiz/71/answer
