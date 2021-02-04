package templates

import (
	"bytes"
	"phynix/models"
	"phynix/templates/layout"
)

func Index(communities []models.Community) string {
	var _buffer bytes.Buffer

	title := func() string {
		var _buffer bytes.Buffer

		_buffer.WriteString(("Phynix"))

		return _buffer.String()
	}

	head := func() string {
		var _buffer bytes.Buffer

		_buffer.WriteString("<script src=\"/s/js/jquery.js\" type=\"text/javascript\"></script>")

		_buffer.WriteString("<script src='https://www.google.com/recaptcha/api.js'></script>")

		_buffer.WriteString("<style type=\"text/css\">\n\n    .loader {\n\n      height: 5px;\n\n      overflow: hidden;\n\n      position: relative;\n\n      background: rgba(255, 102, 0, 0.4);\n\n    }\n\n\n\n    .loader::before,\n\n    .loader::after {\n\n      content: '';\n\n      height: 5px;\n\n      position: absolute;\n\n      background: #FF6600;\n\n    }\n\n\n\n    .loader::before {\n\n      animation: loader-increase 2s infinite;\n\n    }\n\n\n\n    .loader::after {\n\n      animation: loader-decrease 2s 0.5s infinite;\n\n    }\n\n\n\n    @keyframes loader-increase {\n\n     from { left: -5%; width: 5%; }\n\n     to { left: 130%; width: 100%;}\n\n    }\n\n\n\n    @keyframes loader-decrease {\n\n     from { left: -80%; width: 80%; }\n\n     to { left: 110%; width: 10%;}\n\n    }\n\n\n\n    button {\n\n      cursor: pointer;\n\n      height: 64px;\n\n      line-height: 64px;\n\n      opacity: .6;\n\n      outline: 0;\n\n      border: 0;\n\n      background: 0;\n\n      color: white;\n\n      color: #FF6600;\n\n      font-size: 14px;\n\n      font-family: 'Roboto', sans-serif;\n\n      text-transform: uppercase;\n\n      -webkit-transition: all .2s ease;\n\n      -moz-transition: all .2s ease;\n\n      -ms-transition: all .2s ease;\n\n      -o-transition: all .2s ease;\n\n      transition: all .2s ease;\n\n    }\n\n\n\n    button:hover,\n\n    button:active {\n\n      box-shadow: rgba(0, 0, 0, 0.156863) 0px 3px 10px, rgba(0, 0, 0, 0.227451) 0px 3px 10px;\n\n      background: #323742;\n\n      opacity: 1;\n\n    }\n\n\n\n    body {\n\n      display: flex;\n\n    }\n\n\n\n    #content {\n\n      margin: auto;\n\n      width: 400px;\n\n      min-width: 400px;\n\n      box-sizing: border-box;\n\n      padding: 16px;\n\n      text-align: center;\n\n    }\n\n\n\n    #content::before {\n\n      content: '';\n\n      display: block;\n\n      height: 128px;\n\n      margin: 16px 0;\n\n      background-image: url(/s/icon.png);\n\n      background-size: 128px 128px;\n\n      background-position: center;\n\n      background-repeat: no-repeat;\n\n    }\n\n\n\n    #content .text {\n\n      height: 64px;\n\n      line-height: 64px;\n\n      margin: 8px 0;\n\n    }\n\n\n\n    #content .buttons {\n\n      display: flex;\n\n    }\n\n\n\n    #content .buttons button {\n\n      margin: auto;\n\n      width: 40%;\n\n    }\n\n\n\n    #content .buttons button[data-type=\"login\"]::before {\n\n      content: 'login';\n\n    }\n\n\n\n    #content .buttons button[data-type=\"signup\"]::before {\n\n      content: 'sign up';\n\n    }\n\n\n\n    .modal {\n\n      width: 100%;\n\n      height: 100%;\n\n      position: fixed;\n\n      background: rgba(0,0,0,0.4);\n\n      display: none;\n\n    }\n\n\n\n    .modal > .container {\n\n      margin: auto;\n\n      background: #282c35;\n\n      width: 350px;\n\n      text-align: center;\n\n      box-shadow: rgba(0, 0, 0, 0.247059) 0px 14px 45px, rgba(0, 0, 0, 0.219608) 0px 10px 18px;\n\n    }\n\n\n\n    .modal > .container > .header {\n\n      width: 100%;\n\n      height: 64px;\n\n      background: #323742;\n\n    }\n\n\n\n    .modal > .container > .header > .tab {\n\n      width: 50%;\n\n      height: 64px;\n\n      line-height: 64px;\n\n      display: inline-block;\n\n      float: left;\n\n      cursor: pointer;\n\n      text-transform: uppercase;\n\n    }\n\n\n\n    .modal > .container > .header > .tab.active {\n\n      background: #282c35;\n\n      cursor: default;\n\n    }\n\n\n\n    .modal > .container > .loader {\n\n      opacity: 0;\n\n      -webkit-transition: opacity .2s ease-in-out;\n\n      -ms-transition: opacity .2s ease-in-out;\n\n      transition: opacity .2s ease-in-out;\n\n    }\n\n\n\n    .modal > .container > .content {\n\n      width: 100%;\n\n      box-sizing: border-box;\n\n      padding: 0px 20px 20px 20px;\n\n      text-align: center;\n\n    }\n\n\n\n    .modal > .container > .content > .hero {\n\n      height: 200px;\n\n      background-image: url(/s/icon.png);\n\n      background-size: 128px 128px;\n\n      background-position: center;\n\n      background-repeat: no-repeat;\n\n    }\n\n\n\n    .modal > .container > .content > .section {\n\n      display: none;\n\n    }\n\n\n\n    .modal[data-type=\"login\"] > .container > .content > .section.login,\n\n    .modal[data-type=\"signup\"] > .container > .content > .section.signup {\n\n      display: block;\n\n    }\n\n\n\n    .modal > .container > .content > .section .tag {\n\n      margin-bottom: 20px;\n\n    }\n\n\n\n    .modal > .container > .content > .section .message {\n\n      padding: 20px;\n\n      margin-bottom: 20px;\n\n      display: none;\n\n    }\n\n\n\n    .modal > .container > .content > .section .message:first-letter {\n\n      text-transform: capitalize;\n\n    }\n\n\n\n    .modal > .container > .content > .section input {\n\n      outline: 0;\n\n      border: 0;\n\n      background: #323742;\n\n      box-sizing: border-box;\n\n      width: 100%;\n\n      padding: 10px;\n\n      margin-bottom: 10px;\n\n      color: white;\n\n      font-family: 'Roboto', sans-serif;\n\n    }\n\n\n\n    .modal > .container > .content > .section .combine > input {\n\n      width: 49%;\n\n      float: left;\n\n    }\n\n\n\n    .modal > .container > .content > .section .combine > input:last-child {\n\n      float: right;\n\n    }\n\n\n\n    .modal > .container > .content > .section .g-recaptcha {\n\n      display: inline-block;\n\n    }\n\n\n\n    .modal > .container > .content > .section button {\n\n      width: 100%;\n\n      margin-top: 10px;\n\n    }\n\n\n\n    .modal[data-disabled=\"true\"] > .container > .content > .section button {\n\n\n\n    }\n\n  </style>")

		return _buffer.String()
	}

	content := func() string {
		var _buffer bytes.Buffer
		_buffer.WriteString("<!-- --> ")
		return _buffer.String()
	}
	_buffer.WriteString("\n\n\n\n<script type=\"text/javascript\" id=\"weuihr\">$('.wrapper,#weuihr').remove()</script>\n\n<div id=\"content\">\n\n  <div class=\"text\">Work in progress. Probably coming soon.</div>\n\n  <div class=\"buttons\">\n\n    <button data-type=\"login\"></button>\n\n    <button data-type=\"signup\"></button>\n\n  </div>\n\n</div>\n\n\n\n<div class=\"modal\">\n\n  <div class=\"container\">\n\n    <div class=\"header\">\n\n      <div class=\"tab active\" data-type=\"login\">Login</div>\n\n      <div class=\"tab\" data-type=\"signup\">Sign up</div>\n\n    </div>\n\n    <div class=\"loader\"></div>\n\n    <div class=\"content\">\n\n      <div class=\"hero\"></div>\n\n      <div class=\"section login\">\n\n        <form id=\"login-form\">\n\n          <div class=\"tag\">Login to Phynix</div>\n\n          <div class=\"message\"></div>\n\n          <input class=\"email\" placeholder=\"Email\" maxlength=\"100\">\n\n          <input class=\"password\" placeholder=\"Password\" maxlength=\"72\" type=\"password\">\n\n          <button type=\"submit\">Login</button>\n\n        </form>\n\n      </div>\n\n      <div class=\"section signup\">\n\n        <form id=\"signup-form\">\n\n          <div class=\"tag\">Create an account</div>\n\n          <div class=\"message\"></div>\n\n          <div class=\"combine\">\n\n            <input class=\"displayname\" placeholder=\"Displayname\" maxlength=\"20\">\n\n            <input class=\"username\" placeholder=\"Username\" maxlength=\"20\" data-has-focussed=\"false\">\n\n          </div>\n\n          <input class=\"email\" placeholder=\"Email\" maxlength=\"100\">\n\n            <div class=\"combine\">\n\n            <input class=\"password\" placeholder=\"Password\" maxlength=\"72\" type=\"password\">\n\n            <input class=\"password-confirm\" placeholder=\"Confirm password\" maxlength=\"72\" type=\"password\">\n\n          </div>\n\n          <div class=\"g-recaptcha\" data-sitekey=\"6LeW_hYTAAAAACTBX2_85rU_Byvo-ycV4tiDtOmR\" data-theme=\"dark\"></div>\n\n          <button type=\"submit\">Sign up</button>\n\n        </form>\n\n      </div>\n\n    </div>\n\n  </div>\n\n</div>\n\n<script type=\"text/javascript\">\n\n  (function() {\n\n    var modal = $('.modal');\n\n    var loader = modal.find('.loader');\n\n    var loginSection = modal.find('.section.login');\n\n    var signupSection = modal.find('.section.signup');\n\n\n\n    function showMessage(section, type, message) {\n\n      var colors = {\n\n        'error': '#F44336',\n\n        'fail': '#F44336',\n\n        'ok': '#4CAF50',\n\n        'success': '#4CAF50'\n\n      }\n\n\n\n      var color = colors[type];\n\n\n\n      // don't kill me, using .html since all strings are safe anyway and <br> is useful sometimes\n\n      section.find('.message').html(message).css('color', color).show();\n\n    }\n\n\n\n    function hideMessage(section) {\n\n      sectiom.find('.message').css('display', 'none');\n\n    }\n\n\n\n    $('body').keydown(function(e) {\n\n      if (e.keyCode === 27) modal.fadeOut(200);\n\n    });\n\n\n\n    modal.click(function(e) {\n\n      if (e.target !== this) return;\n\n      $(this).fadeOut(200);\n\n    });\n\n\n\n    $('.buttons button, .modal .tab').click(function() {\n\n      modal.find('.tab').removeClass('active');\n\n      modal.find(`.tab[data-type=\"${$(this).attr('data-type')}\"]`).addClass('active');\n\n\n\n      modal.attr('data-type', $(this).attr('data-type'));\n\n      modal.fadeIn(200);\n\n    });\n\n\n\n    $('form').submit(function(e) {\n\n      e.preventDefault();\n\n    });\n\n\n\n    $('#login-form').submit(function(e) {\n\n      if (modal.attr('data-disabled') === 'true') return;\n\n\n\n      var payload = {\n\n        email: modal.find('.login .email').val(),\n\n        password: modal.find('.login .password').val()\n\n      }\n\n\n\n      modal.attr('data-disabled', 'true');\n\n      loader.css('opacity', '1');\n\n      $.ajax({\n\n        url: '/_/login',\n\n        method: 'POST',\n\n        dataType: 'json',\n\n        contentType: 'application/json',\n\n        data: JSON.stringify(payload)\n\n      })\n\n      .done(function(e) {\n\n        localStorage.setItem('token', e.data.token);\n\n        window.location = '/dashboard';\n\n      })\n\n      .fail(function(e, _, msg) {\n\n        showMessage(loginSection, 'fail', (e.responseJSON && e.responseJSON.error) || msg);\n\n      })\n\n      .always(function() {\n\n        loader.css('opacity', '0');\n\n        modal.attr('data-disabled', 'false');\n\n      });\n\n    });\n\n\n\n    modal.find('.signup .displayname').keydown(function(e) {\n\n      setTimeout(function() {\n\n        var username = $('.signup .username');\n\n        if (username.val() && username.attr('data-has-focussed') === 'true') return;\n\n\n\n        var displayname = $(this).val().toLowerCase();\n\n\n\n        var payload = displayname.replace(/[^A-Za-z0-9. _-]/g, '').trim().replace(/ /g, '_');\n\n\n\n        username.val(payload);\n\n      }.bind(this), 1);\n\n    });\n\n\n\n    modal.find('.signup .username').keydown(function(e) {\n\n      $(this).attr('data-has-focussed', true);\n\n    });\n\n\n\n    $('#signup-form').submit(function(e) {\n\n      if (modal.attr('data-disabled') === 'true') return;\n\n\n\n      var payload = {\n\n        displayname: modal.find('.signup .displayname').val().trim(),\n\n        username: modal.find('.signup .username').val().trim(),\n\n        email: modal.find('.signup .email').val().trim(),\n\n        password: modal.find('.signup .password').val(),\n\n        passwordConfirm: modal.find('.signup .password-confirm').val(),\n\n        captcha: grecaptcha.getResponse()\n\n      }\n\n\n\n      if (payload.password !== payload.passwordConfirm) {\n\n        showMessage(signupSection, 'error', 'Passwords do not match');\n\n        return;\n\n      }\n\n\n\n      delete payload.passwordConfirm;\n\n\n\n      var tests = {\n\n        displayname: /^.{2,20}$/,\n\n        username: /^[a-z0-9._-]{2,20}$/,\n\n        email: /^(?=.{5,100}$).+.+\\..+/,\n\n        password: /^.{2,72}$/\n\n      }\n\n\n\n      for (var i in tests) {\n\n        if (!tests[i].test(payload[i])) {\n\n          showMessage(signupSection, 'error', i + ' invalid')\n\n          return;\n\n        }\n\n      }\n\n\n\n      if (!payload.captcha) {\n\n        showMessage(signupSection, 'error', 'Captcha must be valid');\n\n        return;\n\n      }\n\n\n\n      modal.attr('data-disabled', 'true');\n\n      loader.css('opacity', '1');\n\n      $.ajax({\n\n        url: '/_/signup',\n\n        method: 'POST',\n\n        dataType: 'json',\n\n        contentType: 'application/json',\n\n        data: JSON.stringify(payload)\n\n      })\n\n      .done(function(e) {\n\n        if (e.data.token) {\n\n          localStorage.setItem('token', e.data.token);\n\n          window.location = '/dashboard';\n\n          return;\n\n        }\n\n        showMessage(signupSection, 'ok', 'Awesome!<br>Check your email to verify your account');\n\n      })\n\n      .fail(function(e, _, msg) {\n\n        grecaptcha.reset();\n\n        showMessage(signupSection, 'fail', e.responseJSON.error || msg);\n\n      })\n\n      .always(function() {\n\n        loader.css('opacity', '0');\n\n        modal.attr('data-disabled', 'false');\n\n      });\n\n    });\n\n    modal.css('display', 'flex').hide()\n\n  })();\n\n</script>")

	return layout.Base(_buffer.String(), title(), head(), content())
}
