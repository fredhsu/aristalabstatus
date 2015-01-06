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

app.factory("NeutronNet", function($resource) {
    return $resource('/api/openstack/neutron', {})
});
app.factory("NeutronSubnet", function($resource) {
    return $resource('/api/openstack/subnet', {})
});
app.factory("NovaCompute", function($resource) {
    return $resource('/api/openstack/nova', {})
});
app.factory("StackEos", function($resource) {
    return $resource('/api/openstack/eos', {})
});
app.factory("StackReset", function($resource) {
    return $resource('/api/openstack/reset', {})
});

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

app.controller('OpenstackController', function($scope, $log, NeutronNet, NeutronSubnet, NovaCompute, StackEos, StackReset) {
      $scope.netresult = 'No Test';
      $scope.subnetresult = 'No Test';
      $scope.computeresult = 'No Test';
      $scope.eosresult = 'No Test';
      $scope.resetresult = 'No Test';

      $scope.itemClicked = function () {
        $scope.netlabel = "label-info";
        $scope.subnetlabel = "label-info";
        $scope.computelabel = "label-info";
        $scope.eoslabel = "label-info";
        $scope.resetlabel = "label-info";

        $scope.netresult = 'Running';
        $scope.subnetresult = 'Running';
        $scope.computeresult = 'Running';
        $scope.eosresult = 'Running';
        $scope.resetresult = 'Running';

        NeutronNet.query(function(data) {
            if (data[0].Working) {
                $scope.netlabel = "label-success";
                $scope.netresult = 'Passed';
            } else {
                $scope.netlabel = "label-danger";
                $scope.netresult = 'Failed';
            }
        });

        // NeutronSubnet.query(function(data) {
        //     $log.log(data);
        //     if (data[0].Working) {
        //         $scope.subnetlabel = "label-success";
        //         $scope.subnetresult = 'Passed';
        //     } else {
        //         $scope.subnetlabel = "label-danger";
        //         $scope.subnetresult = 'Failed';
        //     }
        // });
        // NovaCompute.query(function(data) {
        //     $log.log(data);
        //     if (data[0].Working) {
        //         $scope.computelabel = "label-success";
        //         $scope.computeresult = 'Passed';
        //     } else {
        //         $scope.computelabel = "label-danger";
        //         $scope.computeresult = 'Failed';
        //     }
        // });
        // StackEos.query(function(data) {
        //     $log.log(data);
        //     if (data[0].Working) {
        //         $scope.eoslabel = "label-success";
        //         $scope.eosresult = 'Passed';
        //     } else {
        //         $scope.eoslabel = "label-danger";
        //         $scope.eosresult = 'Failed';
        //     }
        // });
        // StackReset.query(function(data) {
        //     $log.log(data);
        //     if (data[0].Working) {
        //         $scope.resetlabel = "label-success";
        //         $scope.resetresult = 'Passed';
        //     } else {
        //         $scope.resetlabel = "label-danger";
        //         $scope.resetresult = 'Failed';
        //     }
        // });
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
        .when('/openstack', {
                templateUrl : 'openstack.html',
                controller  : 'OpenstackController'
        })

});
