module.exports = function(grunt)
{
	grunt.initConfig({
		concat: {
			options: {
				process: function(src, filepath)
				{
					return "/*! " + filepath + " */\n" + src + "\n\n";
				}
			},
			app: {
				src: [
					'src/web/static/bower/underscore/underscore-min.js',
					'src/web/static/bower/jquery/dist/jquery.min.js',
					'src/web/static/bower/highcharts/highcharts.js',
					'src/web/static/bower/bootstrap/dist/js/bootstrap.min.js',
					'src/web/static/bower/moment/moment.js',
					'src/web/static/bower/timeago/jquery.timeago.js',
					'src/web/static/js/src/app.js'
				],
				dest: 'src/web/static/js/bin/app.js',
			}
		},

		less: {
			style: {
				src: ['src/web/static/less/style.less'],
				dest: 'src/web/static/css/style.css'
			}
		},

		cssmin: {
			style: {
				src: ['src/web/static/css/style.css'],
				dest: 'src/web/static/css/style.min.css'
			}
		}
	});

	grunt.loadNpmTasks('grunt-contrib-concat');
	grunt.loadNpmTasks('grunt-contrib-less');
	grunt.loadNpmTasks('grunt-contrib-cssmin');

	grunt.registerTask('default', ['less', 'cssmin', 'concat']);
};
