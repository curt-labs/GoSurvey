/**
 * bootstraps angular onto the window.document node
 * NOTE: the ng-app attribute should not be on the index.html when using ng.bootstrap
 */

 if (typeof define !== 'function') {
	// to be able to require file from node
	var define = require('amdefine')(module);
}

require.config({
	// Here paths are set relative to `/source/js` folder
	paths: {
		'angular': '/js/vendor/angular/angular.min',
		'async': '/js/vendor/requirejs-plugins/src/async.min',
		'domReady': '/js/vendor/requirejs-domready/domReady',
		'ngResource': '/js/vendor/angular-resource/angular-resource.min',
		'ngRoute': '/js/vendor/angular-route/angular-route.min',
		'jquery': '/js/vendor/jquery/dist/jquery.min',
		'bootstrap': '/js/vendor/bootstrap-sass-official/vendor/assets/javascripts/bootstrap',
		'app.controllers':'/js/controllers',
		'app.services':'/js/services',
		'app.constants':'/js/config'
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
		'app.controllers': ['angular'],
		'app.constants': ['angular'],
		'app.services': ['angular']
	},
	waitSeconds: 15,
	urlArgs: 'bust=v0.1.0',
	baseUrl: '/js'
});

require([
  'require',
  'angular',
  'app',
  'ngRoute'
], function (require, angular) {
  'use strict';

  /*place operations that need to initialize prior to app start here
   * using the `run` function on the top-level module
   */

  require(['domReady!'], function (document) {
    /* everything is loaded...go! */
    angular.bootstrap(document, ['app']);
  });
});
