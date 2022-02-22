package main

import (
	"encoding/json"
	"fmt"
	"test/rule"
	"time"
)
var s1 =`{
"Id":1,
"Filter": [
    {
      "Id": "1",
      "When": "Req.Action in ['act_user_join','act_meet_create']",
      "Name": " 需要的用户动作"
    }
  ],
  "Action": [
    {
      "Id": "2",
      "When": "Req.City in ['云南','西藏',\" ttt \"]",
      "Then": "Rsp.Write('doSetRedisOk','{\"name\":33333}')",
      "Name": " 需要的 用户动作",
      "Action": [
        {
          "Id": "21",
          "When": "Req.City in ['云南','西藏',\" ttt \"]",
      	  "Then": " Req.Get('yy') + Rsp.Write( '' , '{\"name\":22222}')",
          "Name": " 需要的 21"
        }
      ]
    },
    {
      "Id": "3",
      "When": "Req.City in ['云南','西藏'] ",
      "Then": "Req.Get2('yy') + Rsp.Write('doSetRedis','{\"name\":999999}')",
      "Name": " 需要的 3 动作"
    }
  ]
}
`
var s2 =`{
"Id":2,
"Filter": [
    {
      "Id": "1",
      "When": "true",
      "Name": " 需要的用户动作"
    }
  ],
  "Action": [
    {
      "Id": "2",
      "When": "Req.City in ['云南','西藏','ttt'']",
      "Then": "Rsp.Write('doSetRedis','{\"name\":555}')",
      "Name": " 需要的 用户动作"
    },
    {
      "Id": "3",
      "When": "true",
      "Then": "Req.Get2( name==nil?'--- test111 name':'test2') ",
      "Name": " 需要的 用户动作"
    },
    {
      "Id": "4",
      "When": "true",
      "Then": "Req.Get2( Req.Get('yy11')  )  ",
      "Name": " 需要的 用户动作"
    }
  ]
}
`

// Request 请求
type Request struct {
	Action        string `json:"Action"`
	City          string `json:"City"`
}

func (c *Request ) Get( act string )string {
	fmt.Printf("[Get] %v \n",act  )
	return act + "--Get---"
}
type Request2 struct {
	*Request
	Action2        string `json:"Action"`
}
func (c *Request2 ) Get2( act string )string {
	fmt.Printf("[Get2] %v \n",act  )
	return act + "-Get2---"
}

// Response 响应
type Response struct {
	Id  int
}
func (c *Response ) Write( act, params  string )string  {
	fmt.Printf("[Write] Write id:%v  act: %v, params:%v  \n ", c.Id , act, params  )
	return act + "-hh-"
}

func main(){

	//参数
	params := make(map[string]interface{}, 8)
	req := &Request2{ Request:&Request{Action:"act_user_join", City:"云南"} }
	params["Req"] = req


	fmt.Println("-------- test1  --------- ")
	//测试1
	Test( s1, params )
	fmt.Println("-------- test2  --------- ")
	//测试2
	Test( s2, params )
	fmt.Println("-------- test3 --------- ")

	//测试3
	req.Action = "test"
	req.City   = "uyy"

	Test( s1, params )
	fmt.Println("-------- test4  --------- ")
	//测试4
	Test( s2, params )

	fmt.Println("--------test count  --------------------------------- ")
	{
		TestCnt( s1, params )
		TestCnt( s2, params )
	}
}

// Test 测试
func Test( s string, params map[string]interface{} ){
	val:= &rule.ConfigValue{}
	if err := json.Unmarshal([]byte(s) , val ); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf( "req: %+v \n ", val )
	rsp :=  &Response{ }
	rsp.Id = val.Id
	params["Rsp"] = rsp
	rule.DoRule( val, params )
}

// TestCnt 压测
func TestCnt( s string, params map[string]interface{} ){
	val:= &rule.ConfigValue{}
	if err := json.Unmarshal([]byte(s) , val ); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println( *val )
	rsp :=  &Response{ }
	rsp.Id = val.Id
	params["Rsp"] = rsp
	cnt := 1
	startT := time.Now() //计算当前时间
	for i := 1; i <= cnt; i++ {
		rule.DoRule( val, params )
	}
	tc := time.Since(startT) //计算耗时
	fmt.Printf("pb test  cnt:%v, time cost = %v\n", cnt, tc)

}