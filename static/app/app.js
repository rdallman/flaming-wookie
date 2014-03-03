// sets up app for angular, includes modules
// TODO better/cleaner way of including modules?
// TODO dependencies without including every file in html?
var dashboardApp = angular.module('dashboardApp', ['ngRoute', 'quizControllers', 'classControllers', 'customFilters']);

// routes handler for angular
// TODO better way of setting up routes? also resource handling
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
                        when('/class-form', {
                          templateUrl: '/templates/partials/class-form.html',
                          controller: 'ClassController'
                        }).
                        when('/quiz/:id', {
                          templateUrl: '/templates/partials/quiz.html',
                          controller: 'QuizController'
                        }).
                        when('/quiz/:id/grades', {
                          templateUrl: '/templates/partials/grades.html',
                          controller: 'QuizController'
                        }).
                        when('/classes', {
                          templateUrl: '/templates/partials/classes.html',
                          controller: 'ClassController'
                        }).
                        when('/classes/:id', {
                          templateUrl: '/templates/partials/class.html',
                          controller: 'ClassController'
                        }).
                        when('/classes/:id/edit', {
                          templateUrl: '/templates/partials/class-form.html',
                          controller: 'ClassController'
                        }).
                        when('/quizzes', {
                          templateUrl: '/templates/partials/quizzes.html',
                          controller: 'QuizController'
                        }).
                        otherwise({
                          redirectTo: '/main'
                        });
                    }]);

// clear cache
// REMOVE WHEN IN PRODUCTION
dashboardApp.run(function ($rootScope, $templateCache) {
  $rootScope.$on('$viewContentLoaded', function() {
    $templateCache.removeAll();
  });
});

// custom filters FTW!
filters = angular.module('customFilters', []);
filters.filter('abc', function() {
  return function(input) {
    return String.fromCharCode(input + 65);
  };
});

// main controller for the main dashboard page
// TODO pull out into another file
dashboardApp.controller('MainController', function($scope, $http) {
  $scope.classList = [];
  $scope.quizList = [];
  
  // get classes
  $http({
    method: 'GET',
    url: '/classes'
  }).
  success(function(data) {
    // angular.forEach(data, function(value, key) {
    //   this.push({"Title": key, "id": value});
    // }, $scope.classes);

    $scope.classList = data.info
  }).
  error(function(data) {
    // handle
  });

  // get quizzes
  $http({
    method: 'GET',
    url: '/quiz'
  }).
  success(function(data) {
    // angular.forEach(data, function(value, key) {
    //   this.push({"Title": key, "id": value});
    // }, $scope.classes);
    $scope.quizList = data.info
  }).
  error(function(data) {
    // handle
  });
});

