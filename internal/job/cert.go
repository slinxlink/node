package job

import (
	"time"

	"github.com/slinxlink/node/internal/cert"
	"github.com/slinxlink/node/internal/database"
	"github.com/slinxlink/node/internal/task"
	"github.com/slinxlink/node/internal/util"
)

func CertRenew() {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			time.Sleep(next.Sub(now))
			checkAndRenew()
		}
	}()
}

func checkAndRenew() {
	var list []database.Cert
	database.DB.Where("auto_renew = ? AND mode != ?", true, "manual").Find(&list)

	for _, dbCert := range list {
		if dbCert.ExpireAt.IsZero() {
			continue
		}
		days := time.Until(dbCert.ExpireAt).Hours() / 24
		if days > 10 {
			continue
		}
		util.Info("[cert] 证书即将过期，自动续签: %s (剩余 %.0f 天)", dbCert.Domain, days)
		t := task.New("renew-" + dbCert.Domain)
		go cert.ApplyCert(&dbCert, t)
	}
}
