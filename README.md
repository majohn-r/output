# output

[![GoDoc Reference](https://godoc.org/github.com/majohn-r/output?status.svg)](https://pkg.go.dev/github.com/majohn-r/output)
[![go.mod](https://img.shields.io/github/go-mod/go-version/majohn-r/output)](go.mod)
[![LICENSE](https://img.shields.io/github/license/majohn-r/output)](LICENSE)

[![Release](https://img.shields.io/github/release/majohn-r/output.svg)](https://github.com/majohn-r/output/releases)
[![Code Coverage Report](https://codecov.io/github/majohn-r/output/branch/main/graph/badge.svg)](https://codecov.io/github/majohn-r/output)
[![Go Report Card](https://goreportcard.com/badge/github.com/majohn-r/output)](https://goreportcard.com/report/github.com/majohn-r/output)
[![Build Status](https://img.shields.io/github/actions/workflow/status/majohn-r/output/build.yml?branch=main)](https://github.com/majohn-r/output/actions?query=workflow%3Abuild+branch%3Amain)

- [output](#output)
  - [Installing](#installing)
  - [Basic Usage](#basic-usage)
  - [Documentation](#documentation)
  - [Contributing](#contributing)
    - [Git](#git)
    - [Code Quality](#code-quality)
    - [Commit message](#commit-message)

**output** is a Go library that provides an easy way for command-line oriented
programs to handle console writing, error writing, and logging (but agnostic as
to the choice of logging framework). It also provides a simple way to verify
what is written to those writers.

## Installing

Execute this:

```text
go get github.com/majohn-r/output
```

## Basic Usage

In main, create a **Bus** implementation and a **Logger** implementation. Here is an example that uses the
[https://github.com/sirupsen/logrus](https://github.com/sirupsen/logrus) library
to implement logging:

```go
func main() {
    // the Bus created by output.NewDefaultBus() neither knows nor cares about
    // how logging actually works - that's the purview of the Logger
    // implementation it uses.
    o := output.NewDefaultBus(ProductionLogger{})
    runProgramLogic(o, os.Args)
}

func runProgramLogic(o output.Bus, args []string) {
    // any functions called should have the Bus passed in if they, or any
    // function they call, needs to write output or do any logging
    o.ConsolePrintf("hello world: %v\n", args)
}

type ProductionLogger struct {}

// Trace outputs a trace log message
func (ProductionLogger) Trace(msg string, fields map[string]any) {
    logrus.WithFields(fields).Trace(msg)
}

// Debug outputs a debug log message
func (ProductionLogger) Debug(msg string, fields map[string]any) {
    logrus.WithFields(fields).Debug(msg)
}

// Info outputs an info log message
func (ProductionLogger) Info(msg string, fields map[string]any) {
    logrus.WithFields(fields).Info(msg)
}

// Warning outputs a warning log message
func (ProductionLogger) Warning(msg string, fields map[string]any) {
    logrus.WithFields(fields).Warning(msg)
}

// Error outputs an error log message
func (ProductionLogger) Error(msg string, fields map[string]any) {
    logrus.WithFields(fields).Error(msg)
}

// Panic outputs a panic log message and calls panic()
func (ProductionLogger) Panic(msg string, fields map[string]any) {
    logrus.WithFields(fields).Panic(msg)
}

// Fatal outputs a fatal log message and terminates the program
func (ProductionLogger) Fatal(msg string, fields map[string]any) {
    logrus.WithFields(fields).Fatal(msg)
}
```

In the test code, the output can be checked like this:

```go
func Test_runProgramLogic {
    tests := map[string]struct {
        name string
        args []string
        output.WantedRecording
    }{
        "test case": {
            args: []string{"hi" "12" "true"},
            WantedRecording: output.WantedRecording{
                Console: "hello world: [hi 12 true]",
            },
        },
    }
    for name, tt := range tests {
        t.Run(name, func(t *testing.T) {
            o := NewRecorder()
            runProgramLogic(o, tt.args)
            if issues, ok := o.Verify(tt.WantedRecording); !ok {
                for _, issue := range issues {
                    t.Errorf("runProgramLogic() %s", issue)
                }
            }
        })
    }
}
```

## Documentation

Documentation beyond this file can be obtained by running `./build.sh doc`, or
go here:
[https://pkg.go.dev/github.com/majohn-r/output](https://pkg.go.dev/github.com/majohn-r/output)

## Contributing

### Git

1. Fork the repository (`https://github.com/majohn-r/output/fork`).
2. Create a feature branch (`git checkout -b my-new-feature`).
3. Commit your changes (`git commit -am 'Add some feature'`).
4. Push to the branch (`git push origin my-new-feature`).
5. Create a new Pull Request.

### Code Quality

These are the minimum standards:

1. run [`./build.sh preCommit`] with no errors and 100% coverage on tests.
2. update CHANGELOG.md with a brief description of the change(s).

### Commit message

Reference an issue in the first line of the commit message:

```text
(#1234) fix that nagging problem
```

In the example above, **1234** is the issue number this commit reference.

This library adheres to [Semantic Versioning](https://semver.org/) standards, so
it will be very helpful if the details in the commit message make clear whether
the changes require a minor or major release bump.
