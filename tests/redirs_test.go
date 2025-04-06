package tests

import (
	"fmt"
	"mvdan.cc/sh/v3/syntax"
	"shAudit/visitor"
	"strings"
	"testing"
)

func TestRedirsVisitor_Test(t *testing.T) {
	shCmd := []string{
		"echo 'Hello' > output.txt",                      // 基础输出重定向
		"cat < input.txt",                                // 输入重定向
		"command 2> error.log",                           // 错误重定向
		"command &> all.log",                             // 合并 stdout 和 stderr
		"exec 3> debug.txt; echo 'debug' >&3",            // 自定义文件描述符
		"cat <<EOF\nThis is a here doc\nEOF",             // Here Document
		"grep 'pattern' <<< 'test string'",               // Here String
		"echo 'Discarded' > /dev/null 2>&1",              // 特殊设备
		"LOGFILE=log_$(date).txt; echo 'log' > $LOGFILE", // 变量+命令替换
		"echo 'Error' > \"file with space.txt\"",         // 包含空格的文件名
		"echo 'Hello' > !(0-9)",                          // 排除模式
		"echo 'Hello' > $((x + 3))",                      // 数学表达式
		"echo 'Hello' > {a,'x',\"c\",^EX,$((x + z))}",    // 大括号扩展
		"cat log > >(grep 'error')",                      // 进程替换
	}
	for _, cmd := range shCmd {
		redirs := visitor.NewRedirsVisitor([]string{})
		redirs.RegisterLogger(&visitor.DefaultLogger{}, false)
		in := strings.NewReader(cmd)
		f, err := syntax.NewParser().Parse(in, "")
		if err != nil {
			t.Errorf("Failed to parse command: %v", err)
			return
		}
		syntax.Walk(f, redirs.Visit())
		_, err = redirs.Analyze()
		if err != nil {
			t.Errorf("Error analyzing command: %v", err)
		}
		for _, res := range redirs.GetFounded() {
			fmt.Printf("Found redir: %s\n", res)
		}
	}
}
