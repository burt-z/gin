local key = KEYS[1]
local countKey = key.."cnt"
--输入的验证码
local expectedCode = ARGV[1]
-- 过期时间
local count = tonumber(redis.call("get",countKey))
local code = redis.call("get",key)

if count == nil or count <= 0 then
    return -1
end

if count <= 0 then
    -- key  一直输错
    -- 自己约定的错误码,表示系统错误
    return -1
    --   -2 key 不存在 600-60
elseif expectedCode == code then
    redis.call("set",countKey,-1)
    return 0
else
    -- 输错,但合法,可以验证次数减 1
    redis.call("decr",key)
    return -2
end