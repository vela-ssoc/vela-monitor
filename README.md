# vela-monitor
提供统一的性能\自定义指标\系统资源等的监控组件，用于全方位监控各种资源的状态

## 基本使用
```lua
local cnt =luakit.monitor.metrics.counter("req_cnt", "请求计数器")
local m = luakit.monitor{
    name = "服务名称"
}
m.collectors(luakit.monitor.collectors.cpu{interval = 10,})  -- 添加采集器
m.metrics(cnt)  -- 添加指标
m.PrometheusPull{prom_pull_addr = "0.0.0.0:9100",prom_pull_uri = "/metrics"}  -- 添加适配器
m.start()  -- 启动服务
cnt.incr()  -- 计数器+1
```
## 测试
在`tests/monitor.lua`中可以找到测试用例，运行`tests/service_test.go`中的`TestService`即可启动测试服务。

## 采集器
支持以下内置采集器：
- CPU采集器
- 内存采集器
- 磁盘采集器
- 网络采集器
- 普罗米修斯原生GO指标采集器
- 普罗米修斯原生自身进程指标采集器

### 1. CPU采集器
- **功能**：采集CPU使用率和CPU时间
- **Lua配置示例**：
```lua
local cpu = luakit.monitor.collectors.cpu{
    interval = 10,  -- 采集间隔(秒)
}
```

### 2. 内存采集器
- **功能**：采集内存使用情况
- **Lua配置示例**：
```lua
local mem = luakit.monitor.collectors.mem{
    interval = 10,
}
```

### 3. 磁盘采集器
- **功能**：采集磁盘使用情况和指定目录的磁盘空间
- **Lua配置示例**：
```lua
local disk = luakit.monitor.collectors.disk{
    interval = 1800,  -- 30分钟采集一次
    targets = {"C:\\", "D:\\"}  -- 监控的磁盘目录
}
```

### 4. 网络采集器
- **功能**：采集网络流量和连接数
- **Lua配置示例**：
```lua
local net = luakit.monitor.collectors.net{
    interval = 10,
}
```

### 5. 自身进程采集器
- **功能**：采集自身进程的资源使用情况
- **Lua配置示例**：
```lua
local self_process = luakit.monitor.collectors.self_process{
    interval = 10,
}
```

### 6. 原生GO指标采集器
- **功能**：采集自身go程序的资源使用情况
- **Lua配置示例**：
```lua
local self_process = luakit.monitor.collectors.go{
    interval = 10,
}
```

### 7. 通用采集器
- **功能**：自定义采集器和指标集合
- **Lua配置示例**：
```lua
local c = luakit.monitor.collectors.new{
    name = "自定义采集器",
    help = "自定义采集器描述",
    metrics = {
        metrics.counter('req_cnt', '请求计数器'),
        metrics.counter('req_fail_cnt', '请求失败计数器')
    }
}
```
## 指标

### 1. 原子计数器(atomic counter)
- **用途**：记录事件发生的总次数
- **特性**：
  - 线程安全，支持并发递增
  - 只能增加，不能减少
- **示例**：
```lua
local cnt = luakit.monitor.metrics.counter("req_cnt", "请求计数器")
cnt.incr()  -- 计数器+1
cnt.add(5)   -- 计数器+5
```

### 2. 简单指标(simple gauge)
- **用途**：记录瞬时值
- **特性**：
  - 可设置任意数值
  - 适合记录内存使用率、CPU负载等指标
- **示例**：
```lua
local mem = luakit.monitor.metrics.simple_gauge("mem_usage", "内存使用百分比")
mem.set(75.3)  -- 设置当前内存使用率为75.3%
```

### 3. 速率计算器(rate calculator)
- **用途**：计算指标的变化速率
- **特性**：
  - 基于滑动窗口计算平均值
  - 适合计算QPS、网络吞吐等指标
