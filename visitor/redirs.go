package visitor

import (
	"github.com/yxxchange/shAudit/utils"
	"mvdan.cc/sh/v3/syntax"
	"regexp"
	"strings"
)

// RedirsVisitor is a visitor that will find all redirections in the AST
type RedirsVisitor struct {
	rules []string
	reg   *regexp.Regexp

	founded []string

	logger        Logger
	verboseLogger Logger
}

func NewRedirsVisitor(rules []string) *RedirsVisitor {
	return &RedirsVisitor{
		rules:         rules,
		reg:           regexp.MustCompile(strings.Join(rules, "|")),
		logger:        &DefaultLogger{},
		verboseLogger: WrapVerboseLogger(&DefaultLogger{}, false),
	}
}

func (visit *RedirsVisitor) RegisterLogger(logger Logger, verbose bool) *RedirsVisitor {
	visit.logger = logger
	visit.verboseLogger = WrapVerboseLogger(logger, verbose)
	return visit
}

func (visit *RedirsVisitor) Visit() VisitorFunc {
	return func(node syntax.Node) bool {
		switch n := node.(type) {
		case *syntax.Redirect:
			visit.verboseLogger.Infof("Found Redirs: %s", SDebugf(n))
			var redirsArgs string
			var err error
			if n.Word != nil {
				redirsArgs, err = utils.WordPartToString(n.Word.Parts)
				if err != nil {
					visit.verboseLogger.Errorf("Error parsing redir: %v", err)
					return true
				}
			}
			visit.founded = append(visit.founded, redirsArgs)
		}
		return true
	}
}

func (visit *RedirsVisitor) Analyze() (result interface{}, err error) {
	if len(visit.founded) == 0 {
		return nil, nil
	}
	for _, redir := range visit.founded {
		if !visit.reg.MatchString(redir) {
			continue
		}
		visit.verboseLogger.Infof("Found redir: %s", redir)
		return redir, nil
	}
	visit.founded = nil
	return nil, nil
}

func (visit *RedirsVisitor) GetFounded() []string {
	return visit.founded
}
