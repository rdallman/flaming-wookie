// quiz
var quizApp = angular.module('quizControllers', ['ngRoute']);

quizApp.controller('QuizController', function ($scope, $http, $route, $routeParams, $location) {
  
  $scope.current = -1;
  $scope.id = -1;
	$scope.quiz = {
					title: "",
					questions: []
					};
  if ($routeParams.id !== undefined) {	
    $http({
        method: 'GET',
        url: '/dashboard/quiz/' + $routeParams.id
    }).success(function(data) {
        if (data !== undefined) {
          $scope.quiz = data;
          $scope.id = $routeParams.id
        }
    }).error(function(data) {
        // handle error
    });
  }

	$scope.addQuestion = function(textIn) {
		$scope.quiz.questions.push({text: textIn, correct: -1, answers: []});
		$scope.newQuestion = "";
	}

	$scope.changeTitle = function(text) {
		$scope.quiz.title = text;
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
		$http({
				method: 'POST', 
				url: '/quiz/add', 
				data: JSON.stringify($scope.quiz),
				headers: {'Content-Type': 'application/json'}
			});
	}

  $scope.startQuiz = function() {
    // post to server we're starting the quiz
    $http({
        method: 'PUT',
        url: '/quiz/' + $scope.id + '/state?state=0'
    });
    // set current to start of questions
    $scope.current = 0;
  }

  $scope.nextQuestion = function() {
    // post to server that we're moving on...
    $http({
      method: 'PUT',
      url: '/quiz/' + $scope.id + '/state?state=1'
    });
    $scope.current++;
  }
  
  $scope.endQuiz = function() {
    $http({
      method: 'PUT',
      url: '/quiz/' + $scope.id + '/state?state=-1'
    });
    // redirect to main
    $location.path('/main');
  }
});
