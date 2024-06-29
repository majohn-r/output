package main

import (
	"github.com/goyek/goyek/v2"
	build "github.com/majohn-r/tools-build"
)

const coverageFile = "coverage.out"

var (
	clean = goyek.Define(goyek.Task{
		Name:  "clean",
		Usage: "delete build products",
		Action: func(a *goyek.A) {
			build.Clean([]string{coverageFile})
		},
	})

	_ = goyek.Define(goyek.Task{
		Name:  "coverage",
		Usage: "run unit tests and produce a coverage report",
		Action: func(a *goyek.A) {
			build.GenerateCoverageReport(a, coverageFile)
		},
	})

	_ = goyek.Define(goyek.Task{
		Name:  "doc",
		Usage: "generate documentation",
		Action: func(a *goyek.A) {
			build.GenerateDocumentation(a, []string{"build"})
		},
	})

	format = goyek.Define(goyek.Task{
		Name:  "format",
		Usage: "clean up source code formatting",
		Action: func(a *goyek.A) {
			build.Format(a)
		},
	})

	lint = goyek.Define(goyek.Task{
		Name:  "lint",
		Usage: "run the linter on source code",
		Action: func(a *goyek.A) {
			build.Lint(a)
		},
	})

	nilaway = goyek.Define(goyek.Task{
		Name:  "nilaway",
		Usage: "run nilaway on source code",
		Action: func(a *goyek.A) {
			build.NilAway(a)
		},
	})

	tests = goyek.Define(goyek.Task{
		Name:  "tests",
		Usage: "run unit tests",
		Action: func(a *goyek.A) {
			build.UnitTests(a)
		},
	})

	updateDependencies = goyek.Define(goyek.Task{
		Name:  "updateDependencies",
		Usage: "update dependencies",
		Action: func(a *goyek.A) {
			build.UpdateDependencies(a)
		},
	})

	vulnCheck = goyek.Define(goyek.Task{
		Name:  "vulnCheck",
		Usage: "run vulnerability check on source code",
		Action: func(a *goyek.A) {
			build.VulnerabilityCheck(a)
		},
	})

	_ = goyek.Define(goyek.Task{
		Name:  "preCommit",
		Usage: "run all pre-commit tasks",
		Deps:  goyek.Deps{clean, updateDependencies, lint, nilaway, format, vulnCheck, tests},
	})
)
