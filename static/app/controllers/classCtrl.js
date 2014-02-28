
// class controller for
var classApp = angular.module('classControllers', ['ngRoute']);

classApp.controller('ClassController', function ($scope, $http, $route, $routeParams, $location) {

  $scope.myClass = {
    name: "",
    students: []
  };
  $scope.classes = [];
  $scope.id = -1; 

  $scope.class = {
    name: "",
    students: []
  }

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

  $scope.createClass = function() {
      //alert(angular.toJson($scope.class));
    $http({
      method: 'POST',
      url: '/classes',
      data: angular.toJson($scope.class),
      headers: {'Content-Type': 'application/json'}
    });
    $location.path('/main');
    
  }

  $scope.addStudent = function(email, firstName, lastName) {
    $scope.class.students.push({email: email, fname: firstName, lname: lastName});
    $scope.student.email = "";
    $scope.student.fname = "";
    $scope.student.lname = "";

  }

  

  $scope.changeName = function(text) {
    $scope.myClass.name = text;
  }
  $scope.removeStudent = function(student) {
    $scope.myClass.students.splice($scope.myClass.students.indexOf(student), 1);
  }
});
