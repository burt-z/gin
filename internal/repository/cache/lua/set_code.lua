-- 验证码在 Redis 上的 key
-- phone_code:login:152xxx
local key = KEYS[1]
-- 验证次数,记录了验证了几次
-- phone_code:login:152xxx:cnt
local countKey = key..":cnt"
-- 验证码
local val = ARGV[1]
-- 获取当前key的剩余过期时间
local ttl = tonumber(redis.call("ttl", key))

if ttl == -1 then
    -- key 存在,没有过期时间 - 系统错误
    return -2
elseif ttl == -2 or ttl < 540 then  -- 600-60=540，表示剩余时间小于9分钟时可以重新获取
    -- 设置验证码，过期时间10分钟
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    -- 设置验证次数为3，同样的过期时间
    redis.call("set", countKey, 3)
    redis.call("expire", countKey, 600)
    return 0  -- 成功
else
    -- 验证码获取太频繁
    return -1
end