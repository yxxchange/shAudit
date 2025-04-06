package visitor

import (
	"encoding/json"
	"fmt"
	"mvdan.cc/sh/v3/syntax"
)

type VisitorFunc func(node syntax.Node) bool

// Visitor All visitors should implement this interface to range over the shell AST
type IVisitor interface {
	// Visit returns a function that will be called for each node in the AST
	Visit() VisitorFunc

	// Analyze returns the result of the analysis
	// should collect info from the AST by the method Visit
	Analyze() (result interface{}, err error)
}

type Logger interface {
	// Infof Log logs the message
	Infof(format string, args ...interface{})
	// Errorf logs the error message
	Errorf(format string, args ...interface{})
}

type DefaultLogger struct{}

func (d *DefaultLogger) Infof(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
func (d *DefaultLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

type VerboseLogger struct {
	CustomLogger Logger

	Verbose bool
}

func WrapVerboseLogger(logger Logger, verbose bool) *VerboseLogger {
	return &VerboseLogger{
		CustomLogger: logger,
		Verbose:      verbose,
	}
}

func (v *VerboseLogger) Infof(format string, args ...interface{}) {
	if !v.Verbose {
		return
	}
	format = fmt.Sprintf("[INFO] %s\n", format)
	if v.CustomLogger != nil {
		v.CustomLogger.Infof(format, args...)
	} else {
		fmt.Printf(format, args...)
	}
}
func (v *VerboseLogger) Errorf(format string, args ...interface{}) {
	if !v.Verbose {
		return
	}
	format = fmt.Sprintf("[ERROR] %s\n", format)
	if v.CustomLogger != nil {
		v.CustomLogger.Errorf(format, args...)
	} else {
		fmt.Printf(format, args...)
	}
}

func SDebugf(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%+v", v)
	}
	return string(b)
}
