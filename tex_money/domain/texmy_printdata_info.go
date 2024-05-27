package domain

// 精算機日計表
type SummarySales struct {
	Amount      []int //売上金額
	Count       []int //売上回数
	TotalAmount int   //売上金額（合計）
	TotalCount  int   //売上回数（合計）
}
