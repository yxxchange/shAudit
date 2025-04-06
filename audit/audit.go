package audit

import (
	"fmt"
	"github.com/yxxchange/shAudit/visitor"
	"mvdan.cc/sh/v3/syntax"
	"strings"
)

type Audit struct {
	Visitors []visitor.IVisitor
}

func NewAudit(visitors []visitor.IVisitor) *Audit {
	return &Audit{
		Visitors: visitors,
	}
}

func (a *Audit) Audit(cmd string) error {
	in := strings.NewReader(cmd)
	f, err := syntax.NewParser().Parse(in, "")
	if err != nil {
		return fmt.Errorf("shell command is error: %s", err.Error())
	}
	for _, v := range a.Visitors {
		syntax.Walk(f, v.Visit())
		_, err = v.Analyze()
		if err != nil {
			return err
		}
	}
	return nil
}
