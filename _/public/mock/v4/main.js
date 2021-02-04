// Wrap in anonymous function
// var app = angular.module('phynix');
// app.controller('main', ['$scope'], function($scope) {

// });

$('.icon-toggle').click(function() {
	$(this).toggleClass('active');
});

$('#header .drawer').click(function() {
  $('#drawer').toggleClass('open');
  $('#content').toggleClass('drawer-open');
});

$('#header .sidebar').click(function() {
  $('#sidebar').toggleClass('open');
  $('#content').toggleClass('sidebar-open');
});

$('#sidebar .footer .icon').click(function() {
  $('#sidebar .footer .icon').removeClass('active');
  $(this).addClass('active');
});

$('#content .meta .icon').click(function() {
  $(this).toggleClass('active');
});

$(function() {
  setTimeout(function() { $('#header .sidebar').click() }, 500);
  var msg = $($('#sidebar .messages .message')[0]);
  for (var i = 0; i < 10; i++) {
    $('#sidebar .messages').append(msg.clone());
  }
});

// (function($) {
// 	var module = angular.module('');
// })(jQuery);