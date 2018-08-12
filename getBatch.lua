-- ARGV[1] = queue name
-- ARGV[2] = batch name
-- ARGV[3] = max batch

for i=0,ARGV[3] do
    redis.call("RPOPLPUSH",ARGV[1],ARGV[2])
end
local element = redis.call("LRANGE",ARGV[2],0,-1)

return element