define(['jquery','angular'],function($, angular){
	'use strict';

	var directives = {};

	var initialize = function(angModule){
		$.each(directives, function(name, filter) {
			 angModule.directive(name, filter);
		});
	};

	return {
		initialize: initialize
	};
});