package parser

import (
	"errors"
	"fmt"
	"strings"
)

type FuncParser struct {
	statement        string
	variablesToValue map[string]interface{}
}

func NewFuncParser(statement string) *FuncParser {
	return &FuncParser{
		statement: statement,
	}
}

// func extractVariableValues(statement string) map[string]interface{} {
// 	reg, err := regexp.Compile("{*}")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	ids := reg.FindAllString(statement, -1)
// 	variablesToValue := map[string]interface{}
// 	for _, id := range ids {
// 		// 通过id找value ...
// 	}
// 	return variablesToValue
// }

// func (f *FuncParser) Exec() interface{} {

// }

func (f *FuncParser) Exec() (interface{}, error) {
	// TODO: 对函数语句进行解析，提取其中的变量值

	return f.exec(f.statement)
}

const (
	FuncAND   = "AND"
	FuncFALSE = "FALSE"
	FuncIF    = "IF"
	FuncIFS   = "IFS"
	FuncNOT   = "NOT"
	FuncOR    = "OR"
	FuncTRUE  = "TRUE"
	FuncXOR   = "XOR"
)

func (f *FuncParser) exec(statement string, vals ...interface{}) (interface{}, error) {
	remove := func(statement, funcName string) string {
		statement = strings.TrimPrefix(statement, funcName)
		statement = strings.TrimPrefix(statement, "(")
		statement = strings.TrimPrefix(statement, ")")
		return statement
	}

	switch {
	case strings.HasPrefix(statement, FuncAND):
		statement := remove(statement, FuncAND)
		args := f.parseArgs(statement)
		for _, arg := range args {
			res, err := f.exec(arg)
			if err != nil {
				return nil, err
			}
			v, ok := res.(bool)
			if !ok {
				return nil, errors.New("AND()只能接受返回值为布尔类型的表达式")
			}
			if v == false {
				return false, nil
			}
		}
		return true, nil
	case strings.HasPrefix(statement, FuncFALSE):
		statement := remove(statement, FuncFALSE)
		if len(statement) != 0 {
			return nil, errors.New("FALSE()函数错误，不能存在参数")
		}
		return false, nil
	case strings.HasPrefix(statement, FuncIF):
		statement := remove(statement, FuncIF)
		args := f.parseArgs(statement)
		if len(args) != 3 {
			return nil, errors.New("IF()函数错误，参数有误")
		}
		result, err := f.exec(args[0])
		if err != nil {
			return nil, err
		}
		v, ok := result.(bool)
		if !ok {
			return nil, errors.New("IF()函数的第一个表达式返回值不为布尔类型")
		}
		if v {
			return args[1], nil
		}
		return args[2], nil
	case strings.HasPrefix(statement, FuncIFS):
		statement := remove(statement, FuncIFS)
		args := f.parseArgs(statement)
		if len(args)%2 != 0 {
			return nil, errors.New("IFS()函数的参数个数有误")
		}
		for i := 0; i < len(args); i += 2 {
			result, err := f.exec(args[i])
			if err != nil {
				return nil, err
			}
			v, ok := result.(bool)
			if !ok {
				return nil, fmt.Errorf("IFS()的第%d个表达式返回值不为布尔类型", i+1)
			}
			if v {
				return args[i+1], nil
			}
		}
	default:
		// TODO: 解析表达式
		return nil, nil
	}
}

func (f *FuncParser) parseArgs(statement string) []string {
	var res []string
	var start, end int
	var pass bool = false
	for i := 0; i < len(statement); i++ {
		if pass {
			continue
		}
		if statement[i] == ',' {
			end = i
			res = append(res, statement[start:end])
			start = end + 1
		}
		if statement[i] == '(' {
			pass = true
		}
		if statement[i] == ')' {
			pass = false
		}
	}
	return res
}
