package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	kongcompletion "github.com/jotaen/kong-completion"
	"github.com/posener/complete"
)

func main() {
	// 1. Create kong app, but don’t run arg parsing yet.
	app := kong.Must(&GreetingApp{})

	// 2. Register completions. This must happen before the parsing step, so that
	// tab completion invocations can be intercepted.
	// In this case, we also configure a custom predictor, but that’s optional.
	kongcompletion.Register(app, predictNames)

	// 3. Now, proceed as usual with parsing arguments and running the app.
	ctx, err := app.Parse(os.Args[1:])
	if err != nil {
		app.Printf("%s", err)
		app.Exit(1)
		return
	}
	err = ctx.Run()
	if err != nil {
		app.Printf("%s", err)
		app.Exit(1)
		return
	}
}

// GreetingApp is a sample app. The `Completion` subcommand is optional; it’s purpose is
// to help the user to configure the completions for their shell. (The completions themselves
// would work without this command, though.)
type GreetingApp struct {
	Hello      Hello                     `cmd:"" help:"Prints a greeting"`
	Completion kongcompletion.Completion `cmd:"" help:"Outputs shell code for initialising tab completions"`
}

type Hello struct {
	Strong bool   `help:"Print the greeting in ALL CAPS!!!"`
	Casual bool   `help:"In case you are more laid back, you know"`
	Name   string `arg:"" help:"Personalize the greeting with your name!" predictor:"name" default:"World"`
}

// Custom predictor. (Just as demo how these works.)
var predictNames = kongcompletion.WithPredictor(
	"name",
	complete.PredictSet("Ben", "Liz", "Mark", "Sarah"),
)

func (h *Hello) Run() error {
	formula := "Hello"
	if h.Casual {
		formula = "Howdy"
	}
	phrase := formula + " " + h.Name
	if h.Strong {
		phrase = strings.ToUpper(phrase) + "!!!1"
	}
	fmt.Println(phrase)
	return nil
}
