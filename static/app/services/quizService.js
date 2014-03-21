
angular.module('dashboardApp').factory('quizService', function($http) {
  return {
    getQuizzes: function(cid) {
      return $http({
        method: 'GET',
        url: '/classes/' + cid + '/quiz'
      })
    },

    deleteQuiz: function(qid) {
      return $http({
        method: 'DELETE',
        url: '/quiz/' + qid
      })
    }
  }
});
