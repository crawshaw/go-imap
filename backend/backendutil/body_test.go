package backendutil

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-message"
)

const plainHeader = "MIME-Version: 1.0\r\n" +
	"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
	"\r\n"

const plainBody = "Hello, World!"

const plainMsg = plainHeader + plainBody

var bodyTests = []struct {
	msg     string
	section string
	body    string
}{
	{
		msg:     plainMsg,
		section: "BODY[1]",
		body:    plainBody,
	},
	{
		section: "BODY[1.1]",
		body:    testTextBodyString,
	},
	{
		section: "BODY[1.2]",
		body:    testHTMLBodyString,
	},
	{
		section: "BODY[2]",
		body:    testAttachmentBodyString,
	},
	{
		section: "BODY[HEADER]",
		body:    testHeaderString,
	},
	{
		section: "BODY[1.1.HEADER]",
		body:    testTextHeaderString,
	},
	{
		section: "BODY[2.HEADER]",
		body:    testAttachmentHeaderString,
	},
	{
		section: "BODY[2.MIME]",
		body:    testAttachmentHeaderString,
	},
	{
		section: "BODY[TEXT]",
		body:    testBodyString,
	},
	{
		section: "BODY[1.1.TEXT]",
		body:    testTextBodyString,
	},
	{
		section: "BODY[2.TEXT]",
		body:    testAttachmentBodyString,
	},
	{
		section: "BODY[2.1]",
		body:    "",
	},
	{
		section: "BODY[4]",
		body:    "",
	},
	{
		section: "BODY[2.TEXT]<0.9>",
		body:    testAttachmentBodyString[:9],
	},
}

func TestFetchBodySection(t *testing.T) {
	for _, test := range bodyTests {
		test := test
		t.Run(test.section, func(t *testing.T) {
			msg := test.msg
			if msg == "" {
				msg = testMailString
			}
			e, err := message.Read(strings.NewReader(msg))
			if err != nil {
				t.Fatal("Expected no error while reading mail, got:", err)
			}

			section, err := imap.ParseBodySectionName(imap.FetchItem(test.section))
			if err != nil {
				t.Fatal("Expected no error while parsing body section name, got:", err)
			}

			r, err := FetchBodySection(e, section)
			if test.body == "" {
				if err == nil {
					t.Error("Expected an error while extracting non-existing body section")
				}
			} else {
				if err != nil {
					t.Fatal("Expected no error while extracting body section, got:", err)
				}

				b, err := ioutil.ReadAll(r)
				if err != nil {
					t.Fatal("Expected no error while reading body section, got:", err)
				}

				if s := string(b); s != test.body {
					t.Errorf("Expected body section %q to be \n%s\n but got \n%s", test.section, test.body, s)
				}
			}
		})
	}
}
