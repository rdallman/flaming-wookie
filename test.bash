#!/bin/bash

set -eu

cookies=/tmp/$(basename $0)-$$.txt

curl -X POST \
      -d email=wookie@wookie.com \
          -d password=password \
              -c $cookies http://24.178.89.28:8080/logmein

fail=0

#Create Class with 1 student
output=$(curl -s -X POST -d "{\"name\":\"myclass\", \"students\": [{\"email\":\"mjs8119@gmail.com\", \"fname\":\"maghen\", \"lname\":\"smith\"}]}" -b $cookies http://24.178.89.28:8080/classes)
  cid=$(expr "$output" : '.*"info":\([0-9]*\)')
  #echo $cid

#GET /classes  
curl -s -X GET -b $cookies http://24.178.89.28:8080/classes | (grep -q '"success":false' && fail=1 echo 'expected success: GET classes') 

#GET 'myclass'
curl -s -X GET -b $cookies http://24.178.89.28:8080/classes/$cid | (grep -q '"success":false' && fail=1 echo 'expected success: get class $cid')

#Update 'myclass' name to 'updated'
curl -s -X POST -d "{\"name\":\"updated\"}" -b $cookies http://24.178.89.28:8080/classes/$cid | (grep -q '"success":false' && fail=1 echo 'expected success: update class')

#Create quiz for 'updated'
output1=$(curl -s -X POST -d "{\"title\":\"myquiz\", \"questions\":[{\"text\":\"Are there many like it?\",\"answers\":[\"Yay\",\"Nay\"],\"correct\":0}],\"grades\":{}}" -b $cookies http://24.178.89.28:8080/classes/$cid/quiz)
  qid=$(expr "$output1" : '.*"info":\([0-9]*\)')
  #echo $qid

#Add student to class 'updated'
outputAddStudent=$(curl -s -X POST -d "{\"cid\":158, \"email\":\"email@email.com\",\"fname\":\"M\",\"lname\":\"S\"}" -b $cookies http://24.178.89.28:8080/classes/$cid/student)
sid=$(expr "$outputAddStudent" : '.*"sid":"\(.\{16\}\)')
  #echo $sid
#Update student email
curl -s -X PUT -d "{\"sid\":\"$sid\", \"fname\":\"fnameUPDATE\", \"email\":\"emailUPDATE@email.com\"}" -b $cookies http://24.178.89.28:8080/classes/$cid/student | (grep -q '"success":false' && fail=1 echo 'expected success: update student')
#Delete student
curl -s -X DELETE -d "{\"sid\":\"$sid\"}" -b $cookies http://24.178.89.28:8080/classes/$cid/student | (grep -q '"success":false' && fail=1  echo 'expected success: delete student')

#GET /quiz handleQuizList
curl -s -X GET -b $cookies http://24.178.89.28:8080/quiz | (  grep -q '"success":false'  &&  fail=1 echo 'expected success: list quizzes' )
#GET /quiz/$qid handleQuizGet
curl -s -X GET -b $cookies http://24.178.89.28:8080/quiz/$qid | (grep -q '"success":false' && fail=1 echo 'expected success: get quiz $qid')
##DELETE /quiz/$qid handleQuizDelete
curl -s -X DELETE -b $cookies http://24.178.89.28:8080/quiz/$qid | (grep -q '"success":false' && fail=1 echo 'delete quiz')



#TODO fix these with new frontend
#curl -s -X PUT -d "{\"state\":0}" -b $cookies http://localhost:8080/dashboard/#/quiz/117 | grep -q '"success":true' || (echo "expected success: /quiz 117 state"; fail=1)
#curl -s -X PUT -d "{\"id\":\"ojSqAbnPGIEj7fA5\", \"answer\":1}" http://localhost:8080/quiz/117/answer | grep -q '"success":true' || (echo "expected success: /quiz 117 answer" && fail=1)
#curl -s -X PUT -d "{\"state\":-1}" -b $cookies http://localhost:8080/quiz/117/state | grep -q '"success":true' || (echo "expected success: /quiz 117 state" && fail=1)

#Clean up db
curl -s -X DELETE -b $cookies http://24.178.89.28:8080/classes/$cid | (grep -q '"success":false' && fail=1 echo 'expected success: delete class $cid')
#echo "select * from classes" |  psql -h absker.com -U wooki

if [ "$fail"	-eq 0 ]; then
	echo "pass"
fi
	