package rule

import (
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"runtime"
)

type Filter struct {
	ID        string `json:"Id"`
	When string `json:"When"`
	Name      string `json:"Name"`
}
type Action struct {
	ID        string `json:"Id"`
	When string `json:"When"`
	Then    string `json:"Then"`
	Name      string `json:"Name"`
	Action    []*Action
}

type ConfigValue struct {
	Id 		int
	Filter []*Filter
	Action []*Action
}

//go 异常处理
func DefPanicFun( name string ) {
	if e := recover(); e != nil {
		buf := make([]byte, 1024)
		buf = buf[:runtime.Stack(buf, false)]
		fmt.Printf("[PANIC] express:%v|%v\n%s \n",name ,  e, buf)
	}
}
//词法分析结果 bool缓存
var gBoolExprMap = map[string]*vm.Program{ }
//词法分析结果 缓存
var gExprMap = map[string]*vm.Program{ }

// DoRule 执行规则 
func DoRule( val *ConfigValue, env map[string]interface{} ){

	//参数克隆，确保参数不影响
	params :=  make(map[string]interface{}, len( env ))
	for k,v :=range env{
		params[k] =v
	}
	//过滤器
	for _,v:=range val.Filter {
		rst,err := doFilter( v, params  )
		if  err != nil {
			fmt.Printf("doFilter fail err:%v \n ", err.Error() )
			return
		}
		if !rst {
			return
		}
	}
	//fmt.Println("--------Action---------")
	//执行
	for _,v:=range val.Action {
		rst, err:=doAction( v , params );
		if  err!= nil {
			fmt.Printf("[error] doAction faile err:%v \n  ", err.Error() )
			continue
		}
		fmt.Printf("doAction rst %v \n  ", rst  )
	}
	//fmt.Println("--------rule end ---------")
}
//执行
func doExpress( express string, params map[string]interface{} )( interface{}, error) {
	defer DefPanicFun( express ) //注意这里的 异常处理
	if len( express ) <=0 {
		return false , nil
	}
	var cmdEx *vm.Program
	if f , ok  := gExprMap[express]; ok {
		cmdEx  =f
	} else {
		f, err := expr.Compile(express , expr.Env(params) ,expr.AllowUndefinedVariables())
		if err!= nil {
			return nil, err
		}
		gExprMap[express] = f //?? cache 一下
		cmdEx  =f
	}
	return expr.Run( cmdEx, params )
}
//执行 bool
func doBoolExpress( express string, params map[string]interface{} )( interface{}, error) {
	defer DefPanicFun( express ) //注意这里的 异常处理
	if len( express ) <=0 {
		return false , nil
	}
	var cmdEx *vm.Program
	if f , ok  := gBoolExprMap[express]; ok {
		cmdEx  =f
	} else{
		f, err := expr.Compile(express , expr.Env(params), expr.AsBool(),expr.AllowUndefinedVariables() )
		if err!= nil {
			return false, err
		}
		gBoolExprMap[express] = f
		cmdEx  =f
	}
	return expr.Run( cmdEx, params )
}


func doFilter( filter *Filter,   params map[string]interface{} ) (bool, error ) {

	//fmt.Printf(" do doFilter id:%v, name:%v \n", filter.ID, filter.Name )
	result, err := doBoolExpress( filter.When , params )
	if err != nil {
		return false , err
	}
	rst, ok := result.(bool)
	if !ok  {
		return false , fmt.Errorf("type faile ")
	}
	return rst ,  nil
}

func doAction( act  *Action, params map[string]interface{}  )(bool, error ){
	//fmt.Printf(" do action id:%v, name:%v \n", act.ID, act.Name )
	//执行 When
	{
		result, err := doBoolExpress( act.When , params )
		if err != nil {
			return false , err
		}
		rst, ok := result.(bool)
		if !ok  {
			return false , fmt.Errorf("type faile ")
		}
		if !rst {
			return false ,  nil
		}
	}
	//执行 Then
	{
		_, err := doExpress( act.Then, params )
		if err != nil {
			fmt.Printf("Result faile err:%v \n", err.Error() )
		}
	}
	//执行失败返回
	for _, v2 :=range act.Action{
		doAction( v2, params  )
	}
	return false, nil
}
