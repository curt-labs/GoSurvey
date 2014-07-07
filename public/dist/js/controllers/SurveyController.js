define([],function(){
	

	var ctlr = ['$scope',
		'$route',
		'$routeParams',
		'$http',
		'$location',
		'SurveyService',
	function($scope, $route, $routeParams, $http, $location, SurveyService){
		$scope.alerts = [];
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
		},function(){
			$location.path("/surveys");
		});

		$scope.submitSurvey = function(){
			var resp = {
				user:$scope.survey.user,
				id:$scope.survey.id,
				questions:[]
			};

			for (var i = 0; i < $scope.survey.questions.length; i++) {
				var q = $scope.survey.questions[i];
				var question = {
					id:q.id,
					answer:q.answer
				};
				resp.questions.push(question);
			}
			$scope.alerts = [];
			SurveyService.post(resp).$promise.then(function(data){
				$scope.alerts.push({type: 'success', msg: 'Survey has successfully been submitted.'});
				$('.survey-form').find('input,select,textarea').val('');
			},function(err){
				$scope.alerts.push({type: 'warning',msg:err.data});
			});
		};

		$scope.closeAlert = function(index){
			$scope.alerts.splice(index,1);
		};
	}];

	return ctlr;
});
