// Package multicheker provides functionality to run a set of static analyzers.
//
// This package allows combining several analyzers into one executable file to perform
// simultaneous code checking for compliance with various rules and recommendations.
package multicheker

import (
	"github.com/Azzonya/go-shortener/internal/analyzers/exit"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

// Run starts a set of static analyzers.
//
// The Run function collects and adds static analyzers to the list of analyzers,
// including a set of staticcheck, simple, quickfix and stylecheck analyzers, as well as the
// the exit.New() analyzer, which checks for os.Exit() calls in the main() function.
// Then the function runs multichecker.Main with the obtained set of analyzers.
func Run() {
	checks := append(getStaticCheckAnalyzers(),
		nilfunc.Analyzer,
		exit.New())

	multichecker.Main(checks...)
}

// getStaticCheckAnalyzers returns a set of static analyzers to execute.
//
// The getStaticCheckAnalyzers function collects all available analyzers from the packages
// staticcheck, simple, quickfix and stylecheck and then returns them as a slice
// pointers to *analysis.Analyzer.
func getStaticCheckAnalyzers() []*analysis.Analyzer {
	checks := make([]*analysis.Analyzer, 0, len(staticcheck.Analyzers))
	for _, v := range staticcheck.Analyzers {
		checks = append(checks, v.Analyzer)
	}
	for _, v := range simple.Analyzers {
		checks = append(checks, v.Analyzer)
	}
	for _, v := range quickfix.Analyzers {
		checks = append(checks, v.Analyzer)
	}
	for _, v := range stylecheck.Analyzers {
		checks = append(checks, v.Analyzer)
	}

	return checks
}
