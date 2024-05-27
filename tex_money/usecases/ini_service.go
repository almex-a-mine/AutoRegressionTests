package usecases

import (
	"fmt"
	"sync"
	"tex_money/domain"
	"tex_money/domain/handler"
	"tex_money/pkg/file"

	"gopkg.in/ini.v1"
)

type iniService struct {
	logger   handler.LoggerRepository
	FilePath string
	FileName string
	mu       sync.Mutex
}

func NewIniService(logger handler.LoggerRepository) IniServiceRepository {
	// env name
	FilePath := file.GetCurrentDir()
	env := file.GetEnv("ALMEXPATH")
	if len(env) != 0 {
		env = file.AdjustFileName(env)
		if file.DirExists(env + "ini") {
			FilePath = env + "ini"
		}
	}
	dirPath := file.AdjustFileName(FilePath)
	FileName := fmt.Sprintf("%v%s.ini", dirPath, domain.AppName)
	return &iniService{
		logger:   logger,
		FilePath: FilePath,
		FileName: FileName,
	}
}

// データをiniファイルに書き込む
func (c *iniService) UpdateIni(texCon *domain.TexContext, section string, key string, data string) {
	c.logger.Trace("【%v】START:NewIniService UpdateIni section=%s, key=%s, data=%s", texCon.GetUniqueKey(), section, key, data)

	//書き込み処理で競合が発生しないようにロック処理を追加
	c.mu.Lock()
	defer c.mu.Unlock()

	if fileExists := file.FileExists(c.FileName); fileExists {
		cfg, err := ini.Load(c.FileName)
		if err != nil {
			c.logger.Debug("【%v】- iniファイル取得失敗 err=%v", texCon.GetUniqueKey(), err)
			return
		}
		// iniファイル値チェック
		m := cfg.Section("SYSTEM").Key("MaxLength").String()
		if m == "" {
			c.logger.Debug("【%v】- iniファイルデータ存在チェックデータが空の為、停止しました。", texCon.GetUniqueKey())

		}

		cfg.Section(section).Key(key).SetValue(data)
		if err := cfg.SaveTo(c.FileName); err != nil {
			c.logger.Debug("【%v】- iniファイル書き込み失敗 err=%v", texCon.GetUniqueKey(), err)
			return
		}
	}

	c.logger.Trace("【%v】END:NewIniService UpdateIni", texCon.GetUniqueKey())
}

// データをiniファイルに書き込む
func (c *iniService) MultipleUpdateIni(texCon *domain.TexContext, sectionKeyValue map[string]map[string]string) {
	c.logger.Trace("【%v】START:NewIniService MultipleUpdateIni sectionKeyValue=%v", texCon.GetUniqueKey(), sectionKeyValue)

	//書き込み処理で競合が発生しないようにロック処理を追加
	c.mu.Lock()
	defer c.mu.Unlock()

	if fileExists := file.FileExists(c.FileName); fileExists {
		cfg, err := ini.Load(c.FileName)
		if err != nil {
			c.logger.Debug("【%v】- iniファイル取得失敗 err=%v", texCon.GetUniqueKey(), err)
			return
		}
		// iniファイル値チェック
		m := cfg.Section("SYSTEM").Key("MaxLength").String()
		if m == "" {
			c.logger.Debug("【%v】- iniファイルデータ存在チェックデータが空の為、停止しました。", texCon.GetUniqueKey())

		}

		for sec, keyVal := range sectionKeyValue {
			for key, val := range keyVal {
				cfg.Section(sec).Key(key).SetValue(val)
			}
		}

		if err := cfg.SaveTo(c.FileName); err != nil {
			c.logger.Debug("【%v】- iniファイル書き込み失敗 err=%v", texCon.GetUniqueKey(), err)
			return
		}

	}

	c.logger.Trace("【%v】END:NewIniService MultipleUpdateIni", texCon.GetUniqueKey())
}
