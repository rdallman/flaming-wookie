var dashboardApp = angular.module('dashboardApp', ['ngRoute', 'quizControllers']);

dashboardApp.config(['$routeProvider',
                    function($routeProvider) {
                      $routeProvider.
                        when('/main', {
                          templateUrl: '/templates/partials/main.html',
                          controller: 'MainController'
                        }).
                        when('/quiz-create', {
                          templateUrl: '/templates/partials/quiz-form.html',
                          controller: 'QuizController'
                        }).
                        when('/quiz/:id', {
                          templateUrl: '/templates/partials/quiz.html',
                          controller: 'QuizController'
                        }).
                        otherwise({
                          redirectTo: '/main'
                        });
                    }]);

dashboardApp.controller('MainController', function($scope, $http) {
  $scope.quizzes = [];
  
  $http({
    method: 'GET',
    url: '/dashboard/quiz'
  }).
  success(function(data) {
    angular.forEach(data, function(value, key) {
      this.push({"Title": key, "id": value});
    }, $scope.quizzes);
  }).
  error(function(data) {
    // handle
  });
});

