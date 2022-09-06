# Visual Cron

`visualcron` helps to visualise a cron expression

## Usage

1 argument is required, a string representation of a cron expression

```
$ visualcron "*/15 0 1,15 * 1-5 /usr/bin/find"
```

This will output the expression in table format

```
minute        0 15 30 45
hour          0
day of month  1 15
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   1 2 3 4 5
command       /usr/bin/find
```

## Development

For local development, [Go](http://golang.org) must be installed

Note: A minimum Go version of 1.13.0 is required

## Makefile

A makefile has been generated to ease testing, building and running `visualcron`

The makefile has 3 depenencies

- `make` (yum, apt, chocolatey) to run the make commands 
- `go` (http://golang.org) to run the go related commands
- `gox` (go get github.com/mitchellh/gox) to build binaries

The available commands are

- `test` - run the unit tests and print coverage
- `test-cover` - run the unit test and generate HTML coverage
- `build` - build the binaries for Linux, Mac, and Windows
- `run` - run
- `fmt` - run "go fmt"

## Building

[gox](https://github.com/mitchellh/gox) is used to build cross compatible versions of `visualcron`

For example, the following will build a Linux, Mac and Windows binary in the `builds/` directory

```shell
go get -v github.com/mitchellh/gox
gox -verbose -osarch="linux/amd64 darwin/amd64 windows/amd64" -output "builds/visualcron_{{.OS}}_{{.Arch}}"
```

## Testing

The unit tests can be executed with 

```
go test -v ./...
```

Unit test coverage is currently at 92.4%

## Author

Michael Bell