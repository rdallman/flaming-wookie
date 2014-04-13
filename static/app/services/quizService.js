
angular.module('dashboardApp').factory('quizService', function($http) {
  return {
    getQuizzes: function(cid) {
      return $http({
        method: 'GET',
        url: '/classes/' + cid + '/quiz'
      })
    },

    getGrades: function(qid) {
      return $http({
        method: 'GET',
        url: '/quiz/' + qid + '/grades'
      })
    },

    deleteQuiz: function(qid) {
      return $http({
        method: 'DELETE',
        url: '/quiz/' + qid
      })
    },

    getQuiz: function(qid) {
      return $http({
        method: 'GET',
        url: '/quiz/' + qid
      })
    },

    createQuiz: function(cid, quiz) {
      return $http({
        method: 'POST',
        url: '/classes/' + cid + '/quiz',
        data: angular.toJson(quiz),
        headers: {'Content-Type': 'application/json'}
      });
    }
  }
});
