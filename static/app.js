var quizApp = angular.module('quizApp', []);

quizApp.controller('QuizController', function ($scope, $http) {

	$scope.questions = ["What's my name?"];

	$scope.addQuestion = function(text) {
		$scope.questions.push(text);
	}

});