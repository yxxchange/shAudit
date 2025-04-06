package main

import (
	"shAudit/audit"
	"shAudit/visitor"
)

func main() {
	shllCmds := []string{
		"echo 'log' > ${LOGFILE}",
	}
	auditor := audit.NewAudit([]visitor.IVisitor{
		visitor.NewRedirsVisitor([]string{"echo", ".*log"}),
		visitor.NewCommandVisitor([]string{"echo", ".*wk"}),
	})

	for _, cmd := range shllCmds {
		err := auditor.Audit(cmd)
		if err != nil {
			println("Audit error:", err.Error())
		} else {
			println("Audit success")
		}
	}
}
