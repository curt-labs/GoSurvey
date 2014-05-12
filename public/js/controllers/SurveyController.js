define([],function(){
	'use strict';

	var ctlr = ['$scope', '$route', '$routeParams', 'SurveyService',function($scope, $route, $routeParams, SurveyService){
		SurveyService.get({id: $routeParams.id}, function(survey){
			for (var i = 0; i < survey.questions.length; i++) {
				var answers = [];
				for (var j = 0; j < survey.questions[i].answers.length; j++) {
					var ans = survey.questions[i].answers[j];
					if(ans.data_type == 'multiple'){
						if(survey.questions[i].selects === undefined){
							survey.questions[i].selects = [];
						}
						survey.questions[i].selects.push(ans);
					}else{
						answers.push(ans);
					}
				}
				survey.questions[i].answers = answers;
			}

			$scope.survey = survey;
		});
	}];

	return ctlr;
});
