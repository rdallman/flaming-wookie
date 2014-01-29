

var quizApp = angular.module('quizApp', ['ngRoute'])
	.config(function ($routeProvider, $http, $scope) {
		$routeProvider.when('/quiz', {
			template: 'quiz.html', 
			controller: 'QuizController' 
		}).when('/quiz/:id', {
			template: 'quiz.html',
			controller: 'QuizController',
			resolve: {
				$http({
					method: 'GET',
					url: '/quiz/'
				})
			}
		})
		;
	});