local colt = luakit.monitor.collectors
local metrics = luakit.monitor.metrics
local cpu = colt.cpu {interval = 10,}
local mem = colt.mem {interval = 10,}
local disk = colt.disk {
    interval = 1800,
    targets = {
        "D:\\",
        "E:\\"
    }
}
local net = colt.net { interval = 10,}
local go = colt.go {interval = 10,}
local self_process = colt.self_process {interval = 10,}

local req_cnt = metrics.counter('req_cnt', '总请求计数器')
local success_cnt = metrics.counter('req_success_cnt', '请求成功计数器')

-- 生成一个速率指标
local success_cnt_rate_per5s = success_cnt.gen_rate_metric('req_success_rate', '每秒请求成功数(每隔5秒采集平均值)', 5)

-- 生成一个采集器
local c = colt.new {
    name = "xxx模块采集器",
    help = "xxx模块采集器 详细描述",
    metrics = {
        success_cnt,
        req_cnt,
        success_cnt_rate_per5s,
        metrics.counter('req_fail_cnt', '请求失败计数器'),
    }
}


local alarm = luakit.monitor.alarm.simple("简单告警器")
alarm.addSimple(cpu, "cpu_usage > 10", "cpu使用率超过10%").outputLog()
--alarm.addAvg(cpu,"cpu.usage > 80", "cpu使用率超过80%", 5) -- 5次采集平均值


local m = luakit.monitor {
    name = 'x_monitor',
}

-- 采集器 (相当于指标集合和统一的采集方法)
m.collectors(cpu, mem, disk, net, self_process,c)
--m.collectors(cpu, mem, disk, net)

-- 单独定义的指标(孤儿指标)
-- 注意 不要重复定义指标 name必须唯一
-- m.metrics(success_cnt, req_cnt, success_cnt_rate_per5s)
m.PrometheusPull {
    prom_pull_addr = '0.0.0.0:9100',
    prom_pull_uri = '/metrics',
}
m.PrometheusPush {
    prom_push_gateway = "http://127.0.0.1:8428/prometheus/api/v1/import/prometheus",
    prom_push_interval = 5,
    prom_push_job_name = "rock_limit_metrics",
    prom_push_username = "",
    prom_push_password = "",
}
-- /onekit/monitor/view
-- /onekit/monitor/collect
m.SimplePull("0.0.0.0:9101")

luakit.demotask(c)
m.start()

-- 测试计数器(lua 层)
-- success_cnt.incr()
-- success_cnt.incr()
-- success_cnt.incr()
-- success_cnt.incr()
-- success_cnt.incr()
-- req_cnt.incr()
-- for i = 1, 1000000 do
--     success_cnt.incr()
-- end

-- 测试计数器(go 层)
