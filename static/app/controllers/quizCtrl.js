// quiz
//var quizApp = angular.module('quizControllers', ['ngRoute']);

angular.module('dashboardApp').controller('QuizController', function (quizSessionService, quizService, classService, $scope, $http, $route, $routeParams, $location, flash) {

  // used for traversing quiz in html
  $scope.current = -1;
  $scope.qid = -1;
  $scope.cid;
  $scope.className;
  $scope.grades;

  $scope.classes = [];

  // main model for quiz
  $scope.quiz = {
    title: "",
    questions: [],
    grades: {}
  };

  $scope.quizlist = []

  if ($routeParams.qid !== undefined) {	
    quizService.getQuiz($routeParams.qid).success(function(data) {
      if (data !== undefined) {
        $scope.quiz = data["info"]
        $scope.qid = $routeParams.qid;
        //alert($scope.quiz.grades.length > 0);
      }
    }).error(function(data) {
      // handle error
    });
  }

  if ($routeParams.cid !== undefined) {
    classService.getClass($routeParams.cid).success(function(data){
      if (data !== undefined) {
        $scope.className = data["info"]["name"];
        $scope.cid = $routeParams.cid;
      }
    });

  }

  /*if ($location.$$path === "/quizzes") {
    $http({
      method: 'GET',
      url: '/quiz'
    }).success(function(data) {
      $scope.quizlist = data["info"]

    }).error(function(data) {

    });
  }*/

  /*if ($location.$$path.match(/(\/classes\/[0-9]+\/new-quiz)/)) {
    // creating a quiz
    $scope.cid = $routeParams.cid;
    classService.getClass($scope.cid).success(function(data){
      if (data !== undefined) {
        $scope.className = data["info"]["name"];
      }
    }).error(function(data){
      //handle error
    });
  }*/

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
      quizService.createQuiz($scope.cid, $scope.quiz);
      $location.path('/quizzes');
      document.getElementById("flash").setAttribute("class", "alert alert-success");
      flash("You created a quiz!");
    }
  }

 if ($location.$$path.match(/(\/classes\/[0-9]+\/quiz\/[0-9]+\/grades)$/)) {
      quizService.getGrades($routeParams.qid).success(function(data) {
        $scope.grades = data["info"]
      });
      
  }




  /*
   * SESSION STUFF
   * */

  if ($location.$$path.match(/(\/quiz\/[0-9]+)$/)) {
    // open dat socket
    quizSessionService.startSession($routeParams.qid);
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
