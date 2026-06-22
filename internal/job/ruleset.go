package job

import (
	"os"
	"time"

	"github.com/slinxlink/node/internal/config"
	"github.com/slinxlink/node/internal/sub/template"
	"github.com/slinxlink/node/internal/util"
)

func RulesetRefresh() {
	go func() {
		checkRulesetFiles()

		for {
			now := time.Now()
			daysUntilMonday := (time.Monday - now.Weekday() + 7) % 7
			if daysUntilMonday == 0 {
				daysUntilMonday = 7
			}
			next := time.Date(now.Year(), now.Month(), now.Day()+int(daysUntilMonday), 0, 0, 0, 0, now.Location())
			time.Sleep(next.Sub(now))

			if !config.Config.RulesetAutoUpdate {
				continue
			}
			RefreshRuleset()
		}
	}()
}

func checkRulesetFiles() {
	_, errJson := os.Stat(template.RulesetPath)
	_, errClash := os.Stat(template.ClashTemplatePath)
	_, errSurge := os.Stat(template.SurgeTemplatePath)
	_, errSingBox := os.Stat(template.SingBoxTemplatePath)

	if os.IsNotExist(errJson) {
		RefreshRuleset()
	} else {
		if os.IsNotExist(errClash) {
			template.GenerateClash()
		}
		if os.IsNotExist(errSurge) {
			template.GenerateSurge()
		}
	}

	if os.IsNotExist(errSingBox) {
		template.GenerateSingBox()
	}
}

func RefreshRuleset() {
	count, err := template.DownloadRuleset()
	if err != nil {
		util.Error("[ruleset] 规则下载失败: %v", err)
		return
	}
	util.Info("[ruleset] 规则下载完成，共 %d 条", count)
	template.GenerateClash()
	template.GenerateSurge()
	template.GenerateSingBox()
}
