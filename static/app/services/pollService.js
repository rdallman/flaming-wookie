
angular.module('dashboardApp').factory('pollService', function($http) {
  return {
    getPolls: function(cid) {
      return $http({
        method: 'GET',
        url: '/classes/' + cid + '/polls'
      })
    },

    getResults: function(qid) {
      return $http({
        method: 'GET',
        url: '/poll/' + qid + '/results'
      })
    },

    deletePoll: function(qid) {
      return $http({
        method: 'DELETE',
        url: '/poll/' + qid
      })
    },

    getPoll: function(qid) {
      return $http({
        method: 'GET',
        url: '/poll/' + qid
      })
    },

    createPoll: function(cid, poll) {
      return $http({
        method: 'POST',
        url: '/classes/' + cid + '/poll',
        data: angular.toJson(poll),
        headers: {'Content-Type': 'application/json'}
      });
    }
  }
});
