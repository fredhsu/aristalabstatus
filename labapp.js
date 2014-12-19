


var app = angular.module('labApp', ['ngResource', 'ngRoute']);

app.factory("Get", function($resource) {
    return $resource('http://localhost:8081/switches', {})
});

app.factory("PanTest", function($resource) {
    return $resource('http://localhost:8081/pan', {})
});


// app.controller('SwitchesController', ['$scope', function($scope) {
app.controller('SwitchesController', function($scope, Get) {
    Get.query(function(data) {
        $scope.test = data;
    });
    $scope.switches = [
        {Hostname:'bleaf1', IpIntf:'test', IntfConnected:'testintf', Uptime:'00', Version:'11'}
    ];

});

app.controller('PanController', function($scope, $log, PanTest) {
      $scope.bypassresult = 'No Test';
      $scope.dropresult = 'No Test';

      $scope.itemClicked = function () {
        $scope.bypasslabel = "label-info";
        $scope.droplabel = "label-info";
        $scope.bypassresult = 'Running';
        $scope.dropresult = 'Running';
        PanTest.query(function(data) {
            $log.log(data);
            if (data[0].Working) {
                $scope.bypasslabel = "label-success";
                $scope.bypassresult = 'Passed';
            } else {
                $scope.bypasslabel = "label-danger";
                $scope.bypassresult = 'Failed';
            }
            if (data[1].Working) {
                $scope.droplabel = "label-success";
                $scope.dropresult = 'Passed';
            } else {
                $scope.droplabel = "label-danger";
                $scope.dropresult = 'Failed';
            }
        });
      };
});

app.config(function($routeProvider) {
    $routeProvider
        .when('/', {
            templateUrl : 'home.html'
        })
        .when('/overview', {
                templateUrl : 'overview.html',
                controller  : 'SwitchesController'
        })
        .when('/topology', {
                templateUrl : 'topology.html'
        })
        .when('/pan', {
                templateUrl : 'pan.html',
                controller  : 'PanController'
        })

});