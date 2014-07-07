if(typeof define !== 'function'){
	var define = require('amdefine')(module);
}

require.config({
	paths: {
		'angular': './vendor/angular/angular.min',
		'angular-alert':'./vendor/angular-ui-bootstrap/src/alert/alert',
		'ngResource': './vendor/angular-resource/angular-resource.min',
		'ngRoute': './vendor/angular-route/angular-route.min',
		'jquery': './vendor/jquery/dist/jquery.min',
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
		'ngResource': ['angular'],
		'angular-alert':['angular']
	},
	waitSeconds: 15,
	urlArgs: 'bust=v0.1.2'
});

require([
	'require',
	'jquery',
	'angular',
	'bootstrap',
	'respondJS',
	'nprogress',
	'holder',
	'angular-alert'], function(require, $, angular){
		require(['app'],function(app){
			app.initialize();
		});
});
