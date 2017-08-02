module.exports = function(grunt){

	var jsFilePaths = [
		'public/js/*.js'
	];

	grunt.initConfig({
		pkg: grunt.file.readJSON('package.json'),
		appDir: 'public',
		builtDir: 'public/dist',
		requirejs:{
			main:{
				options:{
					appDir: '<%= appDir %>',
					baseUrl: './js',
					dir: '<%= builtDir %>',
					optimizeCss: 'none', // will be taken care of with compass
					optimize: 'none' // will be taken care of with an uglify task
				}
			}
		},
		uglify:{
			options:{
				banner: '/*! <%= pkg.name %> <%= grunt.template.today("yyyy-mm-dd") %> */\n'
			},
			build:{
				files:(function(){
					var files = [];
					jsFilePaths.forEach(function(val){
						files.push({
							expand: true,
							cwd: '<%= builtDir %>',
							src: val,
							dest: '<%= builtDir %>'
						});
					});

					return false;
				})()
			}
		},
		jshint: {
			options:{
				reporter: require('jshint-stylish')
			},
			all:[
				'Gruntfile.js',
				'<%= appDir %>/js/{,*/}*.js'
			]
		},
		compass:{
			dist:{
				options: {
					sassDir: '<%= appDir %>/sass',
					cssDir: '<%= appDir %>/css',
					banner: '/*! <%= pkg.name %> <%= grunt.template.today("yyyy-mm-dd") %> */\n',
					specify: '<%= appDir %>/sass/*.scss',
					imagesDir: 'public/imgs',
					fontsDir: 'public/fonts',
					outputStyle:'compressed',
					noLineComments: true,
					environment: 'production'
				}
			},
			dev:{
				options: {
					sassDir: '<%= appDir %>/sass',
					cssDir: '<%= appDir %>/css',
					banner: '/*! <%= pkg.name %> <%= grunt.template.today("yyyy-mm-dd") %> */\n',
					specify: '<%= appDir %>/sass/*.scss',
					imagesDir: 'public/imgs',
					fontsDir: 'public/fonts',
					outputStyle:'expanded',
					noLineComments: false,
					environment: 'development'
				}
			}
		},
		watch:{
			scripts:{
				files: (function(){
					var files = [];
					jsFilePaths.forEach(function(val){
						files.push('<%= appDir %>/'+val);
					});
					return files;
				})(),
				tasks:['jshint'],
				options:{
					spawn: false
				}
			},
			compass:{
				files: '<%= appDir %>/sass/*.scss',
				tasks: ['compass:dev'],
				options:{
					spawn:false
				}
			},
			golang:{
				files: ['**/**/*.go'],
				tasks:['build-server:dev'],
				options:{
					livereload:true
				}
			}
		},
		concurrent:{
			prod:{
				tasks:['build-server:dev','watch:scripts','watch:compass','watch:golang']
			},
			options:{
				logConcurrentOutput: true
			}
		},
		'build-server':{
			dev:{
				root: '$GOPATH/src/github.com/curt-labs/GoSurvey'
			}
		}
	});

	grunt.loadNpmTasks('grunt-contrib-requirejs');
	// grunt.loadNpmTasks('jshint-stylish');
	// grunt.loadNpmTasks('requirejs');
	grunt.loadNpmTasks('grunt-contrib-jshint');
	grunt.loadNpmTasks('grunt-contrib-jasmine');
	grunt.loadNpmTasks('grunt-contrib-concat');
	grunt.loadNpmTasks('grunt-contrib-uglify');
	grunt.loadNpmTasks('grunt-contrib-watch');
	grunt.loadNpmTasks('grunt-contrib-compass');
	grunt.loadNpmTasks('grunt-shell');
	grunt.loadNpmTasks('grunt-concurrent');

	grunt.registerTask('default',['jshint','compass:dev','concurrent:prod']);

	grunt.registerTask('prod', ['jshint','requirejs','uglify','compass:dist']);

	grunt.registerMultiTask('build-server', "Compile all the things.", function() {
		var done = this.async(),
				path = require('path'),
				dir = grunt.config(['build-server', this.target, 'root']),
				opts = grunt.config(['build-server', this.target, 'opts']);

		return grunt.util.spawn({
			cmd: "go",
			args: ["run", "watcher.go"],
			opts:{stdio:'inherit'}
		}, function (error, result, code) {
			if (error) {
				grunt.log.error("Failed to build server: " + error + " (stdout='" + result.stdout + "', stderr='" + result.stderr + "')");
				done(false);
			} else {
				grunt.log.writeln("Now you can run " + path.resolve(dir, "../", "index") + " and go to http://0.0.0.0:8080");
				done();
			}
		});
	});
};
