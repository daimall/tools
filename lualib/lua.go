package lualib

import (
	"context"
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	lua "github.com/yuin/gopher-lua"
)

/* lua和golang 之间的交互
LUA：
   1、获取用例环境变量
   2、通过 golang 执行脚本
   3、通过 golang 记录日志
Golang：
   1、调用 lua 启动函数
*/

type lualib struct {
	timeout time.Duration             // 超时时间
	funcs   map[string]lua.LGFunction // golang 共享出来的函数库
	globals map[string]lua.LValue     // 全局变量
	method  string                    // 方法名
	params  []lua.LValue              // 主函数参数名
}

func NewLua() *lualib {
	return &lualib{
		timeout: 60 * time.Second,
		method:  "start",
	}
}

func (ll *lualib) SetTimeout(timeout time.Duration) *lualib {
	ll.timeout = timeout
	return ll
}

func (ll *lualib) SetFuncs(funcs map[string]lua.LGFunction) *lualib {
	ll.funcs = funcs
	return ll
}

func (ll *lualib) SetGlobals(globals map[string]lua.LValue) *lualib {
	ll.globals = globals
	return ll
}

func (ll *lualib) RunLuaScript(filePath string) (err error) {
	// 创建Lua虚拟机
	L := lua.NewState()
	defer L.Close()
	// 设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), ll.timeout)
	defer cancel()
	L.SetContext(ctx)

	// 将Go函数库注册到Lua虚拟机中
	for name, fn := range ll.funcs {
		L.Register(name, fn)
	}
	for k, v := range ll.globals {
		L.SetGlobal(k, v)
	}
	// 加载并执行Lua脚本
	if err := L.DoFile(filePath); err != nil {
		return fmt.Errorf("failed to load lua script: %v", err)
	}
	//先获取lua中定义的函数(用例入口函数)
	fn := L.GetGlobal(ll.method)
	cp := lua.P{
		Fn:      fn,
		NRet:    1,    //表示lua 启动函数有几个返回值，这里定义1个
		Protect: true, //lua脚本出错时，是否panic（false表示panic)
	}

	if err = L.CallByParam(cp, ll.params...); err != nil {
		logs.Error("Exec lua script(%s) failed, %s", filePath, err.Error())
		return err
	}
	lret := L.Get(-1) //获取最后一个返回值
	L.Remove(-1)      //删除最新的一个返回值
	if lret == lua.LBool(false) {
		//用例执行失败
		err = fmt.Errorf("exec lua script(%s) failed, ret is false", filePath)
		logs.Error(err.Error())
		return err
	}
	return nil
}
