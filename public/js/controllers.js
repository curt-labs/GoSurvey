define(['angular'],function(angular){
	'use strict';

	angular.module('app.controllers',['$scope','Surveys'])
		.controller('HomeController',['$scope','Surveys',function($scope, Surveys){
			$scope.twoTimesTwo = 2 * 2;
			console.log('hi');
	}]);
});
