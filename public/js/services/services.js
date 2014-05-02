define(['jquery', 'services/WarrantyService','services/SurveyService'],function($, ws,ss){
	'use strict';

	var services = {
		SurveyService: ss,
		WarrantyService: ws
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