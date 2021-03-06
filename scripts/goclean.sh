#!/bin/bash
# https://gist.github.com/hailiang/0f22736320abe6be71ce
# The script does automatic checking on a Go package and its sub-packages, including:
# 1. gofmt         (http://golang.org/cmd/gofmt/)
# 2. goimports     (https://github.com/bradfitz/goimports)
# 3. golint        (https://github.com/golang/lint)
# 4. go vet        (http://golang.org/cmd/vet)
# 5. race detector (http://blog.golang.org/race-detector)
# 6. test coverage (http://blog.golang.org/cover)

set -e

# Automatic checks
test -z "$(gofmt -l -w .     | tee /dev/stderr)"
#test -z "$(goimports -l -w . | tee /dev/stderr)"
#test -z "$(golint .             | tee /dev/stderr)"
#test -z "$(golint wt/.          | tee /dev/stderr)"
#test -z "$(golint handlers/.    | tee /dev/stderr)"
go vet ./...
echo 'Run tests'
#go test -race ./... #disabled for alpine go142
echo "PKG_CONFIG_PATH: $PKG_CONFIG_PATH"

godep go test ./...

# Run test coverage on each subdirectories and merge the coverage profile.

echo "mode: count" > profile.cov


# Standard go tooling behavior is to ignore dirs with leading underscors
for dir in $(find -L . -maxdepth 10 -not -path './.git*' -not -path '*/_*' -type d);
do
if ls $dir/*.go &> /dev/null; then
	godep go test -covermode=count -coverprofile=$dir/profile.tmp $dir
    if [ -f $dir/profile.tmp ]
    then
        cat $dir/profile.tmp | tail -n +2 >> profile.cov
        rm $dir/profile.tmp
    fi
fi
done

go tool cover -func profile.cov

# To submit the test coverage result to coveralls.io,
# use goveralls (https://github.com/mattn/goveralls)
# goveralls -coverprofile=profile.cov -service=travis-ci
