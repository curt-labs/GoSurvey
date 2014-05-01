/**
 * loads sub modules and wraps them up into the main module
 * this should be used for top-level module definitions only
 */
define([
  'angular',
  './config',
  './services',
  './controllers'
], function (angular) {
  'use strict';

  angular.module('app', [
  	'ngRoute',
    'app.constants',
    'app.controllers'
  ]).config(['$routeProvider', function ($routeProvider) {
  	console.log($routeProvider);
    $routeProvider.when('/',{
    	templateUrl: 'js/views/home.html',
    	controller: 'HomeController'
    });
  }]);

});