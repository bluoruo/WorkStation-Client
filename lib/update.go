package lib

import (
	"client/update"
)

// Update 更新
func Update() {
	//当前版本
	update.AppVer = AppVer
	//检查更新地址
	update.CheckUrl = serverConfig.Server + "/v2/config"
	//执行update
	update.RunClientUpdateCom()
}