- **示例**：
```lua
local req_cnt = luakit.monitor.metrics.counter("req_total", "总请求数")
local qps = req_cnt.gen_rate_metric("req_rate", "每秒请求数", 5)  -- 5秒窗口

-- 添加到监控服务
m.metrics(req_cnt, qps)
```

## 适配器


### Prometheus集成
支持Pull和Push两种模式：

1. Pull模式：
```lua
m.PrometheusPull{
    prom_pull_addr = "0.0.0.0:9100",
    prom_pull_uri = "/metrics"
}
```

2. Push模式：
```lua
m.PrometheusPush{
    prom_push_gateway = "http://127.0.0.1:9091",
    prom_push_interval = 5,
    prom_push_job_name = "job_name"
}
```
### 内部数据采集PULL接口
还在开发中, 目前只有最简易的实现, 以标准json格式输出所有采集的数据信息或者触发采集动作

``` lua
 m.SimplePull("0.0.0.0:9101")
```
#### 接口URI
`/onekit/monitor/collect` 立即执行采集一次, 返回所有采集器的采集数据，格式为JSON。  
`/onekit/monitor/view` 返回所有数据，格式为JSON。    
**注意**: 如果采集器/指标没有配置自动定时采集, 那么里面的一些指标不会更新  

### 内部数据采集PUSH服务
TODO


## 告警器
### 简单告警器
**还在开发中, 部分功能还未完成**  
>只支持直接指标的告警, 暂时不支持复杂运算
- 方法 `outputLog()` 会输出告警日志, 可以通过配置日志输出到文件或其他方式
- 方法 `outputSiem()` 会发送告警到SIEM, 可以通过配置SIEM地址和参数来发送告警
- 方法 `addSuppression(int,int)` 可以添加抑制规则, 第一个参数为抑制时间(秒), 第二个参数为抑制次数, 当在抑制时间内达到抑制次数时, 告警器会抑制告警

使用样例:
```lua
local alarm=luakit.monitor.alarm.simple("简单告警器")
alarm.addSimple(cpu,"cpu_usage > 1", "cpu使用率超过1%").outputLog()
```


## 启动服务
```lua
m.start()
```

## Lua调用完整示例
参考`tests/monitor.lua`中的示例代码

## Go调用Lua定义的采集器
>有时候, 我们需要在GO代码中进行更底层的一些计数和性能统计操作.  
此时可以通过接收调用Lua定义的采集器来实现, 这样使得Go中采集器更加灵活和可扩展.  

以下示例展示了如何在Go代码中接收并使用Lua中定义的采集器对象.
更为详细的示例请参考`tests/demotask.go`中的示例代码
```go
// 预制一个计数器指标 并在相应采集点插桩
var req_cnt *metrics.AtomicCounter 
// 或者预制一个计数器指标list 
// 再或者 实现GeneralCollectorI接口对象 并在相应的采集点插桩


func demoL(L *lua.LState) {
    // 假设Lua脚本已加载并定义了采集器
    c := lua.Check[lua.GenericType](L, L.Get(1))



    if v, ok := c.Unpack().(collector.GeneralCollectorI); ok {
        fmt.Println("采集器名称:", v.Name())
        ms = v.Metrics()
    }

    for _, m := range ms {
        fmt.Println("指标名称:", (*m).Name(), "值:", (*m).Value())
        if v, ok := (*m).(*metrics.AtomicCounter); ok {
            // 把lua中定义的采集器映射到你预置的Go中的指标对象
            // 具体怎样的映射规则, 取决于你自己
            if (*m).Name() == "req_fail_cnt" {
                req_cnt = v
            }
        }
    }

    // 模拟计数器操作
    if req_cnt != nil {
        req_cnt.Add(10)
        fmt.Println("更新后的计数器值:", req_cnt.Value())
    }
}
```

### 说明
- 调用者可以通过接口断言(`collector.GeneralCollectorI`)来获取采集器实例。
- 遍历采集器中的指标列表，根据指标名称进行映射和操作。
- 在操作指标前，建议进行空值检查以避免运行时错误。
