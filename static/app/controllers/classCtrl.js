
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
      $location.path('/main');
    }
    else {
      classService.createClass($scope.class);
      $location.path('/main');
    }
  }

  $scope.addStudent = function(email, firstName, lastName) {
    if ($location.$$path.match(/(\/classes\/[0-9]+\/edit)/)) {
      classService.addStudent(email, firstName, lastName).success(function(data) {
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
