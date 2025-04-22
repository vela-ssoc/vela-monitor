local default_collectors = luakit.monitor.collectors

-- 自定义采集器

local c1 =luakit.monitor.counter{
    name = 'test_counter',
    desc = 'test counter1111111',
    value = 1,
}
-- 原子操作 只支持incr和add
c1.incr()
c1.Add(100)

local c2 =luakit.monitor.gauge{
    name = 'test_simple_gauge',
    desc = 'test gauge1111111',
    value = 123456,
}
-- 简单set值 如读取内存使用量等
c2.set(11111)

local c1_rate_5s=luakit.monitor.rateCalculator{
    metric = c1,
    interval = 5,
    smooth = true,
}
c1_rate_5s.rate()

-- 滚动窗口均值
local c1_avg_5s = luakit.monitor.avgCalculator{
    metric = c1,
    store = 3600,
}
c1_avg_5s.avg(10)

-- 直方图
local c1_histogram = luakit.monitor.histogram{
    name = 'test_histogram',
    metric = c1,
}
print(c1_histogram.data()["top1%"])
print(c1_histogram.data()["low1%"])

-- 对rate的指标进行二次运算
local c1_rate_histogram = luakit.monitor.histogram{
    name = 'test_histogram',
    metric = c1_rate_5s,
    interval = 5,
}
print(c1_rate_histogram.data()["top90%"])


local m = luakit.monitor{
    service_name = 'x_monitor',
    collectors = {
        -- 使用内部采集器
        default_collectors.cpu{
            interval = 10,
        },
        default_collectors.mem{
            interval = 10,
        },
        default_collectors.disk{
            interval = 3600,
        },
        default_collectors.net{
            interval = 10,
        },
        c1,
        c2,
        c1_rate_5s,
    }
}

m.PrometheusPull{
    prom_pull_addr = '0.0.0.0:9100',
    path = '/metrics',
}.run()
m.PrometheusPush {
    address = '192.168.1.2:9091',
    job = 'x_monitor',
    interval = 10,
    username = 'test',
    password = 'test',
}.run()
m.SimplePush{
    job = 'x_monitor',
    interval = 10,
}.run()
m.SimplePull{
    inner_uri = '/x_monitor',
}.run()
m.start()