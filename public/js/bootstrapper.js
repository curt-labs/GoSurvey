if(typeof define !== 'function'){
	var define = require('amdefine')(module);
}

require.config({
	paths: {
		'angular': './vendor/angular/angular.min',
		'ngResource': './vendor/angular-resource/angular-resource.min',
		'ngRoute': './vendor/angular-route/angular-route.min',
		'jquery': './vendor/jquery/dist/jquery.min',
		'html5shiv':'./vendor/html5shiv/dist/html5shiv.min',
		'respondJS': './vendor/respondJS/dest/respond.min',
		'bootstrap': './vendor/bootstrap/dist/js/bootstrap.min',
		'nprogress':'./vendor/nprogress/nprogress',
		'holder':'./vendor/holderjs/holder',
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
	'bootstrap',
	'respondJS',
	'nprogress',
	'holder',
	'html5shiv'], function(require, $, angular){
		require(['app'],function(app){
			app.initialize();
		});
});
