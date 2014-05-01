if(typeof define !== 'function'){
	var define = require('amdefine')(module);
}

require.config({
	paths: {
		'angular': './vendor/angular/angular.min',
		'ngResource': './vendor/angular-resource/angular-resource.min',
		'ngRoute': './vendor/angular-route/angular-route.min',
		'jquery': './vendor/jquery/dist/jquery.min',
		'bootstrap': './vendor/bootstrap-sass-official/vendor/assets/javascripts/bootstrap',
		'templates':'./views'
	},

	shim: {
		'jquery':{
			'exports':'$'
		},
		'angular': {
			'exports': 'angular'
		},
		'bootstrap':['jquery'],
		'ngRoute':['angular'],
		'ngResource': ['angular']
	},
	waitSeconds: 15,
	urlArgs: 'bust=v0.1.0'
});

require([
	'require',
	'jquery',
	'angular',
	'bootstrap'], function(require, $, angular){
		require(['app'],function(app){
			app.initialize();
		});
});