define([
	'jquery',
	'filters/truncate',
	'filters/range'],function($, truncate, ranger){
	'use strict';

	var filters = {
		truncate: truncate,
		ranger: ranger
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