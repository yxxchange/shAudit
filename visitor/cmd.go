package visitor

import (
	"fmt"
	"mvdan.cc/sh/v3/syntax"
	"regexp"
	"strings"
)

type CommandVisitor struct {
	rules []string
	reg   *regexp.Regexp

	founded []string

	logger        Logger
	verboseLogger Logger
}

func NewCommandVisitor(rules []string) *CommandVisitor {
	return &CommandVisitor{
		rules:         rules,
		reg:           regexp.MustCompile(strings.Join(rules, "|")),
		logger:        &DefaultLogger{},
		verboseLogger: WrapVerboseLogger(&DefaultLogger{}, false),
	}
}

func (visit *CommandVisitor) RegisterLogger(logger Logger, verbose bool) *CommandVisitor {
	visit.logger = logger
	visit.verboseLogger = WrapVerboseLogger(logger, verbose)
	return visit
}

func (visit *CommandVisitor) Visit() VisitorFunc {
	return func(node syntax.Node) bool {
		switch n := node.(type) {
		case *syntax.CallExpr:
			visit.verboseLogger.Infof("Found CallExpr: %s", SDebugf(n))
			if len(n.Args) == 0 {
				visit.verboseLogger.Infof("no args in CallExpr")
				return true
			}
			firstWord := n.Args[0] // first word
			if len(firstWord.Parts) == 0 {
				visit.verboseLogger.Infof("no parts in first word")
				return true
			}
			firstPart := firstWord.Parts[0]
			if lit, ok := firstPart.(*syntax.Lit); ok {
				visit.verboseLogger.Infof("Found cmd: %s", lit.Value)
				visit.founded = append(visit.founded, lit.Value)
			}
			return true
		}
		return true
	}
}

func (visit *CommandVisitor) Analyze() (result interface{}, err error) {
	if len(visit.founded) == 0 {
		return nil, nil
	}
	for _, cmd := range visit.founded {
		if visit.reg.MatchString(cmd) {
			return nil, fmt.Errorf("found command: %s", cmd)
		}
	}
	visit.founded = nil
	return nil, nil
}
