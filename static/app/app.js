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
                        otherwise({
                          redirectTo: '/main'
                        });
                    }]);

dashboardApp.controller('MainController', function($scope) {

});

