package templates

import (
	"bytes"
	"github.com/sipin/gorazor/gorazor"
	"math"
	"phynix/models"
	"phynix/templates/layout"
)

func Dashboard(user models.User, communities []models.Community) string {
	var _buffer bytes.Buffer

	title := func() string {
		var _buffer bytes.Buffer

		_buffer.WriteString(("Dashboard"))

		return _buffer.String()
	}

	head := func() string {
		var _buffer bytes.Buffer

		_buffer.WriteString("<script src=\"https://cdnjs.cloudflare.com/ajax/libs/jquery/3.0.0-beta1/jquery.min.js\" type=\"text/javascript\"></script>")

		_buffer.WriteString("<style type=\"text/css\">\n\n    .loader {\n\n      height: 5px;\n\n      overflow: hidden;\n\n      position: relative;\n\n      background: rgba(255, 102, 0, 0.4);\n\n    }\n\n\n\n    .loader::before,\n\n    .loader::after {\n\n      content: '';\n\n      height: 5px;\n\n      position: absolute;\n\n      background: #FF6600;\n\n    }\n\n\n\n    .loader::before {\n\n      animation: loader-increase 2s infinite;\n\n    }\n\n\n\n    .loader::after {\n\n      animation: loader-decrease 2s 0.5s infinite;\n\n    }\n\n\n\n    @keyframes loader-increase {\n\n     from { left: -5%; width: 5%; }\n\n     to { left: 130%; width: 100%;}\n\n    }\n\n\n\n    @keyframes loader-decrease {\n\n     from { left: -80%; width: 80%; }\n\n     to { left: 110%; width: 10%;}\n\n    }\n\n\n\n    body::before {\n\n      content: '';\n\n      position: absolute;\n\n      top: 0;\n\n      left: 0;\n\n      width: 100%;\n\n      height: 400px;\n\n      background-image: linear-gradient(rgba(0, 0, 0, 0.4) 0%, rgba(0, 0, 0, 0.6) 75%, rgba(0, 0, 0, 0.8) 100%), url(/s/banner.jpg);\n\n      background-size: cover;\n\n      background-position-x: center;\n\n      background-repeat: no-repeat;\n\n    }\n\n\n\n    .wrapper {\n\n      max-width: 1586px;\n\n    }\n\n\n\n    @media all and (max-width: 1586px) {\n\n      .wrapper{\n\n        max-width: 1064px;\n\n      }\n\n    }\n\n\n\n    @media all and (max-width: 1064px) {\n\n      .wrapper{\n\n        max-width: 542px;\n\n      }\n\n    }\n\n\n\n    .wrapper > .content .hero {\n\n      width: 128px;\n\n      height: 128px;\n\n      background-image: linear-gradient(rgba(0,0,0,0.2), rgba(0,0,0,0.2)), url(/s/booty.jpg);\n\n      background-size: cover;\n\n      background-position: center;\n\n      display: inline-block;\n\n      margin-top: 36px;\n\n      border-radius: 50%;\n\n    }\n\n\n\n    .wrapper > .content .displayname {\n\n      font-size: 24px;\n\n      margin-bottom: 24px;\n\n    }\n\n\n\n    .wrapper > .content .displayname > .username {\n\n      font-size: 16px;\n\n      color: #808691;\n\n    }\n\n\n\n    .wrapper > .content .search {\n\n      margin: 5px;\n\n      background: #222937;\n\n      display: flex;\n\n    }\n\n\n\n    .wrapper > .content .search > .mdi-magnify {\n\n      line-height: 40px;\n\n      width: 50px;\n\n      flex-shrink: 0;\n\n    }\n\n\n\n    .wrapper > .content .search > input {\n\n      border: 0;\n\n      outline: 0;\n\n      height: 30px;\n\n      flex: 1;\n\n      color: white;\n\n      padding: 5px;\n\n      background: 0;\n\n    }\n\n\n\n    .wrapper > .content .filter {\n\n      display: flex;\n\n      line-height: 40px;\n\n      background: #222937;\n\n      margin: 5px;\n\n    }\n\n\n\n    .wrapper > .content .filter > .section {\n\n      flex: 1;\n\n      box-shadow: inset -1px 0px 0px 0px #0a0a0a;\n\n    }\n\n\n\n    .wrapper > .content .filter > .section:last-child {\n\n      box-shadow: none;\n\n    }\n\n\n\n    .wrapper > .content .filter > .section.active {\n\n      background: #03A9F4;\n\n      box-shadow: none;\n\n    }\n\n\n\n    .wrapper > .content .results {\n\n      overflow: hidden;\n\n    }\n\n\n\n    .wrapper > .content .results > .loader {\n\n      margin: 0px 5px;\n\n    }\n\n\n\n    .wrapper > .content .results > .community {\n\n      display: inline-block;\n\n      box-sizing: border-box;\n\n      background-size: cover;\n\n      background-position: center;\n\n      width: 512px;\n\n      height: 288px;\n\n      margin: 5px;\n\n      padding: 20px;\n\n      text-align: left;\n\n      position: relative;\n\n      border-radius: 2px;\n\n      float: left;\n\n    }\n\n\n\n    .wrapper > .content .results > .community > .title {\n\n      color: #fff;\n\n      font-size: 25px;\n\n      text-overflow: ellipsis;\n\n      white-space: nowrap;\n\n      overflow: hidden;\n\n    }\n\n\n\n    .wrapper > .content .results > .community > .artist {\n\n      color: #e0e0e0;\n\n      font-size: 15px;\n\n    }\n\n\n\n    .wrapper > .content .results > .community > .info {\n\n      bottom: 0;\n\n      left: 0;\n\n      position: absolute;\n\n      font-size: 22px;\n\n      display: block;\n\n      color: #e0e0e0;\n\n      width: 512px;\n\n      box-sizing: border-box;\n\n      padding: 20px;\n\n    }\n\n\n\n    .wrapper > .content .results > .community > .info > .name {\n\n      \n\n    }\n\n\n\n    .wrapper > .content .results > .community > .info > .population,\n\n    .wrapper > .content .results > .community > .info > .waitlist {\n\n      float: right;\n\n    }\n\n\n\n    .wrapper > .content .results > .community > .info > .population {\n\n      margin-left: 10px;\n\n    }\n\n\n\n    .wrapper > .content .results > .community > .info > .waitlist > .mdi,\n\n    .wrapper > .content .results > .community > .info > .population > .mdi {\n\n      margin-right: 10px;\n\n      color: #808691;\n\n    }\n\n  </style>")

		return _buffer.String()
	}

	content := func() string {
		var _buffer bytes.Buffer

		_buffer.WriteString("<div class=\"hero\"></div>")

		_buffer.WriteString("<div class=\"displayname\">\n\n    ")
		_buffer.WriteString(gorazor.HTMLEscape(user.Displayname))
		_buffer.WriteString("\n\n    <div class=\"username\">")
		_buffer.WriteString(gorazor.HTMLEscape(user.Username))
		_buffer.WriteString("</div>\n\n  </div>\n\n  <div class=\"search\">\n\n    <i class=\"mdi mdi-magnify\"></i>\n\n    <input placeholder=\"Search for communities\">\n\n  </div>\n\n  <div class=\"filter\">\n\n    <div class=\"section population active\">Population</div>\n\n    <div class=\"section favourites\">Favourites</div>\n\n  </div>\n\n  <div class=\"results\">\n\n    <!--<div class=\"loader\"></div>-->\n\n    ")

		for _, community := range communities[:int(math.Min(float64(len(communities)), 10))] {

			_buffer.WriteString("<a class=\"community\" href=\"/")
			_buffer.WriteString(gorazor.HTMLEscape(community.Slug))
			_buffer.WriteString("\" style=\"background-image: linear-gradient(rgba(0, 0, 0, 0.4) 0%, rgba(0, 0, 0, 0.6) 75%, rgba(0, 0, 0, 0.8) 100%), url(")
			_buffer.WriteString(gorazor.HTMLEscape(community.Media.Image))
			_buffer.WriteString(")\">\n\n          <div class=\"title\">")
			_buffer.WriteString(gorazor.HTMLEscape(community.Media.Title))
			_buffer.WriteString("</div>\n\n          <div class=\"artist\">")
			_buffer.WriteString(gorazor.HTMLEscape(community.Media.Artist))
			_buffer.WriteString("</div>\n\n          <div class=\"info\">\n\n            <span class=\"name\">")
			_buffer.WriteString(gorazor.HTMLEscape(community.Name))
			_buffer.WriteString("</span>\n\n            <span class=\"population\"><i class=\"mdi mdi-account-multiple\"></i>")
			_buffer.WriteString(gorazor.HTMLEscape(community.Population))
			_buffer.WriteString("</span>\n\n            <span class=\"waitlist\"><i class=\"mdi mdi-account-convert\"></i>")
			_buffer.WriteString(gorazor.HTMLEscape(community.Waitlist))
			_buffer.WriteString("</span>\n\n          </div>\n\n        </a>")

		}

		_buffer.WriteString("\n\n  </div>")
		return _buffer.String()
	}

	return layout.Base(_buffer.String(), title(), head(), content())
}
