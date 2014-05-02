define([], function () {
	

	var service = ['$resource', function ($resource) {
		return $resource('api/survey/:id', {}, {
			query: {method: 'GET', params: {id: 0}, isArray: true}
		});
	}];

	return service;
});