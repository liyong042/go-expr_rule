# 设计参考 Nginx Lua扩展模块 
    https://frxcat.fun/middleware/Nginx/Nginx_Lua_Expansion_module/#ngx-exec

## 行业解决方案
    开源库：govaluate， expr , gengine（B站开源），goja( js脚本 ) , Lua脚本
    金融反欺诈、金融信审等互金领域，由于黑产、羊毛党行业的盛行，风控决策引擎在电商、支付、游戏、
    社交等领域也有了长足的发展，刷单、套现、作弊，凡是和钱相关的业务都离不开风控决策引擎系统的支持保障
    行业内可选方案：一套商业决策引擎系统动辄百万而且需要不断加钱定制，大多数企业最终仍会走上自研道路，市场上有些开源规则引擎项目可参考，
    比较出名的开源规则引擎有drools、urule，都是基于Rete算法，都是基于java代码实现，一些企业会进行二次开发落地生产。
    而这类引擎功能强大但也比较“笨重”，易用性以及定制性并不够好，对其他语言栈二次开发困难

## 快速Demo
```
{
"Id":1,
"Filter": [
    {
      "Id": "1",
      "When": "Req.Action() in ['act_user_join','act_meet_create']",
      "Name": " 需要的用户动作"
    }
  ],
  "Action": [
    {
      "Id": "2",
      "When": "Req.Get( 'City' ) in ['云南','西藏',\" ttt \"]",
      "Then": "Rsp.Write('doSetRedisOk','{\"name\":33333}')",
      "Name": " 需要的 用户动作",
      "Action": [
        {
          "Id": "21",
          "When": "Req.Get( 'City' )  in ['云南','西藏',\" ttt \"]",
      	  "Then": " Req.Get('yy') + Rsp.Write( '' , '{\"name\":22222}')",
          "Name": " 需要的 21"
        }
      ]
    },
    {
      "Id": "3",
      "When": "Req.Get( 'City' )  in ['云南','西藏'] ",
      "Then": "Req.Get2('yy') + Rsp.Write('doSetRedis','{\"name\":999999}')",
      "Name": " 需要的 3 动作"
    }
  ]
}

{
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
      "Name": " 需要的 用户动作",
      "ForEach": {
        "Items": "VarAaray",
        "When": "true",
        "ThenList": [
          "Rsp.Write( k , v  ) "
        ]
      }
    }
  ]
}

```

## 服务变量定义 
    Sys 系统插件入口 
    Req 业务自定义请求
    Rsp 业务自定义响应
    Ctx 业务请求上下文 

## 上下文插件(Ctx) 
    Ctx.RuleId                策略编号
    Ctx.SetGlobal( 'GlbXxx' ) 设置全局变量 ,必须Glb开头 
    Ctx.SetLocal ( 'VarXxx' ) 设置局部变量 ,必须Var开头  
    Ctx.Exit()                退出程序执行 
    Ctx.SubCancel()           子策略取消 
    Ctx.Return()              当前策略返回
    Ctx.DoSubRule( 'xxx')     执行子策略
    Ctx.Log('')               打印日志   
    Ctx.Get(key interface{}) interface{}  对接 Ctx.Value(key interface{}) interface{}

## 系统插件
    Sys.Load  (  'xxx' )      加载自定义插件 
    Sys.Import( 'VarPlgXxx', 'xxx' ) 引入插件  
    Sys.Sleep( 10 )          sleep函数 10毫秒   

## 业务请求  Req, Rsp，业务自定义，可以是 http, pb , json,image 等  
    Req.Action()                           需要实现接口 
    Req.Get( path ) Interface{}            需要实现接口
    Req.Set( path, val interface{} ) bool  需要实现接口 

    Rsp.Write(args ...interface{}) int     需要实现接口 

## 业务自定义私有插件 
    Sys.Load( 'PriXxx' )        插件名称以 Pri开头  
    
