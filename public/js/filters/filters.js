define(['jquery'],function($){
	'use strict';

	var filters = {};

	var initialize = function(angModule){
		$.each(filters, function(name, filter) {
			angModule.filter(name, filter);
		});
	};

	return {
		initialize: initialize
	};
});