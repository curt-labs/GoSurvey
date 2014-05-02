define([
	'jquery',
	'nprogress',
	'routes/routes',
	'controllers/AppController',
	'controllers/HomeController',
	'controllers/WarrantyController',
	'controllers/SurveysController',
	'controllers/SurveyController'],
	function($, NProgress, routes, app, home, warranty, surveys, survey){

		var controllers = {
			home: home,
			warranty: warranty,
			surveys: surveys,
			survey: survey
		};

		var setUpRoutes = function(angModule){
			angModule.config(function($routeProvider) {
				$.each(routes, function(key, val) {
					$routeProvider.when(val.route,{
						template: val.template,
						controller: val.controller,
						title: val.title
					});
				});
				$routeProvider.otherwise({ redirectTo: routes.home.route });
			});
			angModule.run(function($rootScope){
				$rootScope.$on('$routeChangeStart',function(){
					NProgress.start();
				});
				$rootScope.$on('$routeChangeSuccess', function(next, last){
					NProgress.done();
				});
			});
		};

		var initialize = function(angModule){
			angModule.controller('AppController', app);
			$.each(controllers, function(name, ctrl) {
				angModule.controller(name, ctrl);
			});
			setUpRoutes(angModule);
		};

		return {
			initialize: initialize
		};
});
