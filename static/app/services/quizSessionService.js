// service for the quiz sessions, has a websocket connection to the server and provides functions for the
// controllers to interface with the websocket
angular.module('dashboardApp').factory('quizSessionService', function($http, $route, $routeParams) {
  var Service = {};
  var ws;

  Service.startSession = function(qid) {
    ws = new WebSocket("ws://localhost:8080/giveme/" + qid);
    // set up handlers
    ws.onopen = function() {
      console.log("Socket connection open for: " + qid);
    };
    ws.onclose = function() {
      console.log("Socket connection closed.");
    };
  };

  Service.changeState = function(stateIn) {
    // send change state
    ws.send(angular.toJson({state: stateIn}));
    console.log("Changed state: " + stateIn);
  }

  Service.endSession = function() {
    console.log("Sending end session signal...");
    ws.send(angular.toJson({state: -1}));
    ws.close();
  }


  return Service;

});
