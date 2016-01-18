'use strict';
angular.module('irepApp', ['ngRoute','ngSanitize', 'swaggerUi','swagger-client']).run(function(swaggerModules, swagger1ToSwagger2Converter){
        swaggerModules.add(swaggerModules.BEFORE_PARSE, swagger1ToSwagger2Converter);
	}).run(function($rootScope, swaggerClient){
	IrepSchema.basePath = "http://xyz"
	var api = swaggerClient(IrepSchema);
      api.irep01.version().then(function(ver){
        $rootScope.ver = ver;
      });

    })
  .config(function ($routeProvider) {
    $routeProvider
      .when('/', {
        templateUrl: 'views/main.html',
        controller: 'MainCtrl'
      })
      .otherwise({
        redirectTo: '/'
      });
  });
