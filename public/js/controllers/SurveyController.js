define([],function(){
	'use strict';

	var ctlr = ['$scope', '$route', '$routeParams', 'SurveyService',function($scope, $route, $routeParams, SurveyService){
		console.log($routeParams);
		SurveyService.get({id: $routeParams.id}, function(survey){
			$scope.survey = survey;
		});
	}];

	return ctlr;
});
