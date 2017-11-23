package htmltemplateelement

import "github.com/matthewmueller/golly/js"

// js:"HTMLTemplateElement,omit"
type HTMLTemplateElement struct {
	window.HTMLElement
}

// Content
func (*HTMLTemplateElement) Content() (content *window.DocumentFragment) {
	js.Rewrite("$<.Content")
	return content
}
