define([], function(){
	var ranger = [function(){
		return function(input, total){
			total = parseInt(total, 0);
			for(var i = 0; i < total; i++){
				input.push(i);
			}
			return input;
		};
	}];

	return ranger;
});