define([
	'/js/vendor/requirejs-text/text.js!/js/views/home.html',
	'/js/vendor/requirejs-text/text.js!/js/views/data.html'
],function(homeTemplate,dataTemplate){
	return {
		home: {
			title: 'Home',
			route: '/home',
			controller: 'home',
			template: homeTemplate
		},
		data: {
			title: 'Data List',
			route: '/data',
			controller: 'data',
			template: dataTemplate
		}
	};
});