angular.module('dashboardApp').controller('AttendanceController', function (sessionService, attendanceService, classService, $scope, $http, $route, $routeParams, $location, flash) {

	$scope.current = -1;
	// get class
	classService.getClass($routeParams.cid).
    success(function(data) {
      if (data !== undefined) {
        $scope.class = data["info"];
        $scope.cid = $routeParams.cid;
        getAttendance($scope.cid);
        if ($location.$$path.match(/(\/classes\/[0-9]+\/attendance)$/)) {
			// open dat socket
			sessionService.startAttendanceSesh($scope.cid);
		}
      }
    });

    function getAttendance(cid) {
    	attendanceService.getClassAttendance(cid).success(function(data){
    		if (data !== undefined) {
    			$scope.attendance = data["info"];
    		}
    	});
    } 

    $scope.startAttendance = function() {
    	$scope.current = 0;
    	
    	sessionService.changeState($scope.current);
    }

    $scope.endAttendance = function() {
    	sessionService.endSession();
    	$location.path('/main');
    }

    $scope.filterStudents = function(input) {
    	$scope.class.students.forEach(function(student) {
    		if (student["sid"] = input) {
    			return student["fname"] + " " + student["lname"];
    		}
    	})
    }


});