package usecases

import "stub/domain"

type EqualsTexMoneyRepository interface {
	IsEqualStatusCash(a, b domain.StatusCash) bool
}
