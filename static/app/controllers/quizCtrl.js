// quiz
//var quizApp = angular.module('quizControllers', ['ngRoute']);

angular.module('dashboardApp').controller('QuizController', function (quizSessionService, quizService, classService, $scope, $http, $route, $routeParams, $location) {

  // used for traversing quiz in html
  $scope.current = -1;
  $scope.id = -1;
  $scope.classId;
  $scope.className;

  $scope.classes = [];

  // main model for quiz
  $scope.quiz = {
    title: "",
    questions: [],
    grades: {}
  };

  $scope.quizlist = []

  if ($routeParams.id !== undefined) {	
    quizService.getQuiz($routeParams.id).success(function(data) {
      if (data !== undefined) {
        $scope.quiz = data["info"]
        $scope.id = $routeParams.id;
      }
    }).error(function(data) {
      // handle error
    });
  }

  if ($location.$$path === "/quizzes") {
    $http({
      method: 'GET',
      url: '/quiz'
    }).success(function(data) {
      $scope.quizlist = data["info"]
    }).error(function(data) {

    });
  }

  if ($location.$$path.match(/(\/classes\/[0-9]+\/new-quiz)/)) {
    // creating a quiz
    $scope.classId = $routeParams.cid;
    classService.getClass($scope.classId).success(function(data){
      if (data !== undefined) {
        $scope.className = data["info"]["name"];
      }
    }).error(function(data){
      //handle error
    });
  }

  $scope.addQuestion = function() {
    if ($scope.questionform.$valid) {
      $scope.quiz.questions.push({text: $scope.question.text, correct: -1, answers: []});
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
    $scope.quiz.questions.splice($scope.quiz.questions.indexOf(question), 1);
  }

  $scope.removeAnswer = function(question, answer) {
    question.answers.splice(question.answers.indexOf(answer), 1);
  }

  $scope.postQuiz = function() {
    if ($scope.quizform.$valid) {
      quizService.createQuiz($scope.classId, $scope.quiz);
      $location.path('/main');
    }
  }




  /*
   * SESSION STUFF
   * */

  if ($location.$$path.match(/(\/quiz\/[0-9]+)/)) {
    // open dat socket
    quizSessionService.startSession($routeParams.id);
  }

  // start the quiz, send state of 0 to the server
  $scope.startQuiz = function() {
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
    quizSessionService.changeState($scope.current);
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
    quizSessionService.changeState($scope.current);
  }

  $scope.endQuiz = function() {
    /*
     *$http({
     *  method: 'PUT',
     *  url: '/quiz/' + $scope.id + '/state',
     *  data: angular.toJson({state: -1}),
     *  headers: {'Content-Type': 'application/json'}
     *});
     */
    quizSessionService.endSession();
    // redirect to main
    $location.path('/main');
  }
});
