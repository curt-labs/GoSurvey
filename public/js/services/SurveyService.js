define([], function () {
	'use strict';

	var service = ['$resource', function ($resource) {
		return $resource('api/survey/:id', {id: '@id'},{
			query:{method:'GET', isArray: false},
			get:{method: 'GET', isArray:false},
			post:{method: 'POST'}
		});
	}];

	return service;
});