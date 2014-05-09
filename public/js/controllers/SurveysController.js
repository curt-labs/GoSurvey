define(['holder'],function(holder){
	'use strict';

	var ctlr = ['$scope', '$route', 'SurveyService',function($scope, $route, SurveyService){
		holder.run();
		SurveyService.query(function(surveys){
			$scope.surveys = surveys.surveys;
		});
	}];

	return ctlr;
});