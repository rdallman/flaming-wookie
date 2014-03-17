
// class controller for
var classApp = angular.module('classControllers', ['ngRoute']);

classApp.controller('ClassController', function ($scope, $http, $route, $routeParams, $location) {

  $scope.classes = [];
  $scope.id = -1; 
  $scope.quizList = [];

  $scope.class = {
    name: "",
    students: []
  }
  // used for determining buttons on the view
  $scope.editing = false;

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
          $scope.editing = true;
        }
      }
    }).error(function(data) {
      // handle error
    });

    // get quizzes
    $http({
      method: 'GET',
      url: '/classes/' + $routeParams.id + '/quiz'
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
  }

  $scope.createClass = function() {
    if ($location.$$path.match(/(\/classes\/[0-9]+\/edit)/)) {
      $location.path('/main');
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
        data: angular.toJson({cid: parseInt($routeParams.id), email: email, fname: firstName, lname: lastName}),
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

  $scope.deleteQuiz = function(qid, index) {
    $http({
      method: 'DELETE',
      url: '/quiz/' + qid
    }).success(function(data) {
      $scope.quizList.splice(index, 1);
    }).error(function(data) {
        alert("Error: Could not delete quiz");
      });
  }
});
