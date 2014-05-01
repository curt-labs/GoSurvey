define(['angular'],function(angular){
	'use strict';

	angular.module('app.services', ['ngResource'])
	  .factory('AngularIssues', function($resource){
	    return $resource('https://api.github.com/repos/angular/angular.js/issues/:number',
	      {number: '@number'},
	      {getIssue: {method: 'GET', params: {number: 0}}}
	    )
	}).value('version', '0.1');
});