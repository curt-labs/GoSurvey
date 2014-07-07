define(['holder'],function(holder){
	

	var ctlr = ['$scope', '$route', 'SurveyService',function($scope, $route, SurveyService){
		holder.run();
		$scope.isDisabled = function(dir){
			if(dir === 'newer' && $scope.current_page === 1){
				return 'disabled';
			}else if(dir === 'older' && $scope.current_page == $scope.pages.length){
				return 'disabled';
			}
		};
		$scope.nextPage = function(){
			$scope.current_page++;
			SurveyService.query({page:$scope.current_page},function(surveys){
				$scope.surveys = surveys.surveys;
			});
		};

		SurveyService.query(function(surveys){
			$scope.current_page = 0;
			var len = Math.ceil(surveys.total_surveys / surveys.total_results);
			$scope.pages = new Array(len);
			$scope.surveys = surveys.surveys;
			$scope.current_page++;
		});
	}];

	return ctlr;
});