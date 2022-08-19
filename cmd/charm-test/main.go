package main

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/digitalocean/doctl/commands/charm"
	"github.com/digitalocean/doctl/commands/charm/confirm"
	"github.com/digitalocean/doctl/commands/charm/input"
	"github.com/digitalocean/doctl/commands/charm/template"
	"github.com/digitalocean/doctl/commands/charm/textbox"
)

func main() {
	confirm.New("proceed?", confirm.WithDefaultChoice(confirm.Yes)).Prompt()
	i := input.New("app name:", input.WithRequired())
	i.Prompt()

	fmt.Println(
		charm.Checkmark, charm.CheckmarkSuccess,
	)

	fmt.Println(
		charm.TextSuccess.WithString("woo!"), charm.TextSuccess.S("woo 2!"),
	)

	if err := template.PrintE(heredoc.Doc(`
		--- template ---
		This is an example template.

		Another line.

		{{ success "maybe some success output" }}
		{{ success checkmark }} just the checkmark.
		{{ success (join " " (checkmark) "good job!") }}
		{{ error (join " " (checkmark) "we're both confused.") }}
		{{ warning (print promptPrefix " try again?") }}
		{{ error (join " " (crossmark) "there we go.") }}

		{{ success (bold "full send let's go!!!!") }}
		{{ bold (success "full send let's go!!!!") }}
		
		{{ bold (underline "underline behaves very strangely") }}
		{{ underline (bold "underline behaves very strangely") }}
		
		{{ success (underline "underline behaves very strangely") }}
		{{ underline (success "underline behaves very strangely") }}

		{{ newTextBox.Success.S "i'm in a box!" }}
	`), nil); err != nil {
		panic(err)
	}

	img := "yeet/yote:dev"
	dur := 23*time.Minute + 37*time.Second
	fmt.Fprintf(
		textbox.New().Success(),
		"%s Successfully built %s in %s",
		charm.CheckmarkSuccess,
		charm.TextSuccess.S(img),
		charm.TextWarning.S(dur.Truncate(time.Second).String()),
	)

	if err := template.BufferedE(textbox.New().Success(), heredoc.Doc(`
		{{ success checkmark }} Successfully built {{ success  .img }} in {{ warning (duration .dur) }}`,
	), map[string]any{
		"img": img,
		"dur": dur,
	}); err != nil {
		panic(err)
	}

	template.Buffered(textbox.New().Success(), heredoc.Doc(`
		{{ success checkmark }} Successfully built {{ success  .img }} in {{ warning (duration .dur) }}`,
	), map[string]any{
		"img": img,
		"dur": dur,
	})
}