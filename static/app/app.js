// sets up the module for the entire app
// TODO better/cleaner way of including modules?
// TODO dependencies without including every file in html?
var dashboardApp = angular.module('dashboardApp', ['ngRoute', 'customFilters', 'flash']);

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
                        when('/poll/:qid', {
                          templateUrl: '/templates/partials/poll.html',
                          controller: 'PollController'
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
                        when('/classes/:cid/attendance', {
                          templateUrl: '/templates/partials/attendance.html',
                          controller: 'AttendanceController'
                        }).
                        when('/classes/:cid/view-attendance', {
                          templateUrl: '/templates/partials/view-attendance.html',
                          controller: 'AttendanceController'
                        }).
                        when('/classes/:cid/quiz/:qid', {
                          templateUrl: '/templates/partials/quiz.html',
                          controller: 'QuizController'
                        }).
                        when('/classes/:cid/new-quiz', {
                          templateUrl: '/templates/partials/quiz-form.html',
                          controller: 'QuizController'
                        }).
                        when('/classes/:cid/new-poll', {
                          templateUrl: '/templates/partials/poll-form.html',
                          controller: 'PollController'
                        }).
                        when('/classes/:cid/edit', {
                          templateUrl: '/templates/partials/class-form.html',
                          controller: 'ClassController'
                        }).
                        when('/classes/:cid/poll/:qid/results', {
                          templateUrl: '/templates/partials/poll-results.html',
                          controller: 'PollController'
                        }).
                        otherwise({
                          redirectTo: '/main'
                        });
                    }]);

// clear cache
// REMOVE WHEN IN PRODUCTION
dashboardApp.run(function ($rootScope, $templateCache, flash) {
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

filters.filter('filterStudentName', function() {
  return function(input, students) {
    for (i = 0; i < students.length; i++) {
      if (students[i].sid = input) {
        return students[i].fname + " " + students[i].lname;
      }
    }
  };
});

filters.filter('filterPollResults', function() {
  return function(input, index, pIndex) {
    count = 0;
    angular.forEach(input[pIndex], function(value, key) {
      if (value == index) {
        count++;
      }
    });
    return count;
  };
});

