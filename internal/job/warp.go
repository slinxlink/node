package job

import (
	"time"

	"github.com/slinxlink/node/internal/core"
	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/service"
	"github.com/slinxlink/node/internal/util"
)

func WarpAutoUpdate() {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			time.Sleep(next.Sub(now))
			checkAndUpdate()
		}
	}()
}

func checkAndUpdate() {
	var warp database.Warp
	if err := database.DB.First(&warp).Error; err != nil || warp.AutoUpdate <= 0 {
		return
	}

	days := int(time.Since(warp.UpdatedAt).Hours() / 24)
	if days < warp.AutoUpdate {
		return
	}

	var endpoint database.Endpoint
	if err := database.DB.Where("tag = ?", "warp").First(&endpoint).Error; err != nil {
		return
	}

	util.Info("[warp] 自动更新IP (周期 %d 天)", warp.AutoUpdate)
	data, err := service.WarpRegister()
	if err != nil {
		util.Info("[warp] 自动更新IP失败: %v", err)
		return
	}

	endpoint.Address = data.Address
	endpoint.PrivateKey = data.PrivateKey
	endpoint.PublicKey = data.PublicKey
	endpoint.PeerAddress = data.PeerAddress
	endpoint.PeerPort = data.PeerPort
	endpoint.PeerPublicKey = data.PeerPublicKey
	endpoint.Reserved = data.Reserved

	if err := database.DB.Save(&endpoint).Error; err != nil {
		util.Info("[warp] 更新端点失败: %v", err)
		return
	}

	go core.Default.Apply()
}
