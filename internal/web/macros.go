package web

import (
	"html/template"

	"golang.org/x/text/message"
)

func FormatInt(value int) string {
	printer := message.NewPrinter(message.MatchLanguage("en"))

	return printer.Sprintf("%d", value)
}

func FormatFloat(value float64) string {
	printer := message.NewPrinter(message.MatchLanguage("en"))

	return printer.Sprintf("%0.2f", value)
}

var DefaultMacros = template.FuncMap{
	"FmtInt":   FormatInt,
	"FmtFloat": FormatFloat,
}
