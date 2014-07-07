define(['jquery',
	'services/WarrantyService',
	'services/SurveyService',
	'services/PrizeService'],function($, ws,ss, ps){
	

	var services = {
		SurveyService: ss,
		WarrantyService: ws,
		PrizeService: ps
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