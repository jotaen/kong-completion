module github.com/jotaen/kong-completion

go 1.20

require (
	github.com/alecthomas/kong v0.7.1
	github.com/posener/complete v1.2.3
	github.com/riywo/loginshell v0.0.0-20200815045211-7d26008be1ab
	github.com/stretchr/testify v1.7.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	v1.0.2 // Published accidentally.
	v1.0.1 // Published accidentally.
	v1.0.0 // Published accidentally.
)
