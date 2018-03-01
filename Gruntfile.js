module.exports = function(grunt) {
    grunt.initConfig({
        concat: {
            options: {
                process: function(src, filepath) {
                    return "/*! " + filepath + " */\n" + src + "\n\n"
                },
                stripBanners: true,
            },
            app: {
                src: [
                    "web/static/bower/underscore/underscore-min.js",
                    "web/static/bower/jquery/dist/jquery.min.js",
                    "web/static/bower/highcharts/highcharts.js",
                    "web/static/bower/bootstrap/dist/js/bootstrap.min.js",
                    "web/static/bower/moment/moment.js",
                    "web/static/bower/timeago/jquery.timeago.js",
                    "web/static/js/src/app.js",
                ],
                dest: "web/static/js/bin/app.js",
            },
        },

        uglify: {
            app: {
                options: {
                    preserveComments: "some",
                    compress: false,
                    mangle: false,
                },
                src: ["web/static/js/bin/app.js"],
                dest: "web/static/js/bin/app.min.js",
            },
        },

        less: {
            style: {
                src: ["web/static/less/style.less"],
                dest: "web/static/css/style.css",
            },
        },

        cssmin: {
            style: {
                src: ["web/static/css/style.css"],
                dest: "web/static/css/style.min.css",
            },
        },
    })

    grunt.loadNpmTasks("grunt-contrib-concat")
    grunt.loadNpmTasks("grunt-contrib-uglify")
    grunt.loadNpmTasks("grunt-contrib-less")
    grunt.loadNpmTasks("grunt-contrib-cssmin")

    grunt.registerTask("default", ["less", "cssmin", "concat", "uglify"])
}
