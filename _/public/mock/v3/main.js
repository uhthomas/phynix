// Wrap in anonymous function
// var app = angular.module('phynix');
// app.controller('main', ['$scope'], function($scope) {

// });

$('#header .meta').click(function() {
	$('#content').toggleClass('active');
	$('#sidebar').toggleClass('active');
	$(this).toggleClass('active');
});

$('#sidebar .footer .icon').click(function() {
	$('#sidebar .footer .icon').removeClass('active');
	$(this).addClass('active');
});

// (function($) {
// 	var module = angular.module('');
// })(jQuery);