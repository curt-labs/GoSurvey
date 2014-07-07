define([
	'jquery',
	'angular',
	'ngResource',
	'ngRoute',
	'angular-alert',
	'controllers/controllers',
	'services/services',
	'filters/filters',
	'directives/directives'],
	function($, angular, resource, route, alert, controllers, services, filters, directives){
		

		var initialize = function(){
			var mainModule = angular.module('app',['ngResource', 'ngRoute', 'ui.bootstrap.alert']);
			services.initialize(mainModule);
			controllers.initialize(mainModule);
			filters.initialize(mainModule);
			directives.initialize(mainModule);

			angular.bootstrap(window.document, ['app']);
		};

		return {
			initialize: initialize
		};

});