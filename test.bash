
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


curl -X GET -b $cookies http://localhost:8080/quiz | grep -q '"success":true' || echo "expected succes: /quiz"
curl -X GET -b $cookies http://localhost:8080/classes | grep -q '"success":true' || echo "expected succes: /classes"
curl -X POST -d "{\"name\":\"myclass\", \"students\": [{\"email\":\"mjs8119@gmail.com\", \"fname\":\"maghen\", \"lname\":\"smith\"}]}" -b $cookies http://localhost:8080/classes | grep -q '"success":true' || echo "expected succes: /classes"
curl -X GET -b $cookies http://localhost:8080/classes/152/quiz \
    | grep -q '"success":true' || echo "expected succes: /class 71 quiz"
curl -X POST -d "{\"title\":\"quiz\", \"questions\":[{\"text\":\"Are there many like it?\",\"answers\":[\"Yay\",\"Nay\"],\"correct\":0}],\"grades\":{}}" -b $cookies http://localhost:8080/classes/152/quiz \
      | grep -q '"success":true' || echo "expected succes: post quiz"
#curl -X POST -d "{\"email\":\"mjs8119@aol.com\",\"fname\":\"M\",\"lname\":\"S\"}" -b $cookies http://localhost:8080/classes/152/student | grep -q '"success":true' || echo "expected succes: /add student to class 71"
curl -X GET -b $cookies http://localhost:8080/quiz/110 | grep -q '"success":true' || echo "expected succes: /quiz 90"
curl -X PUT -d "{\"state\":0}" -b $cookies http://localhost:8080/quiz/110/state | grep -q '"success":true' || echo "expected succes: /quiz 91 state"
curl -X PUT -d "{\"id\":\"0H2+49Ld/ty5ZQE5\", \"answer\":1}" http://localhost:8080/quiz/110/answer | grep -q '"success":true' || echo "expected succes: /quiz 91 answer"
curl -X PUT -d "{\"state\":-1}" -b $cookies http://localhost:8080/quiz/110/state | grep -q '"success":true' || echo "expected succes: /quiz 91 state"

