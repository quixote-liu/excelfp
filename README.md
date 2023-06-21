# execlfp
解析并执行excel函数表达式，返回相应的结果

excel函数信息：https://www.lanrenexcel.com/excel-functions-list/

excel的函数表达式示例：
```
IFS(AND({a}=1,{b}=0,{c}=1), "停止", AND({a}=1,{b}>0,{c}>0), "运行", AND({a}=0), "离线")
```
