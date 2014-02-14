
// class controller for
var classApp = angular.module('classControllers', ['ngRoute']);

classApp.controller('ClassController', function ($scope, $http, $route, $routeParams, $location) {

  $scope.myClass = {
    name: "",
    students: []
  };
  $scope.classes = [];
  $scope.id = -1; 

  $http({
    method: 'GET',
    url: '/classes'
  }).
  success(function(data) {
  console.log(data);  
    angular.forEach(data["info"], function(value, key) {
      this.push(value)
    }, $scope.classes);
  }).
  error(function(data) {
    // handle
  });

  if ($routeParams.id !== undefined) {	
    $http({
      method: 'GET',
      url: '/classes/' + $routeParams.id
    }).success(function(data) {
      if (data !== undefined) {
      $scope.myClass = data["info"]
      $scope.id = $routeParams.id;
      }
    }).error(function(data) {
      // handle error
    });
  }

  $scope.postClass = function() {
    $http({
      method: 'POST', 
    url: '/classes', 
    data: angular.toJson($scope.myClass),
    headers: {'Content-Type': 'application/json'}
    });
    $location.path('/classes');
  }
  $scope.addStudent = function(name, email) {
    $scope.myClass.students.push({name: name, email:email});
    $scope.newStudent = "";
  }

  $scope.changeName = function(text) {
    $scope.myClass.name = text;
  }
  $scope.removeStudent = function(student) {
    $scope.myClass.students.splice($scope.myClass.students.indexOf(student), 1);
  }
});
