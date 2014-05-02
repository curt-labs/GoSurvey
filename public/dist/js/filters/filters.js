define([
	'jquery',
	'filters/truncate'],function($, truncate){
	

	var filters = {
		truncate: truncate
	};

	var initialize = function(angModule){
		$.each(filters, function(name, filter) {
			angModule.filter(name, filter);
		});
	};

	return {
		initialize: initialize
	};
});