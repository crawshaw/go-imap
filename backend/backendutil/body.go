package backendutil

import (
	"bytes"
	"errors"
	"io"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-message"
)

var errNoSuchPart = errors.New("backendutil: no such message body part")

// FetchBodySection extracts a body section from a message.
func FetchBodySection(e *message.Entity, section *imap.BodySectionName) (imap.Literal, error) {
	// First, find the requested part using the provided path
	for i := 0; i < len(section.Path); i++ {
		n := section.Path[i]

		mr := e.MultipartReader()
		if mr == nil {
			if len(section.Path) == 1 && section.Path[0] == 1 {
				// RFC 3501 Section 6.4.5 states:
				//
				//	Every message has at least one part number.
				//
				// So treat a request for BODY[1] on a non-multipart
				// message as a request for the entire text body.
				break
			}
			return nil, errNoSuchPart
		}

		for j := 1; j <= n; j++ {
			p, err := mr.NextPart()
			if err == io.EOF {
				return nil, errNoSuchPart
			} else if err != nil {
				return nil, err
			}

			if j == n {
				e = p
				break
			}
		}
	}

	// Then, write the requested data to a buffer
	b := new(bytes.Buffer)

	// Write the header
	mw, err := message.CreateWriter(b, e.Header)
	if err != nil {
		return nil, err
	}
	defer mw.Close()

	switch section.Specifier {
	case imap.TextSpecifier:
		// The header hasn't been requested. Discard it.
		b.Reset()
	case imap.EntireSpecifier:
		if len(section.Path) > 0 {
			// When selecting a specific part by index, IMAP servers
			// return only the text, not the associated MIME header.
			b.Reset()
		}
	}

	// Write the body, if requested
	switch section.Specifier {
	case imap.EntireSpecifier, imap.TextSpecifier:
		if _, err := io.Copy(mw, e.Body); err != nil {
			return nil, err
		}
	}

	var l imap.Literal = b
	if section.Partial != nil {
		l = bytes.NewReader(section.ExtractPartial(b.Bytes()))
	}
	return l, nil
}
