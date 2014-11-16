


var app = angular.module('labApp', ['ngResource']);
// var app = angular.module('labApp', []);

app.factory("Get", function($resource) {
    return $resource('http://localhost:8081/switches', {})
});


// app.controller('SwitchesController', ['$scope', function($scope) {
app.controller('SwitchesController', function($scope, Get) {
    Get.query(function(data) {
    //     $scope.switches = data;
        $scope.test = data;
    });
    $scope.switches = [
        {Hostname:'bleaf1', IpIntf:'test', IntfConnected:'testintf', Uptime:'00', Version:'11'}
    ];

});