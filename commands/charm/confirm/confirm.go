package confirm

import (
	"github.com/digitalocean/doctl/commands/charm"
	"github.com/digitalocean/doctl/commands/charm/template"
	"github.com/erikgeiser/promptkit/confirmation"
)

type Choice string

const (
	Yes       Choice = "yes"
	No        Choice = "no"
	Undecided Choice = "undecided"
)

func toChoice(v *bool) Choice {
	if v == nil {
		return Undecided
	}

	if *v {
		return Yes
	} else {
		return No
	}
}

func fromChoice(v Choice) confirmation.Value {
	switch v {
	case Yes:
		return confirmation.Yes
	case No:
		return confirmation.No
	default:
		return confirmation.Undecided
	}
}

type Prompt struct {
	text          string
	choice        Choice
	displayResult DisplayResult
}

type Option func(*Prompt)

func WithDefaultChoice(c Choice) Option {
	return func(p *Prompt) {
		p.choice = c
	}
}

type DisplayResult int

const (
	DisplayResultNormal DisplayResult = iota
	DisplayResultEphemeral
	DisplayResultEphemeralYes
	DisplayResultEphemeralNo
)

func WithDisplayResult(v DisplayResult) Option {
	return func(p *Prompt) {
		p.displayResult = v
	}
}

func New(text string, opts ...Option) *Prompt {
	p := &Prompt{
		text: text,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

var promptTemplate = `
{{- highlight promptPrefix }} {{ .Prompt -}}
{{ if .YesSelected -}}
	{{- print (bold (print " " pointerRight "yes ")) " no" -}}
{{- else if .NoSelected -}}
	{{- print "  yes " (bold (print pointerRight "no")) -}}
{{- else -}}
	{{- "  yes  no" -}}
{{- end -}}
`

var resultTemplate = `
{{- if RenderResult .FinalValue -}}
{{- if.FinalValue -}}{{success promptPrefix}}{{else}}{{error promptPrefix}}{{end}}
{{- print " " .Prompt " " -}}
	{{- if .FinalValue -}}
		{{- success "yes" -}}
	{{- else -}}
		{{- error "no" -}}
	{{- end }}
{{- end -}}
`

func (p *Prompt) Prompt() (Choice, error) {
	input := confirmation.New(p.text, fromChoice(p.choice))
	tfs := template.Funcs(charm.Colors)
	tfs["RenderResult"] = func(finalValue bool) bool {
		switch p.displayResult {
		case DisplayResultEphemeralNo:
			return finalValue
		case DisplayResultEphemeralYes:
			return !finalValue
		default:
			return true
		}
	}
	input.ExtendedTemplateFuncs = tfs
	input.Template = promptTemplate
	input.ResultTemplate = resultTemplate

	v, err := input.RunPrompt()
	if err != nil {
		if err.Error() == "no decision was made" {
			return Undecided, err
		}
		return "", nil
	}
	return toChoice(&v), nil
}