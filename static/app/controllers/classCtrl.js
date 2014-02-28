
// class controller for
var classApp = angular.module('classControllers', ['ngRoute']);

classApp.controller('ClassController', function ($scope, $http, $route, $routeParams, $location) {

  $scope.classes = [];
  $scope.id = -1; 

  $scope.class = {
    name: "",
    students: []
  }

  // class list
  if ($routeParams.id === undefined) {
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
  }
  // specific class
  else {
    $http({
      method: 'GET',
      url: '/classes/' + $routeParams.id
    }).success(function(data) {
      if (data !== undefined) {
        $scope.class = data["info"]
        $scope.id = $routeParams.id;
        // we're editing the class
        // TODO better way to check path that has id param...
        if ($location.$$path.match(/(\/classes\/[0-9]+\/edit)/)) {
        }
      }
    }).error(function(data) {
      // handle error
    });
  }

  $scope.createClass = function() {
    if ($location.$$path.match(/(\/classes\/[0-9]+\/edit)/)) {

    }
    else {
      
      $http({
        method: 'POST',
        url: '/classes',
        data: angular.toJson($scope.class),
        headers: {'Content-Type': 'application/json'}
      });
      $location.path('/main');
    }
  }

  $scope.addStudent = function(email, firstName, lastName) {
    if ($location.$$path.match(/(\/classes\/[0-9]+\/edit)/)) {
      $http({
        method: 'POST',
        url: '/classes/' + $routeParams.id + '/student',
        data: angular.toJson({cid: $routeParams.id, email: email, fname: firstName, lname: lastName}),
        headers: {'Content-Type': 'application/json'}
      }).success(function(data) {
        $scope.class.students.push({email: email, name: {first: firstName, last: lastName}}); 
      }).error(function(data) {
        alert("Error: Could not add student");
      });
    }
    else {
      $scope.class.students.push({email: email, fname: firstName, lname: lastName});
      $scope.student.email = "";
      $scope.student.fname = "";
      $scope.student.lname = "";    }

  }

  

  $scope.changeName = function(text) {
    $scope.class.name = text;
  }
  $scope.removeStudent = function(student) {
    $scope.class.students.splice($scope.class.students.indexOf(student), 1);
  }
});
