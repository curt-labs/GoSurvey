define(['routes/routes'],function(routes){
	
	var appCtrl = ['$scope',function($scope){
		$scope.navigation = routes;
	}];

	return appCtrl;
});