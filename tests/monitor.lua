local default_collectors = luakit.monitor.collectors

local alarm=luakit.monitor.alarm.simple("简单告警器")

local metrics = luakit.monitor.metrics
local cpu = default_collectors.cpu {
    interval = 10,
}
alarm.addSimple(cpu,"cpu_usage > 1", "cpu使用率超过1%")
--alarm.addAvg(cpu,"cpu.usage > 80", "cpu使用率超过80%", 5) -- 5次采集平均值

local mem = default_collectors.mem {
    interval = 10,
}
local disk = default_collectors.disk {
    interval = 1800,
}
local net = default_collectors.net {
    interval = 10,
}
local go = default_collectors.go {
    interval = 10,
}

local self_process = default_collectors.self_process {
    interval = 10,
}

local req_cnt = metrics.counter('req_success_cnt', '总请求计数器')
local success_cnt = metrics.counter('req_success_cnt', '请求成功计数器')

-- 生成一个速率指标
local success_cnt_rate_per5s = success_cnt.gen_rate_metric('req_success_rate', '每秒请求成功数(每隔5秒采集平均值)', 5)

local m = luakit.monitor {
    name = 'x_monitor',
}

m.collectors(cpu, mem, disk, net,self_process)
--m.collectors(cpu, mem, disk, net)

m.metrics(success_cnt,req_cnt, success_cnt_rate_per5s)
m.PrometheusPull {
    prom_pull_addr = '0.0.0.0:9100',
    prom_pull_uri = '/metrics',
}
 m.PrometheusPush{
     prom_push_gateway = "http://127.0.0.1:8428/prometheus/api/v1/import/prometheus",
     prom_push_interval = 5,
     prom_push_job_name = "rock_limit_metrics",
     prom_push_username = "",
     prom_push_password = "",
 }
 m.SimplePull("0.0.0.0:9101")
m.start()
success_cnt.incr()
success_cnt.incr()
success_cnt.incr()
success_cnt.incr()
success_cnt.incr()
req_cnt.incr()
for i = 1, 1000000 do
    success_cnt.incr()
end
