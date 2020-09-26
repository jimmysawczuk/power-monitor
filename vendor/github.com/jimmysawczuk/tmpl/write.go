package main

import (
	"bytes"
	"encoding/json"
	html "html/template"
	"io"
	text "text/template"

	"github.com/pkg/errors"
	"github.com/tdewolff/minify"
	htmlminify "github.com/tdewolff/minify/html"
)

func writeHTML(in string, o payload, min bool, out io.Writer) error {
	tmpl, err := html.New("output").Funcs(o.tmplfuncs()).Parse(in)
	if err != nil {
		return errors.Wrap(err, "compile template")
	}

	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, o); err != nil {
		return errors.Wrap(err, "execute template")
	}

	if min {
		by := buf.Bytes()
		buf.Reset()

		m := minify.New()
		hm := htmlminify.DefaultMinifier
		hm.KeepDocumentTags = true

		hm.Minify(m, buf, bytes.NewReader(by), nil)
	}

	io.Copy(out, buf)

	return nil
}

func writeJSON(in string, o payload, min bool, out io.Writer) error {
	tmpl, err := text.New("output").Funcs(o.tmplfuncs()).Parse(in)
	if err != nil {
		errors.Wrap(err, "compile template")
	}

	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, o); err != nil {
		errors.Wrap(err, "execute template")
	}

	dst := &bytes.Buffer{}
	if min {
		if err := json.Compact(dst, buf.Bytes()); err != nil {
			return errors.Wrap(err, "json: compact")
		}
	} else {
		if err := json.Indent(dst, buf.Bytes(), "", "    "); err != nil {
			return errors.Wrap(err, "json: indent")
		}
	}

	io.Copy(out, dst)

	return nil
}

func writeText(in string, o payload, out io.Writer) error {
	tmpl, err := text.New("output").Funcs(o.tmplfuncs()).Parse(in)
	if err != nil {
		return errors.Wrap(err, "compile template")
	}

	if err := tmpl.Execute(out, o); err != nil {
		return errors.Wrap(err, "execute template")
	}

	return nil
}
