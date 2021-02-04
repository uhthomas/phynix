package templates

import (
	"bytes"
	"github.com/sipin/gorazor/gorazor"
	"phynix/templates/layout"
)

func Verify(text string, token string) string {
	var _buffer bytes.Buffer

	title := func() string {
		var _buffer bytes.Buffer

		_buffer.WriteString(("Verifcation"))

		return _buffer.String()
	}

	head := func() string {
		var _buffer bytes.Buffer

		_buffer.WriteString("<style type=\"text/css\">\n\n    .wrapper > .content .hero {\n\n      height: 200px;\n\n      background-image: url(/s/icon.png);\n\n      background-size: 128px 128px;\n\n      background-position: center;\n\n      background-repeat: no-repeat;\n\n    }\n\n\n\n    .wrapper > .content > .title {\n\n      font-size: 24px;\n\n    }\n\n\n\n    .wrapper > .content > .title > .text {\n\n      font-size: 18px;\n\n      margin-top: 10px;\n\n    }\n\n\n\n    .wrapper > .content > .title > .text a {\n\n      color: #2196f3;\n\n    }\n\n  </style>")

		_buffer.WriteString("<script type=\"text/javascript\">\n\n    (function() {\n\n      var token = \"")
		_buffer.WriteString(gorazor.HTMLEscape(token))
		_buffer.WriteString("\";\n\n      if (!token) return;\n\n      localStorage.setItem('token', token);\n\n      setTimeout(function() {\n\n        window.location = '/dashboard';\n\n      }, 5 * 1000);\n\n    })();\n\n  </script>")

		return _buffer.String()
	}

	content := func() string {
		var _buffer bytes.Buffer

		_buffer.WriteString("<div class=\"hero\"></div>")

		_buffer.WriteString("<div class=\"title\">\n\n    Verification\n\n    <div class=\"text\">")
		_buffer.WriteString((text))
		_buffer.WriteString("</div>\n\n  </div>")
		return _buffer.String()
	}

	return layout.Base(_buffer.String(), title(), head(), content())
}
