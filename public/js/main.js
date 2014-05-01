requirejs.config({
	baseUrl: '/js/vendor'
});

require(['common'],function(common){
	console.log(common);
});