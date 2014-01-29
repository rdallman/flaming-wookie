// quiz
var quizApp = angular.module('quizControllers', []);

quizApp.controller('QuizController', function ($scope, $http) {

	$scope.quiz = {
					title: "",
					questions: []
					};
	

	$scope.addQuestion = function(text) {
		$scope.quiz.questions.push({questionText: text, answers: []});
		$scope.newQuestion = "";
	}

	$scope.changeTitle = function(text) {
		$scope.quiz.title = text;
	}

	$scope.addAnswer = function(question, text) {
		question.answers.push(text);
		//$scope.newAnswers[question] = "";

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

});