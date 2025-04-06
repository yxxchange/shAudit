package utils

import (
	"fmt"
	"mvdan.cc/sh/v3/syntax"
	"strings"
)

func WordPartToString(parts []syntax.WordPart) (string, error) {
	redirsArgs := make([]string, 0)
	for _, part := range parts {
		stringer := NewStringer(part)
		if stringer == nil {
			return "", fmt.Errorf("unsupported type %T", part)
		}
		str, err := stringer.String()
		if err != nil {
			return "", err
		}
		if str != "" {
			redirsArgs = append(redirsArgs, str)
		}
	}
	return strings.Join(redirsArgs, " "), nil
}

type Stringer interface {
	String() (string, error)
}

func NewStringer(node syntax.Node) Stringer {
	switch n := node.(type) {
	case *syntax.Lit:
		return &Lit{
			lit: n,
		}
	case *syntax.SglQuoted:
		return &SglQuoted{
			sgl: n,
		}
	case *syntax.DblQuoted:
		return &DblQuoted{
			dbl: n,
		}
	case *syntax.ParamExp:
		return &ParamExp{
			param: n,
		}
	case *syntax.CmdSubst:
		return &CmdSubst{
			cmdSubst: n,
		}
	case *syntax.ArithmExp:
		return &ArithmExp{
			arithm: n,
		}
	case *syntax.ProcSubst:
		return &ProcSubst{
			procSubst: n,
		}
	case *syntax.ExtGlob:
		return &ExtGlob{
			extGlob: n,
		}
	default:
		return nil
	}
}

type Lit struct {
	lit *syntax.Lit
}

func (l Lit) String() (string, error) {
	if l.lit == nil {
		return "", nil
	}
	return l.lit.Value, nil
}

type ParamExp struct {
	param *syntax.ParamExp
}

func (p ParamExp) String() (string, error) {
	if p.param == nil {
		return "", nil
	}
	if p.param.Short {
		return fmt.Sprintf("$%s", p.param.Param.Value), nil
	}
	return fmt.Sprintf("${%s}", p.param.Param.Value), nil
}

type SglQuoted struct {
	sgl *syntax.SglQuoted
}

func (s SglQuoted) String() (string, error) {
	if s.sgl == nil {
		return "", nil
	}
	return fmt.Sprintf("'%s'", s.sgl.Value), nil
}

type DblQuoted struct {
	dbl *syntax.DblQuoted
}

func (d DblQuoted) String() (string, error) {
	if d.dbl == nil {
		return "", nil
	}
	content, err := WordPartToString(d.dbl.Parts)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("\"%s\"", content), nil
}

type CmdSubst struct {
	cmdSubst *syntax.CmdSubst
}

func (c CmdSubst) String() (string, error) {
	if c.cmdSubst == nil {
		return "", nil
	}
	// TODO: handle cmdSubst
	return fmt.Sprintf("$( SUBCOMMAND )"), nil
}

type ArithmExp struct {
	arithm *syntax.ArithmExp
}

func (a ArithmExp) String() (string, error) {
	// 检查表达式是否存在
	if a.arithm.X == nil {
		return "", nil
	}

	// 根据语法形式生成前缀
	var prefix string
	if a.arithm.Bracket {
		prefix = "$["
	} else {
		prefix = "$(("
		if a.arithm.Unsigned {
			prefix += "# " // 添加无符号标志（如 mksh 的 $((# expr)）
		}
	}

	// 生成表达式字符串（递归处理 ArithmExpr）
	exprStr := " EXPRESSION " // TODO: 处理表达式

	// 根据语法形式生成后缀
	var suffix string
	if a.arithm.Bracket {
		suffix = "]"
	} else {
		suffix = "))"
	}

	return prefix + exprStr + suffix, nil
}

type ProcSubst struct {
	procSubst *syntax.ProcSubst
}

func (p ProcSubst) String() (string, error) {
	if p.procSubst == nil {
		return "", nil
	}

	if p.procSubst.Op == syntax.CmdIn {
		return fmt.Sprintf("<( PROCSUBST )"), nil
	}
	if p.procSubst.Op == syntax.CmdOut {
		return fmt.Sprintf(">( PROCSUBST )"), nil
	}
	return "", fmt.Errorf("unsupported proc subst type: %T", p.procSubst)
}

type ExtGlob struct {
	extGlob *syntax.ExtGlob
}

func (e ExtGlob) String() (string, error) {
	if e.extGlob == nil {
		return "", nil
	}
	if e.extGlob.Op == syntax.GlobZeroOrOne {
		return fmt.Sprintf("?(%s)", e.extGlob.Pattern.Value), nil
	}
	if e.extGlob.Op == syntax.GlobZeroOrMore {
		return fmt.Sprintf("*(%s)", e.extGlob.Pattern.Value), nil
	}
	if e.extGlob.Op == syntax.GlobOneOrMore {
		return fmt.Sprintf("+(%s)", e.extGlob.Pattern.Value), nil
	}
	if e.extGlob.Op == syntax.GlobOne {
		return fmt.Sprintf("@(%s)", e.extGlob.Pattern.Value), nil
	}
	if e.extGlob.Op == syntax.GlobExcept {
		return fmt.Sprintf("!(%s)", e.extGlob.Pattern.Value), nil
	}
	return "", fmt.Errorf("unsupported ext glob type: %T", e.extGlob)
}

type BraceExp struct {
	braceExp *syntax.BraceExp
}

func (b BraceExp) String() (string, error) {
	if b.braceExp == nil {
		return "", nil
	}

	arr := make([]string, 0)
	for _, part := range b.braceExp.Elems {
		if str, err := WordPartToString(part.Parts); err == nil {
			arr = append(arr, str)
		} else {
			return "", err
		}
	}

	return fmt.Sprintf("{%s}", strings.Join(arr, ",")), nil
}
