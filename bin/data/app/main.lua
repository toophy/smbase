module("main", package.seeall)

-- local 
function OnAppBegin()
	LogDebug("App 启动")

	ts:PostEvent("main","Eon_Qiguan",3000,{["log"]="咕咕鸟在鸣叫!"})
end

-- local 
function OnAppEnd()
	LogDebug("App 结束")
end

-- local 
function Eon_Qiguan(t)
	LogInfo(t["log"])
end
