define([],function(){
	'use strict';

	var service = ['$resource',function($resource){
		return $resource('/api/prize/:id', {id: '@id'},{
			query:{method:'GET', isArray:false},
			current:{method:'GET', isArray:false}
		});
	}];

	return service;
});