var app = angular.module('labApp', ['ngResource', 'ngRoute']);
var host = 'http://172.22.206.54:8081'

app.factory("Get", function($resource) {
    return $resource('/status', {})
});

app.factory("PanTest", function($resource) {
    return $resource('/pan', {})
});

app.factory("PanWebTest", function($resource) {
    return $resource('/panweb', {})
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

app.controller('PanController', function($scope, $log, PanTest, PanWebTest) {
      $scope.webresult = 'No Test';
      $scope.bypassresult = 'No Test';
      $scope.dropresult = 'No Test';

      $scope.itemClicked = function () {
        $scope.weblabel = "label-info";
        $scope.bypasslabel = "label-info";
        $scope.droplabel = "label-info";
        $scope.webresult = 'Running';
        $scope.bypassresult = 'Running';
        $scope.dropresult = 'Running';
        PanWebTest.query(function(data) {
            if (data[0].Working) {
                $scope.weblabel = "label-success";
                $scope.webresult = 'Passed';
            } else {
                $scope.weblabel = "label-danger";
                $scope.webresult = 'Failed';
            }
        });

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

app.config(function($routeProvider, $httpProvider) {
    $httpProvider.defaults.useXDomain = true;
    delete $httpProvider.defaults.headers.common['X-Requested-With'];
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