
#!/bin/bash

#integration testing
#function teacherLogIn {}
#function teacherLogOut {}
#function teacherLogOut {}
#function teacherCreateClass {}
#function teacherEditClass {}
#function teacherAddStudents {}

#handler testing
#function viewQuizzes {}
#function viewClasses {}
#function viewSpecificQuiz {}
#function createQuiz {}
#function changeState {}
#function answerQuestion {}

set -eu

cookies=/tmp/$(basename $0)-$$.txt

curl -X POST \
      -d email=wookie@wookie.com \
          -d password=password \
              -c $cookies http://localhost:8080/logmein


fail=0

#echo "You're about to delete everything in the db!!"



output=$(curl -s -X POST -d "{\"name\":\"myclass\", \"students\": [{\"email\":\"mjs8119@gmail.com\", \"fname\":\"maghen\", \"lname\":\"smith\"}]}" -b $cookies http://localhost:8080/classes)
  cid=$(expr "$output" : '.*"info":\([0-9]*\)')
  #echo $cid

curl -s -X GET -b $cookies http://localhost:8080/classes | grep -q '"success":true' || (echo "expected succes: /classes" && fail=1) 
curl -s -X GET -b $cookies http://localhost:8080/classes/$cid/quiz | grep -q '"success":true' || (echo "expected succes: /class $cid  quiz" && fail=1)
output1=$(curl -s -X POST -d "{\"title\":\"quiz\", \"questions\":[{\"text\":\"Are there many like it?\",\"answers\":[\"Yay\",\"Nay\"],\"correct\":0}],\"grades\":{}}" -b $cookies http://localhost:8080/classes/$cid/quiz)
  qid=$(expr "$output1" : '.*"info":\([0-9]*\)')
  #echo $qid

curl -s -X GET -b $cookies http://localhost:8080/quiz | grep -q '"success":true'|| (echo "expected succes: /quiz" && fail=1) 
#curl -X POST -d "{\"email\":\"mjs8119@aol.com\",\"fname\":\"M\",\"lname\":\"S\"}" -b $cookies http://localhost:8080/classes/152/student | grep -q '"success":true' || echo "expected succes: /add student to class 71"
curl -s -X GET -b $cookies http://localhost:8080/quiz/$qid | grep -q '"success":true' || (echo "expected succes: /quiz 117" && fail=1)
curl -s -X PUT -d "{\"state\":0}" -b $cookies http://localhost:8080/quiz/117/state | grep -q '"success":true' || (echo "expected succes: /quiz 117 state" && fail=1)
curl -s -X PUT -d "{\"id\":\"ojSqAbnPGIEj7fA5\", \"answer\":1}" http://localhost:8080/quiz/117/answer | grep -q '"success":true' || (echo "expected succes: /quiz 117 answer" && fail=1)
curl -s -X PUT -d "{\"state\":-1}" -b $cookies http://localhost:8080/quiz/117/state | grep -q '"success":true' || (echo "expected succes: /quiz 117 state" && fail=1)

echo "DELETE from quiz where qid=$qid" |  psql -h absker.com -U wookie > /dev/null
echo "select * from quiz" |  psql -h absker.com -U wookie 

echo "DELETE from classes where cid=$cid" |  psql -h absker.com -U wookie > /dev/null
echo "select * from classes" |  psql -h absker.com -U wookie

if [ $fail -eq 0 ]; then
  echo "pass" 
fi
