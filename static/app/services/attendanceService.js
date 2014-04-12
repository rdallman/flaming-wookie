angular.module('dashboardApp').factory('attendanceService', function($http, $route, $routeParams) {
	var Service = {}

	Service.getClassAttendance = function(cid) {
		return $http({
			method: 'GET',
			url: '/classes/' + cid + '/attendance'
		});
	};

	return Service;

});