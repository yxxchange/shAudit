package tests

import (
	"mvdan.cc/sh/v3/syntax"
	"shAudit/visitor"
	"strings"
	"testing"
)

func TestCommandVisitor(t *testing.T) {
	shllCmds := []string{
		"ps -ef | grep hello",
		"ps -ef | grep 'process' | awk '{print $2}'",
		"VAR=123; echo $VAR",
		"echo 'hello world'",
		"echo \"hello world\"",
		"echo hello world",
		"ls -l /tmp",
		"cat <<EOF\nThis is a here document.\nEOF",
		"ssh user@host \"ls -l /remote/path\"",
		"arr=(a b c); echo ${arr[@]}",
		"myfunc() { echo \"Hello from $1\"; }\nmyfunc \"World\"",
	}
	for _, cmd := range shllCmds {
		in := strings.NewReader(cmd)
		f, err := syntax.NewParser().Parse(in, "")
		if err != nil {
			t.Errorf("Failed to parse command: %v", err)
			return
		}

		rules := []string{
			"ssh",
			".*wk",
		}
		v := visitor.NewCommandVisitor(rules).RegisterLogger(&visitor.DefaultLogger{}, false)
		syntax.Walk(f, v.Visit())
		_, err = v.Analyze()
		if err != nil {
			t.Errorf("Error analyzing command: %v", err)
		}
	}
}
