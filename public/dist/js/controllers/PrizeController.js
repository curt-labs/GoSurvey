define([],function(){
	

	var ctlr = ['$scope','PrizeService', function($scope, PrizeService){
		PrizeService.get({id: 'current'},function(prize){
			$scope.prize = prize;
		});
	}];

	return ctlr;
});