define([
	'jquery',
	'routes/routes',
	'controllers/AppController',
	'controllers/HomeController',
	'controllers/DataController'],
	function($, routes, app, home, data){

		var controllers = {
			home: home,
			data: data
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
				$rootScope.$on('$routeChangeSuccess', function(next, last){
					console.log('Navigating from ', last);
					console.log('Navigating to ', next);
				})
			})
		};

		var initialize = function(angModule){
			angModule.controller('AppController', app)
			$.each(controllers, function(name, ctrl) {
				angModule.controller(name, ctrl);
			});
			setUpRoutes(angModule);
		};

		return {
			initialize: initialize
		};
});