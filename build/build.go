package main

import (
	"github.com/goyek/goyek/v3"
	build "github.com/majohn-r/tools-build"
)

const (
	coverageFile           = "coverage.out"
	taskClean              = "clean"
	taskCoverage           = "coverage"
	taskDeadCode           = "deadcode"
	taskDoc                = "doc"
	taskFix                = "fix"
	taskFormat             = "format"
	taskLint               = "lint"
	taskNilAway            = "nilaway"
	taskTests              = "tests"
	taskUpdateDependencies = "updateDependencies"
	taskVulnerabilityCheck = "vulnCheck"
)

var (
	clean = goyek.Define(goyek.Task{
		Name:  taskClean,
		Usage: "delete build products",
		Action: func(a *goyek.A) {
			if !build.TaskDisabled(taskClean) {
				build.Clean([]string{coverageFile})
			}
		},
	})

	_ = goyek.Define(goyek.Task{
		Name:  taskCoverage,
		Usage: "run unit tests and produce a coverage report",
		Action: func(a *goyek.A) {
			if !build.TaskDisabled(taskCoverage) {
				build.GenerateCoverageReport(a, coverageFile)
			}
		},
	})

	_ = goyek.Define(goyek.Task{
		Name:  taskDeadCode,
		Usage: "run deadcode analysis",
		Action: func(a *goyek.A) {
			if !build.TaskDisabled(taskDeadCode) {
				build.Deadcode(a)
			}
		},
	})

	_ = goyek.Define(goyek.Task{
		Name:  taskDoc,
		Usage: "generate documentation",
		Action: func(a *goyek.A) {
			if !build.TaskDisabled(taskDoc) {
				build.GenerateDocumentation(a, []string{"build"})
			}
		},
	})

	_ = goyek.Define(goyek.Task{
		Name:  taskFix,
		Usage: "run go fix",
		Action: func(a *goyek.A) {
			if !build.TaskDisabled(taskFix) {
				build.GoFix(a)
			}
		},
	})

	format = goyek.Define(goyek.Task{
		Name:  taskFormat,
		Usage: "clean up source code formatting",
		Action: func(a *goyek.A) {
			if !build.TaskDisabled(taskFormat) {
				build.Format(a)
			}
		},
	})

	lint = goyek.Define(goyek.Task{
		Name:  taskLint,
		Usage: "run the linter on source code",
		Action: func(a *goyek.A) {
			if !build.TaskDisabled(taskLint) {
				build.Lint(a)
			}
		},
	})

	nilaway = goyek.Define(goyek.Task{
		Name:  taskNilAway,
		Usage: "run nilaway on source code",
		Action: func(a *goyek.A) {
			if !build.TaskDisabled(taskNilAway) {
				build.NilAway(a)
			}
		},
	})

	tests = goyek.Define(goyek.Task{
		Name:  taskTests,
		Usage: "run unit tests",
		Action: func(a *goyek.A) {
			if !build.TaskDisabled(taskTests) {
				build.UnitTests(a)
			}
		},
	})

	updateDependencies = goyek.Define(goyek.Task{
		Name:  taskUpdateDependencies,
		Usage: "update dependencies",
		Action: func(a *goyek.A) {
			if !build.TaskDisabled(taskUpdateDependencies) {
				build.UpdateDependencies(a)
			}
		},
	})

	vulnCheck = goyek.Define(goyek.Task{
		Name:  taskVulnerabilityCheck,
		Usage: "run vulnerability check on source code",
		Action: func(a *goyek.A) {
			if !build.TaskDisabled(taskVulnerabilityCheck) {
				build.VulnerabilityCheck(a)
			}
		},
	})

	_ = goyek.Define(goyek.Task{
		Name:  "preCommit",
		Usage: "run all pre-commit tasks",
		Deps:  goyek.Deps{clean, updateDependencies, lint, nilaway, format, vulnCheck, tests},
	})
)
