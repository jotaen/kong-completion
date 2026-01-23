# Tab Completion for kong

`kong-completion` is a drop-in library that provides tab completion for CLI apps written with [kong](https://github.com/alecthomas/kong).

It currently supports the following shells:

- Bash
- Zsh
- Fish

`kong-completion` provides two main functionalities:

- It makes a kong app able to intercept and respond to tab completion requests. The completions are automatically derived from kong annotations. They can optionally be enhanced or adjusted with custom predictors.
- Since users have to manually activate the completion functionality in their shell, `kong-completion` provides a subcommand that instructs them how to achieve this.

## Get Started

See [the code of the sample app](./example/greet.go) for how to use the library.

In case you want to compile and run the demo app, keep in mind that completions only work for binaries in your $PATH, not for local ones (e.g. with `./` prefix). You also have to activate the completions first.

## API Reference

For flags and commands of your kong app, you can specify the following parameters in the annotation:

- `completion-enabled`
  - Whether this command or flag should be eligible for completions. By default, all flags and commands are eligible unless they are hidden. You can override this behaviour via this annotation parameter.
  - Possible values: `true`, `false`
  - Default value: derived from kong’s `hidden` flag – i.e., if the flag is hidden, it by default isn’t available for completion.
  - Usage example: `completion-enabled:"true"`

For the `Completion` subcommand specifically (as provided by this library), you can specify the following parameters in the annotation:

- `completion-shell-default`
  - Whether completions should fall back to the shell’s default completion behaviour, e.g. to complete file paths.
  - Possible values: `true`, `false`
  - Default value: `true`
  - Usage example: `completion-shell-default:"false"`

## About

`kong-completion` is free and open-source software, distributed under the [MIT license](./LICENSE.txt).

This library was originally based on [kongplete](https://github.com/WillAbides/kongplete).
