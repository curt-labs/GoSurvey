define([],function(){
	'use strict';
	var ctlr = ['$scope','SurveyService', function($scope, SurveyService){
		SurveyService.query({page:'1', count:'4'},function(surveys){
			$scope.surveys = surveys.surveys;
		});
	}];

	return ctlr;
});