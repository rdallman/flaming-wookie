# Wookie JSON Protocol v1.0

This document will describe the exposed RESTful methods for the Wookie Quiz
service. Methods will be defined by HTTP method, URL, authentication method as
well as input / output format. Each subsection will describe the methods exposed for
certain types of information. But first, a primer on commonalities.

When Input and Output are mentioned in this document, it is in reference to the
Body of the HTTP request.

## JSON Input

The server will expect a JSON object as input. I.e. every top level field should
be inside of a JSON object, for instance:

```json
{
  "name": "me"
}
```

Many of these requests are also expected to be signed with a cookie, obtained
from hitting the login URL with a valid already registered email and password (this
can be done from a browser or simply over HTTP and stored otherwise (see: Login). 

The specific fields that each method is expecting will be defined below.

## JSON Response

Each response from the server will contain a "success" field, which is a boolean
value. The rest of the JSON is based off whether "success" is true or false.

When success is true, there may be a field called "info" that contains the
actual payload data from the server, in JSON form:

```json
{
  "success": true,
  "info": Object or List of Objects
}
```

When success is false, there will be a field called "message" that contains the
error message that caused the request to fail. Of the form:

```json
{
  "success": false,
  "message": string
}
```

## Authentication Methods
## Quiz Session
## Class Methods

### Get Class List

##### URL

GET /classes

##### Authentication

Valid cookie

##### Input

None

##### Success output

```json
{
  "success": true,
  "info":
  [
    { 
      "name": string,
      "students": []map[string]string
    },
    ...
  ]
}
```

Where:

"name": A string representing the name of the class

"students": A list of students containing email, last name, first name, and student id.

### Create Class

##### URL

POST /classes

##### Authentication

Valid cookie

##### Input

```json
{
  "name": string,
  "students": {
    id : name,
    ...
  }
}
```

Where:

"name" is a string representing the class' name

DEPRECATED this will be changing to hashes as unique ID

"students" is a dictionary of student ID numbers to their names (string:string)

##### Success output

```json
{
  "success": true
  "info": cid
}
```

Where:

"info" is the class id in the database.

### Get Class
### Edit Class 
##### URL

POST /classes/{cid}/

Where:

{cid} is a positive integer representing a valid class id

##### Authentication

Valid cookie for owner of given {cid}

##### Input

```json
{
  "name": string,
}
```

Where:

"name" is a string representing the new class name

##### Success Output

```json
{
  "success": true
}
```


### Delete Class
### Add Student

TODO DEPRECATE TO SEND EMAIL VALIDATION AND/OR MOVE TO EDIT CLASS

##### URL

POST /classes/{cid}/student

Where:

{cid} is a positive integer representing a valid class id

##### Authentication

Valid cookie for owner of given {cid}

##### Input

```json
{
  "name": string,
  "sid": string
}
```

Where:

"name" is a string representing the student's name

"sid" is a student id number

##### Success Output

```json
{
  "success": true
  "info": student
}
```

Where:

"info" is the students first name, last name, email, and sid.


### Delete Student
##### URL

DELETE /classes/{cid}/students

Where:

{cid} is a positive integer representing a valid class id

##### Authentication

Valid cookie for owner of given {cid}

##### Input

```json
{
  "sid": string
}
```

Where:

"sid" is a student id number

##### Success Output

```json
{
  "success": true
}
```


### List Quizzes for Class

##### URL

GET /classes/{cid}/quiz

Where:

{cid} is a positive integer representing a valid class id

##### Authentication

Valid cookie for owner of given {cid}

##### Input

None

##### Success Output

```json
{
  "success": true,
  "info":
  [
    {
      "title": string,
      "qid": int,
      "name": string
      "showGrades": 
    },
    ...
  ]
}
```

Where:

"title" is a string representing the title of the quiz

"qid" is a positive int representing the quiz id #

"name" is the name of the class the quiz is for //TODO REDUNDANT, only needed for /quiz

"showGrades" is string of students and their respective grades
## Quiz Methods

### Create Quiz

##### URL

POST /classes/{cid}/quiz

Where:

{cid} is a positive integer representing a valid class id

##### Authentication

Valid cookie for owner of given {cid}

##### Input

```json
{
  "questions": [
    {
      "text": string,
      "answers": [
        string,
        ...
      ],
      "correct": int
    },
    ...
  ],
  "grades": {
    string : int,
    ...
  }
}
```

Where:

"questions": a list of questions with fields:

  "text": a string representing the question.

  "answers": a list of strings representing possible answers.

  "correct": an int representing the index of the correct answer in above answer list

__DEPRACATED__ below will be changing to something else (hashes, same idea)

"grades": an object mapping student ID #'s (string) to their 0-100 grade percentage (int)

Note: "grades" will most likely be sent in as an empty object upon creation,
this is fine (currently they'll get overwritten later anyway).

##### Success Output

```json
{
  "success": true
  "info": int
}
```

Where:

"info" is the quiz id in the database.

### List Quizzes [for teacher]

##### URL

GET /quiz

Where:

{cid} is a positive integer representing a valid class id

##### Authentication

Valid cookie

##### Input

None

##### Success Output

```json
{
  "success": true,
  "info":
  [
    {
      "title": string,
      "qid": int,
      "name": string
    },
    ...
  ]
}
```

Where:

"title" is a string representing the title of the quiz

"qid" is a positive int representing the quiz id #

"name" is the name of the class the quiz is for //TODO REDUNDANT, only needed for /quiz

### Get Quiz

##### URL

GET /quiz/{qid}

{qid} is a positive integer representing a valid quiz id

##### Authentication

Valid cookie for owner of class which quiz belongs to

##### Input

None

##### Success Output

```json
{
  "questions": [
    {
      "text": string,
      "answers": [
        string,
        ...
      ],
      "correct": int
    },
    ...
  ],
  "grades": {
    string : int,
    ...
  }
}
```

Where:

"questions": a list of questions with fields:

  "text": a string representing the question.

  "answers": a list of strings representing possible answers.

  "correct": an int representing the index of the correct answer in above answer list

__DEPRACATED__ below will be changing to something else (hashes, same idea)

"grades": an object mapping student ID #'s (string) to their 0-100 grade percentage (int)

### Delete Quiz

##### URL

DELETE /quiz/{qid}

Where:

{qid} is a positive integer representing a valid quiz id

##### Authentication

Valid cookie for owner of class which quiz belongs to

##### Input

None

##### Success Output

```json
{
  "success": true
}
```

## Quiz Session API

### Change State (incl. Start)

Teacher's interactions with the quiz session

##### URL

PUT /quiz/{qid}/state

Where:

{qid} is a positive integer representing a valid quiz id

##### Authentication

Valid cookie for owner of class which quiz belongs to

##### Input

```json
{
  "state": int
}
```

Where:

"state" is an integer representing the current 'state' of the quiz. To be sent
in by the teacher to tell the server the current state.

Send `{"state": 0}` to start the quiz. Always first request.

Send `{"state": 1..n}` to go to the next question; where n is the number of questions.

Send `{"state": -1}` to end the quiz and compute grades. To be sent after last question.

##### Success Output

```json
{
  "success": true
}
```

### Send Answer

Students' interactions with the quiz session. The server knows which question
the session is on, so the student's answer will be routed accordingly without
them having any knowledge of the current question. 

##### URL

PUT /quiz/{qid}/answer

##### Authentication

Valid user id for class which quiz belongs to in input

##### Input

```json
{
  "id": string,
  "answer": int
}
```

Where:

"id" is a string representing the user's unique ID for the class which this
particular class belongs to.

"answer" is an integer that represents the index of the selected answer

##### Success Output

```json
{
  "success": true
}
```


### METHOD
##### URL
##### Authentication
##### Input
##### Success Output


/{cid}/quiz

map[cid]

{"answer": 1, "id", 3123123123123}
{"state": 0, "qid": 12} COOKIE
{"success":true/false} FOR STATE/ANSWER
CONNECT
DISCONNECT
