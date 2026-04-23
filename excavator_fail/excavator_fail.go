package fail

fail

/*
This is a non-compiling file that has been added to explicitly ensure that CI fails.
It also contains the command that caused the failure and its output.
Remove this file if debugging locally.

./godelw verify failed after updating godel plugins and assets

Command that caused error:
./godelw verify --skip-test --skip-lint

Output:
Running format...
Running mod...
Running generate...
Error: imports: exit status 1
failed to execute [/tmp/_conjure/conjure-4.51.0-extract-dir/conjure-4.51.0/bin/conjure compile imports/imports.yml /tmp/3122983396/out.json]
Output:
Error: LinkageError occurred while loading main class com.palantir.conjure.cli.ConjureCli
	java.lang.UnsupportedClassVersionError: com/palantir/conjure/cli/ConjureCli has been compiled by a more recent version of the Java Runtime (class file version 61.0), this version of the Java Runtime only recognizes class file versions up to 55.0

github.com/palantir/godel-conjure-plugin/v6/ir-gen-cli-bundler/conjureircli.RunWithParams
	/repo/vendor/github.com/palantir/godel-conjure-plugin/v6/ir-gen-cli-bundler/conjureircli/run.go:147
github.com/palantir/godel-conjure-plugin/v6/ir-gen-cli-bundler/conjureircli.InputPathToIRWithParams
	/repo/vendor/github.com/palantir/godel-conjure-plugin/v6/ir-gen-cli-bundler/conjureircli/run.go:79
github.com/palantir/godel-conjure-plugin/v6/ir-gen-cli-bundler/conjureircli.InputPathToIR
	/repo/vendor/github.com/palantir/godel-conjure-plugin/v6/ir-gen-cli-bundler/conjureircli/run.go:64
main.run
	/repo/integration_test/testgenerated/generate.go:54
main.main
	/repo/integration_test/testgenerated/generate.go:33
runtime.main
	/go/go-dists/go1.26.2/src/runtime/proc.go:290
runtime.goexit
	/go/go-dists/go1.26.2/src/runtime/asm_amd64.s:1771
exit status 1
generate.go:19: running "go": exit status 1
Error: failed to run go generate in "/repo/integration_test/testgenerated": exit status 1
Running license...
Running distgo-task...
Failed tasks:
	generate

*/
