package mail

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"text/template"

	gomail "gopkg.in/gomail.v2"
)

var (
	workers  = 10
	sendChan = make(chan *gomail.Message, workers)
	dialer   = gomail.NewPlainDialer("smtp.gmail.com", 587, "thomas@phynix.io", "bsuxspalyoovxbyh")
)

func init() {
	for i := 0; i < workers; i++ {
		go sendLoop()
	}
}

func Send(message *gomail.Message) {
	go func() {
		sendChan <- message
	}()
}

func SendTemplate(message *gomail.Message, name string, data interface{}) error {
	content, err := ioutil.ReadFile(filepath.Join("_", "template", "mail", name+".tmpl"))
	if err != nil {
		return err
	}

	t, err := template.New(name).Parse(fmt.Sprintf("%s", content))
	if err != nil {
		return err
	}

	var d bytes.Buffer

	if err := t.Execute(&d, data); err != nil {
		return err
	}

	message.SetBody("text/html; charset=utf-8", d.String())

	go func() {
		sendChan <- message
	}()

	return nil
}

func sendLoop() {
	for message := range sendChan {
		if err := dialer.DialAndSend(message); err != nil {
			fmt.Println(err.Error())
		}
	}
}
