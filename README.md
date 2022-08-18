# Tab Completion for kong

This library provides tab completion for the [CLI library kong](https://github.com/alecthomas/kong).

It currently supports the following shells:

- Bash
- Zsh
- Fish

This library provides two main functionalities:

- It makes a kong app able to intercept tab completion requests. The completions are automatically generated from the annotations for subcommands and flags. Completions can optionally be enhanced or adjusted with custom predictors.
- It provides a subcommand that helps users to activate the program’s tab completions in their shell.

## Get Started

See [the code of the sample app](./example/greet.go) for how to use the library.

In case you want to compile and run the demo app, keep in mind that completions only work for binaries in your $PATH, not for local ones (e.g. with `./` prefix).

## About

`kong-completion` is free and open-source software, distributed under the [MIT license](./LICENSE.txt).

This library was originally based on [kongplete](https://github.com/WillAbides/kongplete). The functionality is the same, but the APIs are designed differently.
