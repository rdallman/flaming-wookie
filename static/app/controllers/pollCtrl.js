
angular.module('dashboardApp').controller('PollController', function (sessionService, pollService, classService, $scope, $http, $route, $routeParams, $location, flash) {

  // used for traversing poll in html
  $scope.current = -1;
  $scope.qid = -1;
  $scope.cid;
  $scope.class;
  $scope.className;
  $scope.results;

  $scope.classes = [];

  // main model for poll
  $scope.poll = {
    title: "",
    questions: [],
    grades: {}
  };

  $scope.polllist = []

  if ($routeParams.qid !== undefined) {	
    pollService.getPoll($routeParams.qid).success(function(data) {
      if (data !== undefined) {
        $scope.poll = data["info"]
        $scope.qid = $routeParams.qid;
        pollService.getResults($routeParams.qid).success(function(data) {
          $scope.results = data["info"];
        });
      }
    }).error(function(data) {
      // handle error
    });
  }

  if ($routeParams.cid !== undefined) {
    classService.getClass($routeParams.cid).success(function(data){
      if (data !== undefined) {
        $scope.className = data["info"]["name"];
        $scope.class = data["info"];
        $scope.cid = $routeParams.cid;
      }
    });
  }

  $scope.addQuestion = function() {
    if ($scope.questionform.$valid) {
      $scope.poll.questions.push({text: $scope.question.text, answers: []});
      $scope.question.text = "";
    }
  }

  $scope.addAnswer = function(question, text) {
    question.answers.push(text);
    //$scope.newAnswers[question] = "";

  }

  $scope.setCorrectAnswer = function(question, ansIndex) {
    question.correct = ansIndex;
  }

  $scope.removeQuestion = function(question) {
    $scope.poll.questions.splice($scope.poll.questions.indexOf(question), 1);
  }

  $scope.removeAnswer = function(question, answer) {
    question.answers.splice(question.answers.indexOf(answer), 1);
  }

  $scope.postPoll = function() {
    if ($scope.pollform.$valid) {
      pollService.createPoll($scope.cid, $scope.poll).success(function(data) {
        $location.path('/classes/' + $scope.cid);
      });
      
      //document.getElementById("flash").setAttribute("class", "alert alert-success");
      //flash("You created a poll!");
    }
  }

  /*
   * SESSION STUFF
   * */

  if ($location.$$path.match(/(\/poll\/[0-9]+)$/)) {
    // open dat socket
    sessionService.startSession($routeParams.qid);
  }

  // start the quiz, send state of 0 to the server
  $scope.startPoll = function() {
    // post to server we're starting the quiz
    /*
     *$http({
     *  method: 'PUT',
     *  url: '/quiz/' + $scope.id + '/state',
     *  data: angular.toJson({state: 0}),
     *  headers: {'Content-Type': 'application/json'}
     *});
     */
    $scope.current = 0;
    sessionService.changeState($scope.current);
  }

  $scope.nextQuestion = function() {
    // post to server that we're moving on...
    /*
     *$http({
     *  method: 'PUT',
     *  url: '/quiz/' + $scope.id + '/state',
     *  data: angular.toJson({state: $scope.current}),
     *  headers: {'Content-Type': 'application/json'}
     *});
     */
    $scope.current++;
    sessionService.changeState($scope.current);
  }

  $scope.endPoll = function() {
    /*
     *$http({
     *  method: 'PUT',
     *  url: '/quiz/' + $scope.id + '/state',
     *  data: angular.toJson({state: -1}),
     *  headers: {'Content-Type': 'application/json'}
     *});
     */
    sessionService.endSession();
    // redirect to main
    $location.path('/main');
  }
});
