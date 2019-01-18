conjure-go
==========
[![](https://godoc.org/github.com/palantir/conjure-go?status.svg)](http://godoc.org/github.com/palantir/conjure-go)

Go generator for [Conjure](https://github.com/palantir/conjure).

Overview
--------
`conjure-go` takes a Conjure intermediate representation (IR) as input and writes the Go source files that implement
the types and services defined in the input IR.

Usage
-----
* `conjure-go [--output <output-dir>] input-ir-file`: writes the Go files for the Conjure IR file provided as input. 
  Uses the directory specified by `--output` as the base directory for writing the output (uses the working directory if
  unspecified).  

Update verification spec
------------------------
`conjure-go` tests its implementation using the specification defined by [`conjure-verification`](https://github.com/palantir/conjure-verification/).
The version used for verification is specified by the value of the `conjureVerifierVersion` constant in [conjure-go-verifier/generate.go].
To change/update the verifier version, change the constant and run `./godelw generate` -- this will download the test
cases and IR spec, regenerate the Conjure files used by the tests based on the new definition and change the version of
the Docker image used in the CircleCI configuration.

godel plugin
------------
[`godel-conjure-plugin`](https://github.com/palantir/godel-conjure-plugin) packages `conjure-go` as a
gödel plugin and provides abstractions to simplify the process of using `conjure-go` in a project.
