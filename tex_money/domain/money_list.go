package domain

type MoneyList struct {
	M10000 int // m:メイン
	M5000  int
	M2000  int
	M1000  int
	M500   int
	M100   int
	M50    int
	M10    int
	M5     int
	M1     int
	S500   int // s:サブ
	S100   int
	S50    int
	S10    int
	S5     int
	S1     int
	A10000 int // a:あふれ
	A5000  int
	A2000  int
	A1000  int
	A500   int
	A100   int
	A50    int
	A10    int
	A5     int
	A1     int
}

func NewMoneyList(before [EXTRA_CASH_TYPE_SHITEI]int) *MoneyList {
	return &MoneyList{
		M10000: before[0],
		M5000:  before[1],
		M2000:  before[2],
		M1000:  before[3],
		M500:   before[4],
		M100:   before[5],
		M50:    before[6],
		M10:    before[7],
		M5:     before[8],
		M1:     before[9],
		S500:   before[10],
		S100:   before[11],
		S50:    before[12],
		S10:    before[13],
		S5:     before[14],
		S1:     before[15],
		A10000: before[16],
		A5000:  before[17],
		A2000:  before[18],
		A1000:  before[19],
		A500:   before[20],
		A100:   before[21],
		A50:    before[22],
		A10:    before[23],
		A5:     before[24],
		A1:     before[25],
	}
}
