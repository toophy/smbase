--
-- 本脚本不得修改, 不得删除
--
-- ts : 本线程对象(唯一)
--


-- Set up a custom loader
-- package.loaders[2] = function(name) return System:LoadModule(name) end

-- 禁止动态库和all-in-one库加载
package.loaders[3] = nil
package.loaders[4] = nil

-- 禁止加载其他库
package.loadlib = function() end

-- 禁止OS部分函数
local osOriginal = os

os = 
{
	date = os.date,
	time = os.time, 
	setlocale = os.setlocale,
	clock = os.clock, 
	difftime = os.difftime,
}

for k, v in pairs(osOriginal) do
	if not os[k] and type(v) == "function" then
		os[k] = function() end
	end
end

-- 线程ID
Tid_world       = 0    -- 世界线程
Tid_screen_1    = 1    -- 场景线程1
Tid_screen_2    = 2    -- 场景线程2
Tid_screen_3    = 3    -- 场景线程3
Tid_screen_4    = 4    -- 场景线程4
Tid_screen_5    = 5    -- 场景线程5
Tid_screen_6    = 6    -- 场景线程6
Tid_screen_7    = 7    -- 场景线程7
Tid_screen_8    = 8    -- 场景线程8
Tid_screen_9    = 9    -- 场景线程9
Tid_net_1       = 10   -- 网络线程1
Tid_net_2       = 11   -- 网络线程2
Tid_net_3       = 12   -- 网络线程3
Tid_db_1        = 13   -- 数据库线程1
Tid_db_2        = 14   -- 数据库线程2
Tid_db_3        = 15   -- 数据库线程3
Tid_last        = 16   -- 最终线程ID

Evt_lay1_time   = 160000 -- 第一层事件池最大支持时间(毫秒)

-- 当前目录
local curr_dir = ""
local path_obj = io.popen("cd")  --如果不在交互模式下，前面可以添加local 
curr_dir = path_obj:read("*all"):sub(1,-3)    --path存放当前路径
path_obj:close()   --关掉句柄

-- 自定义print函数, 指向线程普通信息日志函数
function print(...)
   local result = ""
   for i, v in ipairs{...} do
       result = result .. v .. ' '
   end
   ts:LogInfo(result)
end

-- 日志 : 调试信息
function LogDebug(...)
   local result = ""
   for i, v in ipairs{...} do
       result = result .. v .. ' '
   end
   ts:LogDebug(result)
end

-- 日志 : 普通信息
function LogInfo(...)
   local result = ""
   for i, v in ipairs{...} do
       result = result .. v .. ' '
   end
   ts:LogInfo(result)
end

-- 日志 : 警告
function LogWarn(...)
   local result = ""
   for i, v in ipairs{...} do
       result = result .. v .. ' '
   end
   ts:LogWarn(result)
end

-- 日志 : 普通错误
function LogError(...)
   local result = ""
   for i, v in ipairs{...} do
       result = result .. v .. ' '
   end
   ts:LogError(result)
end

-- 日志 : 严重错误
function LogFatal(...)
   local result = ""
   for i, v in ipairs{...} do
       result = result .. v .. ' '
   end
   ts:LogFatal(result)
end


-- 根据模块名获得脚本模块
function Mod(ModuleName)
  return package.loaded[ModuleName]
end

-- 打印table
function PrintTable(root)
  local cache = {  [root] = "." }
  local function _dump(t,space,name)
    local temp = {}
    for k,v in pairs(t) do
      local key = tostring(k)
      if cache[v] then
        table.insert(temp,"* " .. key .. " {" .. cache[v].."}")
      elseif type(v) == "table" then
        local new_key = name .. "." .. key
        cache[v] = new_key
        table.insert(temp,"+ " .. key .. _dump(v,space .. (next(t,k) and "|" or " " ).. string.rep(" ",#key),new_key))
      else
        table.insert(temp,"- " .. key .. " [" .. tostring(v).."]")
      end
    end
    return table.concat(temp,"\n"..space)
  end
  print(_dump(root, "",""))
end

-- 获取当前目录
function GetCurrDir()
  return curr_dir
end


-- 初始化随机函数
math.randomseed(os.time())
math.random(1, 100)
