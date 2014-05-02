define(['jquery','services/SurveyService'],function($,ss){
	'use strict';

	var services = {
		SurveyService: ss
	};

	var initialize = function(angModule){
		$.each(services, function(name, service) {
			angModule.factory(name, service);
		});
	};

	return {
		initialize: initialize
	};

});