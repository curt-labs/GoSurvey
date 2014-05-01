define([
	'jquery',
	'angular',
	'controllers/controllers',
	'services/services',
	'filters/filters',
	'directives/directives'],
	function($, angular, controllers, services, filters, directives){
		'use strict';

		var initialize = function(){
			var mainModule = angular.module('app',['ngResource']);
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