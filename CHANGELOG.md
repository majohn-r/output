# Changelog

This project uses [semantic versioning](https://semver.org/); be aware that, until the major version becomes non-zero,
[this proviso](https://semver.org/#spec-item-4) applies.

Key to symbols

- ❗ breaking change
- 🐛 bug fix
- ⚠️ change in behavior, may surprise the user
- 😒 change is invisible to the user
- 🆕 new feature

## v0.9.4

_ release `2025-10-09`_

- 😒 updated build dependencies

## v0.9.3

_ release `2025-08-29`_

- 😒 updated build dependencies

## v0.9.2

_ release `2025-07-19`_

- 😒 updated build dependencies

## v0.9.1

_ release `2025-07-19`_

- 😒 updated build and code dependencies

## v0.9.0

_release `2024-08-27`_

- ❗remove `Bus` interface functions

    - `WriteCanonicalConsole(string, ...any)`
    - `WriteCanonicalError(string, ...any)`
    - `WriteConsole(string, ...any)`
    - `WriteError(string, ...any)`

## v0.8.0

_release `2024-08-26`_

- 🆕add new functions to the `Bus` interface:

  - `ConsolePrintln(string)`
  - `ConsolePrintf(string, ...any)`
  - `ErrorPrintln(string)`
  - `ErrorPrintf(string, ...any)`

- ⚠️deprecated these functions in the `Bus` interface:

  - `WriteCanonicalConsole(string, ...any)` (use `ConsolePrintln` or `ConsolePrintf` instead)
  - `WriteConsole(string, ...any)` (use `ConsolePrintln` or `ConsolePrintf` instead)
  - `WriteCanonicalError(string, ...any)` (use `ErrorPrintln` or `ErrorPrintf` instead)
  - `WriteError](string, ...any)` (use `ErrorPrintln` or `ErrorPrintf` instead)

## v0.7.0

_release `2024-08-25`_

- 🆕add support for lists (bulleted and numeric) in the console and error channels; adds new functions to the **Bus**
interface:

    - `BeginConsoleList(bool)`
    - `EndConsoleList()`
    - `ConsoleListDecorator() *ListDecorator`
    - `BeginErrorList(bool)`
    - `EndErrorList()`
    - `ErrorListDecorator() *ListDecorator`


## v0.6.0

_release `2024-08-24`_

- ⚠️add colon (**:**) to the set of sentence-terminating characters

## v0.5.4

_release `2024-08-21`_

- 😒no changes for consumers

## v0.5.3

_release `2024-06-29`_

- 😒no changes for consumers

## v0.5.2

_release `2024-06-13`_

- 🐛improved output of file and line number in `*Recorder.Report()` output

## v0.5.1

_release `2024-06-13`_

- 🆕add file and line number to `*Recorder.Report()` output

## v0.5.0

_release `2024-05-30`_

- 🆕add `Tab() uint8`, `IncrementTab(uint8)`, and `DecrementTab(uint8)` functions to `Bus` interface

## v0.4.0

_release `2024-05-23`_

- 🆕add `(r *Recorder) Report(TestingReporter, string, WantedRecording)` function

## v0.3.4

_release `2024-05-20`_

- 😒no changes for consumers

## v0.3.3

_release `2024-04-19`_

- 😒no changes for consumers

## v0.3.2

_release `2024-03-10`_

- 😒no changes for consumers

## v0.3.1

_release `2024-01-04`_

- 😒no changes for consumers

## v0.3.0

_release `2023-11-14`_

- 🆕add `IsConsoleTTY() bool` and `IsErrorTTY() bool` functions to the `Bus` interface

## v0.2.0

_release `2023-11-04`_

- ❗remove `LogWriter() Logger` function from the `Bus` interface

## v0.1.3

_release `2023-10-02`_

- 😒no changes for consumers

## v0.1.2

_release `2023-02-07`_

- 😒no changes for consumers

## v0.1.1

_release `2022-12-07`_

- 😒no changes for consumers

## v0.1.0

_release `2022-11-04`_

- 🆕initial release