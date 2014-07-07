define([
	'/js/vendor/requirejs-text/text.js!/js/views/home.html',
	'/js/vendor/requirejs-text/text.js!/js/views/surveys/index.html',
	'/js/vendor/requirejs-text/text.js!/js/views/surveys/survey.html',
	'/js/vendor/requirejs-text/text.js!/js/views/warranty.html',
	'/js/vendor/requirejs-text/text.js!/js/views/policy.html',
	'/js/vendor/requirejs-text/text.js!/js/views/terms.html'
],function(homeTemplate, surveysTemplate, surveyTemplate, warrantyTemplate, policyTemplate, termsTemplate){
	return {
		home: {
			title: 'Home',
			route: '/home',
			controller: 'home',
			template: homeTemplate
		},
		warranty: {
			title: 'Warranty',
			route: '/warranty',
			controller: 'warranty',
			template: warrantyTemplate
		},
		surveys: {
			title: 'Surveys',
			route: '/surveys',
			controller: 'surveys',
			template: surveysTemplate
		},
		survey:{
			title: 'Survey',
			route: '/surveys/:id',
			controller: 'survey',
			template: surveyTemplate
		},
		policy:{
			title: 'Policy Privacy',
			route: '/policy',
			controller: 'policy',
			template: policyTemplate
		},
		terms:{
			title: 'Terms and Conditions',
			route: '/terms',
			controller: 'terms',
			template: termsTemplate
		}
	};
});