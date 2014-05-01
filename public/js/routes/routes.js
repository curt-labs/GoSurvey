define([
	'/js/vendor/requirejs-text/text.js!templates/home.html',
	'/js/vendor/requirejs-text/text.js!templates/data.html'
],function(homeTemplate,dataTemplate){
	return {
		home: {
			title: 'Home',
			route: '/home',
			controller: 'home',
			template: homeTemplate
		},
		creation: {
			title: 'Data List',
			route: '/data',
			controller: 'data',
			template: dataTemplate
		}
	};
})