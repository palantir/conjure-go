package fail

fail

/*
This is a non-compiling file that has been added to explicitly ensure that CI fails.
It also contains the command that caused the failure and its output.
Remove this file if debugging locally.

go mod operation failed. This may mean that there are legitimate dependency issues with the "go.mod" definition in the repository and the updates performed by the gomod check. This branch can be cloned locally to debug the issue.

Command that caused error:
./godelw exec -- go get github.com/bodgit/sevenzip github.com/pierrec/lz4/v4@v4.1.25

Output:
go: github.com/bodgit/sevenzip@upgrade (v1.6.4) requires github.com/pierrec/lz4/v4@v4.1.26, not github.com/pierrec/lz4/v4@v4.1.25
Error: exit status 1

*/
