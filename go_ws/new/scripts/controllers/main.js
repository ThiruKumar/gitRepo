'use strict';

angular.module('irepApp')
  .controller('MainCtrl', ['$scope', '$sce', function ($scope, $sce) {

$scope.timeRange = {
  fromTime: new Date('2015','07','20','00','00'),
  toTime: new Date('2015','07','21','00','00'),
}

var getGrafanaUrl = function() {
  var u = "http://192.168.70.9:3000/dashboard-solo/db/new-dashboard?"+
//  var u = "http://192.168.70.7:3000/dashboard/db/new-dashboard?"+
          "panelId=1&"+
//          "fullscreen&"+
          "theme=light&"+
          "from="+$scope.timeRange.fromTime.getTime().toString()+
          "&to="+$scope.timeRange.toTime.getTime().toString()
//  return u
  return $sce.trustAsResourceUrl(u);
}

$scope.getGrafanaUrl = getGrafanaUrl

$scope.url = '/apidocs.json';
// error management
$scope.myErrorHandler = function(data, status){
    alert('failed to load swagger: '+status);
};
// transform try it request
$scope.myTransform = function(request){
    request.headers['api_key'] = 's5hredf5hy41er8yhee58';
};
  }]);
