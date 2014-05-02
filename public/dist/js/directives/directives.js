define(['jquery','angular'],function($, angular){
	

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