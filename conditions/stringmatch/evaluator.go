package stringmatch

import (
	"errors"
	"strings"
	"text/scanner"
)

const (
	AND = "AND" // 逻辑与操作符
	OR  = "OR"  // 逻辑或操作符
)

// Calculate 函数解析并计算给定的布尔表达式字符串。
// 参数:
//
//	str: 待计算的布尔表达式字符串，可以包含布尔值、逻辑操作符（AND、OR）、括号和标识符。
//	stackSize: 布尔值和操作符栈的最大大小。
//	boolFunc: 一个用于从标识符获取布尔值的回调函数。
//
// 返回值:
//
//	r: 计算结果的布尔值。
//	b: 如果在计算过程中出现错误，则为一个描述错误的字符串；如果没有错误，则为nil。
func Calculate(str string, stackSize int, boolFunc func(string) bool) (r bool, b error) {
	defer func() {
		if err := recover(); err != nil {
			b = errors.New(err.(string))
			return
		}
	}()

	// 创建用于存储布尔值的栈和操作符的栈
	var boolStack = newStack(true, stackSize)
	var symbolStack = newStack("", stackSize)

	// 初始化字符串扫描器
	var s scanner.Scanner
	s.Mode = scanner.ScanIdents
	s.Init(strings.NewReader(str))
	tok := s.Scan()
	tt := s.TokenText()

	// 开始解析表达式
	for tok != scanner.EOF {
		switch tok {
		case scanner.Ident:
			if tt == AND || tt == OR {
				if symbolStack.IsEmpty() {
					symbolStack.Push(tt)
				} else {
					symbol := symbolStack.Peek()
					if tt == "(" {
						symbolStack.Push(tt)
					}
					if getPriority(tt) < getPriority(symbol) {
						symbolStack.Push(tt)
					} else {
						symbol := symbolStack.Pop()
						bool1 := boolStack.Pop()
						bool2 := boolStack.Pop()
						ret := getSum(bool1, bool2, symbol)
						boolStack.Push(ret)
						symbolStack.Push(tt)
					}
				}
			} else {
				boolStack.Push(boolFunc(tt))
			}
		case '(':
			symbolStack.Push("(")
		case ')':
			// 处理右括号，执行相应计算
			for {
				symbol := symbolStack.Pop()
				if symbol != "(" {
					bool1 := boolStack.Pop()
					bool2 := boolStack.Pop()
					ret := getSum(bool1, bool2, symbol)
					boolStack.Push(ret)
				} else {
					break
				}
			}
		default:
			boolStack.Push(boolFunc(tt))
		}
		tok, tt = s.Scan(), s.TokenText()
	}

	// 执行剩余的计算
	for symbolStack.Len() != 0 {
		if boolStack.Len() < 2 {
			return false, errors.New("err7 表达式错误")
		}
		symbol := symbolStack.Pop()
		data1 := boolStack.Pop()
		data2 := boolStack.Pop()
		ret := getSum(data1, data2, symbol)
		boolStack.Push(ret)
	}

	// 验证计算结果并返回
	if boolStack.Len() != 1 {
		return false, errors.New("err10 表达式有误")
	}
	ret := boolStack.Peek()
	return ret, nil
}

// getSum 根据给定的操作符执行布尔运算。
func getSum(cond1 bool, cond2 bool, symbol string) bool {
	switch symbol {
	case AND:
		return cond1 && cond2
	case OR:
		return cond1 || cond2
	}
	return false
}

// getPriority 获取操作符的优先级，用于决定是否需要计算。
func getPriority(symbol string) int {
	switch symbol {
	case AND:
		return 3
	case OR:
		return 2
	case "(":
		return 8
	case ")":
		return 1
	default:
		return 0
	}
}
