#!/bin/bash

set -eu

cookies=/tmp/$(basename $0)-$$.txt

trap "rm -f $cookies" EXIT

curl -X POST \
  -d email=wookie@wookie.com \
  -d password=password \
  -c $cookies http://localhost:8080/logmein

echo -e "1\n" && curl -X GET -b $cookies http://localhost:8080/quiz
echo -e "2\n" && curl -X GET -b $cookies http://localhost:8080/classes
#curl -X POST -d "{\"name\":\"myclass\", \"students\": {\"123456789\":\"Reed\",\"12345\":\"You\"}}" -b $cookies http://localhost:8080/classes
echo -e "3\n" && curl -X GET -b $cookies http://localhost:8080/classes/11/quiz
#curl -X POST -d "{\"title\":\"This is my quiz\",\"cid\":11,\"info\": {\"questions\":[{\"text\":\"Are there many like it?\",\"answers\":[\"Yay\",\"Nay\"],\"correct\":0}],\"grades\":{\"9025038951111\":0}}}" -b $cookies http://localhost:8080/classes/11/quiz
#curl -X POST -d "{\"name\":\"Maghen\", \"sid\":\"1\"}" -b $cookies http://localhost:8080/classes/11/student
echo -e "4\n" && curl -X GET -b $cookies http://localhost:8080/quiz/71
echo -e "5\n" && curl -X PUT -d "{\"state\":0}" -b $cookies http://localhost:8080/quiz/71/state
echo -e "6\n" && curl -X PUT -d "{\"id\":\"12345\", \"answer\":0}" http://localhost:8080/quiz/71/answer

