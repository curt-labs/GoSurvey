define([], function () {
	'use strict';

	var service = ['$resource', function ($resource) {
		return $resource('api/survey/:id', {id: '@id'}, {});
	}];

	return service;
});