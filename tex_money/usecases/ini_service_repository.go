package usecases

import "tex_money/domain"

type IniServiceRepository interface {
	UpdateIni(texCon *domain.TexContext, section string, key string, data string)
	MultipleUpdateIni(texCon *domain.TexContext, sectionKeyValue map[string]map[string]string)
}
