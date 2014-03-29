
// class controller for
//var classApp = angular.module('classControllers', ['ngRoute']);

angular.module('dashboardApp').controller('ClassController', function (classService, quizService, $scope, $http, $route, $routeParams, $location) {

  $scope.classList = [];
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
    classService.getClasses().
    success(function(data) {
      console.log(data);
      angular.forEach(data["info"], function(value, key) {
        this.push(value)
      }, $scope.classList);
    }).
    error(function(data) {
      // handle
    });
  }
  // specific class
  else {
    classService.getClass($routeParams.id).
    success(function(data) {
      if (data !== undefined) {
        $scope.class = data["info"]
        $scope.id = $routeParams.id;
        // we're editing the class
        // TODO better way to check path that has id param...
        if ($location.$$path.match(/(\/classes\/[0-9]+\/edit)/)) {
          $scope.editing = true;
        }
      }
    }).
    error(function(data) {
      // handle error
    });

    // get quizzes
    quizService.getQuizzes($routeParams.id).
    success(function(data) {
      $scope.quizList = data.info
    }).
    error(function(data) {
      // handle
    });
  }

  $scope.createClass = function() {
    if ($location.$$path.match(/(\/classes\/[0-9]+\/edit)/)) {
      // TODO save class updates
      $location.path('/classes/' + $routeParams.id);
    }
    else {
      if ($scope.classform.$valid) {
        classService.createClass($scope.class);
        $location.path('/main');
      }
    }
  }

  $scope.addStudent = function() {
    if ($location.$$path.match(/(\/classes\/[0-9]+\/edit)/)) {
      if ($scope.studentform.$valid) {
        classService.addStudent($scope.student.email, $scope.student.fname, $scope.student.lname).success(function(data) {
          $scope.class.students.push({email: $scope.student.email, fname: $scope.student.fname, lname: $scope.student.lname});
          $scope.student.email = "";
          $scope.student.fname = "";
          $scope.student.lname = "";
        }).error(function(data) {
          alert("Error: Could not add student");
        });
      }
    }
    else {
      if ($scope.studentform.$valid) {
        $scope.class.students.push({email: $scope.student.email, fname: $scope.student.fname, lname: $scope.student.lname});
        $scope.student.email = "";
        $scope.student.fname = "";
        $scope.student.lname = "";
      }
    }

  }

  // TODO add in check for deleting from class if editting existing class
  $scope.removeStudent = function(student) {
    $scope.class.students.splice($scope.class.students.indexOf(student), 1);
  }

  $scope.deleteQuiz = function(qid, index) {
    quizService.deleteQuiz(qid).success(function(data) {
      $scope.quizList.splice(index, 1);
    }).error(function(data) {
        alert("Error: Could not delete quiz");
      });
  }
});
