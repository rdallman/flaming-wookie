// sets up the module for the entire app
// TODO better/cleaner way of including modules?
// TODO dependencies without including every file in html?
var dashboardApp = angular.module('dashboardApp', ['ngRoute', 'customFilters']);

// routes handler for angular
// TODO better way of setting up routes? also resource handling
dashboardApp.config(['$routeProvider',
                    function($routeProvider) {
                      $routeProvider.
                        when('/main', {
                          templateUrl: '/templates/partials/main.html',
                          controller: 'ClassController'
                        }).
                        when('/quiz-create', {
                          templateUrl: '/templates/partials/quiz-form.html',
                          controller: 'QuizController'
                        }).
                        when('/class-form', {
                          templateUrl: '/templates/partials/class-form.html',
                          controller: 'ClassController'
                        }).
                        when('/quiz/:qid', {
                          templateUrl: '/templates/partials/quiz.html',
                          controller: 'QuizController'
                        }).
                        when('/classes/:cid/quiz/:qid/grades', {
                          templateUrl: '/templates/partials/grades.html',
                          controller: 'QuizController'
                        }).
                        when('/classes', {
                          templateUrl: '/templates/partials/classes.html',
                          controller: 'ClassController'
                        }).
                        when('/classes/:cid', {
                          templateUrl: '/templates/partials/class.html',
                          controller: 'ClassController'
                        }).
                        when('/classes/:cid/quiz/:qid', {
                          templateUrl: '/templates/partials/quiz.html',
                          controller: 'QuizController'
                        }).
                        when('/classes/:cid/new-quiz', {
                          templateUrl: '/templates/partials/quiz-form.html',
                          controller: 'QuizController'
                        }).
                        when('/classes/:cid/edit', {
                          templateUrl: '/templates/partials/class-form.html',
                          controller: 'ClassController'
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
filters.filter('gradeNA', function() {
  return function(input) {
    if (input < 0) {
      return "NA";
    } else {
      return input;
    }
  };
});

