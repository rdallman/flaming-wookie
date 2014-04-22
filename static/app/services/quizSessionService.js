// service for the quiz sessions, has a websocket connection to the server and provides functions for the
// controllers to interface with the websocket
angular.module('dashboardApp').factory('sessionService', function($http, $route, $routeParams) {
  var Service = {};
  var ws;

  Service.startSession = function(qid) {
    ws = new WebSocket("ws://24.178.89.28:8080/giveme/" + qid);
    // set up handlers
    ws.onopen = function() {
      console.log("Socket connection open for quiz/poll: " + qid);
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

  Service.startAttendanceSesh = function(cid) {
    ws = new WebSocket("ws://24.178.89.28:8080/takeAttendance/" + cid);

    ws.onopen = function() {
      console.log("Socket connection open for attendance: " + cid);
    };
    ws.onclose = function() {
      console.log("Socket connection closed.");
    };
  }


  return Service;

});
