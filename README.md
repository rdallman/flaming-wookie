# iClickr Replacement

Senior Design Project at Aubrun University

War Damn Sweet Tea

### TODO

make sure:

get postgresql installed on your machine if you don't:

`% brew install postgres`

get dependencies:

```
% go get github.com/gorilla/mux
% go get github.com/lib/pq
```

the low hanging fruit: 

* drop tables on first go round (design first)
* make html all purrty and make stylesheet for css and template partials like
  nav, imports
* cookies on login (gorilla has nice stuff for this...)
* if no cookie, show login; if cookie, show mgmt stuff in nav
* encapsulate some new .go files

high hanging fruit:

* hook up register (javascript validate w/ regex?)
* make login actually work (and display if otherwise)
* quiz pages (js laden? use cookie to figure out start/not ability?)
* admin panel (view/make classes & quizzes)
* student "register" for class (URI scheme?)
* JSON protocol for answering quiz
* QR code generation / reading

thought splurge: 

select students from class join quiz using qid where id = $ID 
for authentication of student for quiz... so quiz has FK classid?

send Quiz Doc to javascript:
js start quiz:
  every 30 seconds:
    next question
    update state for server

server accept answers?
  obvious:
    map[question]map[student]answer

so http://wook.ie/{quiz} is very complex:
  if we get a cookie and this belongs to them
    give "start"
  else
    ask for student id and PIN

...or just have /teacher/{quiz}... way easier


