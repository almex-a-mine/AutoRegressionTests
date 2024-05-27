package calculation

import (
	"reflect"
	"testing"
)

func TestMoneyCalculation_Add(t *testing.T) {
	type fields struct {
		m10000 int
		m5000  int
		m2000  int
		m1000  int
		m500   int
		m100   int
		m50    int
		m10    int
		m5     int
		m1     int
		s500   int
		s100   int
		s50    int
		s10    int
		s5     int
		s1     int
		a10000 int
		a5000  int
		a2000  int
		a1000  int
		a500   int
		a100   int
		a50    int
		a10    int
		a5     int
		a1     int
	}
	type args struct {
		plus [26]int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   [26]int
	}{
		{
			name: "足し算",
			fields: fields{
				m10000: 1,
				m5000:  2,
				m2000:  3,
				m1000:  4,
				m500:   5,
				m100:   6,
				m50:    7,
				m10:    8,
				m5:     9,
				m1:     10,
				s500:   11,
				s100:   12,
				s50:    13,
				s10:    14,
				s5:     15,
				s1:     16,
				a10000: 17,
				a5000:  18,
				a2000:  19,
				a1000:  20,
				a500:   21,
				a100:   22,
				a50:    23,
				a10:    24,
				a5:     25,
				a1:     26,
			},
			args: args{
				plus: [26]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			},
			want: [26]int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MoneyCalculation{
				m10000: tt.fields.m10000,
				m5000:  tt.fields.m5000,
				m2000:  tt.fields.m2000,
				m1000:  tt.fields.m1000,
				m500:   tt.fields.m500,
				m100:   tt.fields.m100,
				m50:    tt.fields.m50,
				m10:    tt.fields.m10,
				m5:     tt.fields.m5,
				m1:     tt.fields.m1,
				s500:   tt.fields.s500,
				s100:   tt.fields.s100,
				s50:    tt.fields.s50,
				s10:    tt.fields.s10,
				s5:     tt.fields.s5,
				s1:     tt.fields.s1,
				a10000: tt.fields.a10000,
				a5000:  tt.fields.a5000,
				a2000:  tt.fields.a2000,
				a1000:  tt.fields.a1000,
				a500:   tt.fields.a500,
				a100:   tt.fields.a100,
				a50:    tt.fields.a50,
				a10:    tt.fields.a10,
				a5:     tt.fields.a5,
				a1:     tt.fields.a1,
			}
			if got := c.Add(tt.args.plus); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoneyCalculation_Exchange(t *testing.T) {
	type fields struct {
		m10000 int
		m5000  int
		m2000  int
		m1000  int
		m500   int
		m100   int
		m50    int
		m10    int
		m5     int
		m1     int
		s500   int
		s100   int
		s50    int
		s10    int
		s5     int
		s1     int
		a10000 int
		a5000  int
		a2000  int
		a1000  int
		a500   int
		a100   int
		a50    int
		a10    int
		a5     int
		a1     int
	}
	type args struct {
		amount     int
		changeType int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   [26]int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MoneyCalculation{
				m10000: tt.fields.m10000,
				m5000:  tt.fields.m5000,
				m2000:  tt.fields.m2000,
				m1000:  tt.fields.m1000,
				m500:   tt.fields.m500,
				m100:   tt.fields.m100,
				m50:    tt.fields.m50,
				m10:    tt.fields.m10,
				m5:     tt.fields.m5,
				m1:     tt.fields.m1,
				s500:   tt.fields.s500,
				s100:   tt.fields.s100,
				s50:    tt.fields.s50,
				s10:    tt.fields.s10,
				s5:     tt.fields.s5,
				s1:     tt.fields.s1,
				a10000: tt.fields.a10000,
				a5000:  tt.fields.a5000,
				a2000:  tt.fields.a2000,
				a1000:  tt.fields.a1000,
				a500:   tt.fields.a500,
				a100:   tt.fields.a100,
				a50:    tt.fields.a50,
				a10:    tt.fields.a10,
				a5:     tt.fields.a5,
				a1:     tt.fields.a1,
			}
			if got := c.Exchange(tt.args.amount, tt.args.changeType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Exchange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoneyCalculation_GetTotalAmount(t *testing.T) {
	type fields struct {
		m10000 int
		m5000  int
		m2000  int
		m1000  int
		m500   int
		m100   int
		m50    int
		m10    int
		m5     int
		m1     int
		s500   int
		s100   int
		s50    int
		s10    int
		s5     int
		s1     int
		a10000 int
		a5000  int
		a2000  int
		a1000  int
		a500   int
		a100   int
		a50    int
		a10    int
		a5     int
		a1     int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "合計金額",
			fields: fields{
				m10000: 1,
				m5000:  2,
				m2000:  3,
				m1000:  4,
				m500:   5,
				m100:   6,
				m50:    7,
				m10:    8,
				m5:     9,
				m1:     10,
				s500:   11,
				s100:   12,
				s50:    13,
				s10:    14,
				s5:     15,
				s1:     16,
				a10000: 17,
				a5000:  18,
				a2000:  19,
				a1000:  20,
				a500:   21,
				a100:   22,
				a50:    23,
				a10:    24,
				a5:     25,
				a1:     26,
			},
			want: 373407,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MoneyCalculation{
				m10000: tt.fields.m10000,
				m5000:  tt.fields.m5000,
				m2000:  tt.fields.m2000,
				m1000:  tt.fields.m1000,
				m500:   tt.fields.m500,
				m100:   tt.fields.m100,
				m50:    tt.fields.m50,
				m10:    tt.fields.m10,
				m5:     tt.fields.m5,
				m1:     tt.fields.m1,
				s500:   tt.fields.s500,
				s100:   tt.fields.s100,
				s50:    tt.fields.s50,
				s10:    tt.fields.s10,
				s5:     tt.fields.s5,
				s1:     tt.fields.s1,
				a10000: tt.fields.a10000,
				a5000:  tt.fields.a5000,
				a2000:  tt.fields.a2000,
				a1000:  tt.fields.a1000,
				a500:   tt.fields.a500,
				a100:   tt.fields.a100,
				a50:    tt.fields.a50,
				a10:    tt.fields.a10,
				a5:     tt.fields.a5,
				a1:     tt.fields.a1,
			}
			if got := c.GetTotalAmount(); got != tt.want {
				t.Errorf("GetTotalAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoneyCalculation_Subtract(t *testing.T) {
	type fields struct {
		m10000 int
		m5000  int
		m2000  int
		m1000  int
		m500   int
		m100   int
		m50    int
		m10    int
		m5     int
		m1     int
		s500   int
		s100   int
		s50    int
		s10    int
		s5     int
		s1     int
		a10000 int
		a5000  int
		a2000  int
		a1000  int
		a500   int
		a100   int
		a50    int
		a10    int
		a5     int
		a1     int
	}
	type args struct {
		minus [26]int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   [26]int
	}{
		{
			name: "引き算",
			fields: fields{
				m10000: 1,
				m5000:  2,
				m2000:  3,
				m1000:  4,
				m500:   5,
				m100:   6,
				m50:    7,
				m10:    8,
				m5:     9,
				m1:     10,
				s500:   11,
				s100:   12,
				s50:    13,
				s10:    14,
				s5:     15,
				s1:     16,
				a10000: 17,
				a5000:  18,
				a2000:  19,
				a1000:  20,
				a500:   21,
				a100:   22,
				a50:    23,
				a10:    24,
				a5:     25,
				a1:     26,
			},
			args: args{
				minus: [26]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			},
			want: [26]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MoneyCalculation{
				m10000: tt.fields.m10000,
				m5000:  tt.fields.m5000,
				m2000:  tt.fields.m2000,
				m1000:  tt.fields.m1000,
				m500:   tt.fields.m500,
				m100:   tt.fields.m100,
				m50:    tt.fields.m50,
				m10:    tt.fields.m10,
				m5:     tt.fields.m5,
				m1:     tt.fields.m1,
				s500:   tt.fields.s500,
				s100:   tt.fields.s100,
				s50:    tt.fields.s50,
				s10:    tt.fields.s10,
				s5:     tt.fields.s5,
				s1:     tt.fields.s1,
				a10000: tt.fields.a10000,
				a5000:  tt.fields.a5000,
				a2000:  tt.fields.a2000,
				a1000:  tt.fields.a1000,
				a500:   tt.fields.a500,
				a100:   tt.fields.a100,
				a50:    tt.fields.a50,
				a10:    tt.fields.a10,
				a5:     tt.fields.a5,
				a1:     tt.fields.a1,
			}
			if got := c.Subtract(tt.args.minus); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Subtract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoneyCalculation_add(t *testing.T) {
	type fields struct {
		m10000 int
		m5000  int
		m2000  int
		m1000  int
		m500   int
		m100   int
		m50    int
		m10    int
		m5     int
		m1     int
		s500   int
		s100   int
		s50    int
		s10    int
		s5     int
		s1     int
		a10000 int
		a5000  int
		a2000  int
		a1000  int
		a500   int
		a100   int
		a50    int
		a10    int
		a5     int
		a1     int
	}
	type args struct {
		m *MoneyCalculation
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *MoneyCalculation
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MoneyCalculation{
				m10000: tt.fields.m10000,
				m5000:  tt.fields.m5000,
				m2000:  tt.fields.m2000,
				m1000:  tt.fields.m1000,
				m500:   tt.fields.m500,
				m100:   tt.fields.m100,
				m50:    tt.fields.m50,
				m10:    tt.fields.m10,
				m5:     tt.fields.m5,
				m1:     tt.fields.m1,
				s500:   tt.fields.s500,
				s100:   tt.fields.s100,
				s50:    tt.fields.s50,
				s10:    tt.fields.s10,
				s5:     tt.fields.s5,
				s1:     tt.fields.s1,
				a10000: tt.fields.a10000,
				a5000:  tt.fields.a5000,
				a2000:  tt.fields.a2000,
				a1000:  tt.fields.a1000,
				a500:   tt.fields.a500,
				a100:   tt.fields.a100,
				a50:    tt.fields.a50,
				a10:    tt.fields.a10,
				a5:     tt.fields.a5,
				a1:     tt.fields.a1,
			}
			if got := c.add(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoneyCalculation_subtract(t *testing.T) {
	type fields struct {
		m10000 int
		m5000  int
		m2000  int
		m1000  int
		m500   int
		m100   int
		m50    int
		m10    int
		m5     int
		m1     int
		s500   int
		s100   int
		s50    int
		s10    int
		s5     int
		s1     int
		a10000 int
		a5000  int
		a2000  int
		a1000  int
		a500   int
		a100   int
		a50    int
		a10    int
		a5     int
		a1     int
	}
	type args struct {
		m *MoneyCalculation
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *MoneyCalculation
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MoneyCalculation{
				m10000: tt.fields.m10000,
				m5000:  tt.fields.m5000,
				m2000:  tt.fields.m2000,
				m1000:  tt.fields.m1000,
				m500:   tt.fields.m500,
				m100:   tt.fields.m100,
				m50:    tt.fields.m50,
				m10:    tt.fields.m10,
				m5:     tt.fields.m5,
				m1:     tt.fields.m1,
				s500:   tt.fields.s500,
				s100:   tt.fields.s100,
				s50:    tt.fields.s50,
				s10:    tt.fields.s10,
				s5:     tt.fields.s5,
				s1:     tt.fields.s1,
				a10000: tt.fields.a10000,
				a5000:  tt.fields.a5000,
				a2000:  tt.fields.a2000,
				a1000:  tt.fields.a1000,
				a500:   tt.fields.a500,
				a100:   tt.fields.a100,
				a50:    tt.fields.a50,
				a10:    tt.fields.a10,
				a5:     tt.fields.a5,
				a1:     tt.fields.a1,
			}
			if got := c.subtract(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("subtract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoneyCalculation_toExchange(t *testing.T) {
	type fields struct {
		m10000 int
		m5000  int
		m2000  int
		m1000  int
		m500   int
		m100   int
		m50    int
		m10    int
		m5     int
		m1     int
		s500   int
		s100   int
		s50    int
		s10    int
		s5     int
		s1     int
	}
	type args struct {
		amount int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *MoneyCalculation
	}{
		{
			name: "No.427 1000円札を1枚、500円玉を3枚、100円玉を6枚入金後の売上回収",
			fields: fields{
				m1000: 1,
				m500:  3,
				m100:  6,
			},
			args: args{
				amount: 3100,
			},
			want: &MoneyCalculation{
				m1000: 1, // 1000円分
				m500:  3, // 1500円分
				m100:  6, //  600円分
			},
		},
		{
			name: "カセットに金額がゼロの場合のテスト",
			fields: fields{
				m2000: 0,
				m1000: 0,
				m500:  0,
				m100:  0,
				m50:   0,
				m10:   0,
				m5:    0,
				m1:    0,
			},
			args: args{
				amount: 0,
			},
			want: &MoneyCalculation{},
		},
		{
			name: "負の金額が指定された場合のテスト",
			fields: fields{
				m2000: 5,
				m1000: 5,
				m500:  5,
				m100:  5,
				m50:   5,
				m10:   5,
				m5:    5,
				m1:    5,
			},
			args: args{
				amount: -5000,
			},
			want: &MoneyCalculation{
				m2000: -2, // -4000円分
				m1000: -1, // -1000円分
			},
		},
		{
			name: "2000円札優先チェック",
			fields: fields{
				m5000: 1,
				m2000: 1,
				m1000: 3,
			},
			args: args{
				amount: 7000,
			},
			want: &MoneyCalculation{
				m2000: 1, // 2000円分
				m5000: 1, // 5000円分
			},
		},
		{
			name: "カセット内の金額がちょうど要求された金額に等しい場合",
			fields: fields{
				m2000: 2,
				m1000: 1,
				m500:  1,
			},
			args: args{
				amount: 4500,
			},
			want: &MoneyCalculation{
				m2000: 2, // 4000円分
				m500:  1, // 500円分
			},
		},
		{
			name: "複数の金種が混在していて、要求された金額を正確に払い出せる場合",
			fields: fields{
				m2000: 3,
				m1000: 2,
				m500:  2,
				m100:  5,
				m50:   5,
				m10:   10,
				m5:    10,
				m1:    20,
			},
			args: args{
				amount: 8765,
			},
			want: &MoneyCalculation{
				m2000: 3, // 6000円分
				m1000: 2, // 2000円分
				m500:  1, // 500円分
				m100:  2, // 200円分
				m50:   1, // 50円分
				m10:   1, // 10円分
				m5:    1, // 5円分
			},
		},
		{
			name: "カセット内の金額が要求された金額よりも多い場合",
			fields: fields{
				m2000: 10,
				m1000: 10,
				m500:  10,
				m100:  10,
				m50:   10,
				m10:   10,
				m5:    10,
				m1:    10,
			},
			args: args{
				amount: 12345,
			},
			want: &MoneyCalculation{
				m2000: 6, // 12000円分
				m1000: 0, // 0円分
				m500:  0, // 0円分
				m100:  3, // 300円分
				m50:   0, // 0円分
				m10:   4, // 40円分
				m5:    1, // 5円分
			},
		},
		// 他のテストケースを追加することができます。
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MoneyCalculation{
				m10000: tt.fields.m10000,
				m5000:  tt.fields.m5000,
				m2000:  tt.fields.m2000,
				m1000:  tt.fields.m1000,
				m500:   tt.fields.m500,
				m100:   tt.fields.m100,
				m50:    tt.fields.m50,
				m10:    tt.fields.m10,
				m5:     tt.fields.m5,
				m1:     tt.fields.m1,
				s500:   tt.fields.s500,
				s100:   tt.fields.s100,
				s50:    tt.fields.s50,
				s10:    tt.fields.s10,
				s5:     tt.fields.s5,
				s1:     tt.fields.s1,
			}
			if got := c.toExchange(tt.args.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q. MoneyCalculation.toExchange() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestMoneyCalculation_toIntTbl26(t *testing.T) {
	type fields struct {
		m10000 int
		m5000  int
		m2000  int
		m1000  int
		m500   int
		m100   int
		m50    int
		m10    int
		m5     int
		m1     int
		s500   int
		s100   int
		s50    int
		s10    int
		s5     int
		s1     int
		a10000 int
		a5000  int
		a2000  int
		a1000  int
		a500   int
		a100   int
		a50    int
		a10    int
		a5     int
		a1     int
	}
	type args struct {
		i *MoneyCalculation
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   [26]int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MoneyCalculation{
				m10000: tt.fields.m10000,
				m5000:  tt.fields.m5000,
				m2000:  tt.fields.m2000,
				m1000:  tt.fields.m1000,
				m500:   tt.fields.m500,
				m100:   tt.fields.m100,
				m50:    tt.fields.m50,
				m10:    tt.fields.m10,
				m5:     tt.fields.m5,
				m1:     tt.fields.m1,
				s500:   tt.fields.s500,
				s100:   tt.fields.s100,
				s50:    tt.fields.s50,
				s10:    tt.fields.s10,
				s5:     tt.fields.s5,
				s1:     tt.fields.s1,
				a10000: tt.fields.a10000,
				a5000:  tt.fields.a5000,
				a2000:  tt.fields.a2000,
				a1000:  tt.fields.a1000,
				a500:   tt.fields.a500,
				a100:   tt.fields.a100,
				a50:    tt.fields.a50,
				a10:    tt.fields.a10,
				a5:     tt.fields.a5,
				a1:     tt.fields.a1,
			}
			if got := c.toIntTbl26(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toIntTbl26() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCassette(t *testing.T) {
	type args struct {
		moneyTbl [26]int
	}
	tests := []struct {
		name string
		args args
		want MoneyCalculationRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCassette(tt.args.moneyTbl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCassette() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toCassette(t *testing.T) {
	type args struct {
		before [26]int
	}
	tests := []struct {
		name string
		args args
		want *MoneyCalculation
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toCassette(tt.args.before); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toCassette() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoneyCalculation_toExchangeBill(t *testing.T) {
	type fields struct {
		m10000 int
		m5000  int
		m2000  int
		m1000  int
		m500   int
		m100   int
		m50    int
		m10    int
		m5     int
		m1     int
		s500   int
		s100   int
		s50    int
		s10    int
		s5     int
		s1     int
		a10000 int
		a5000  int
		a2000  int
		a1000  int
		a500   int
		a100   int
		a50    int
		a10    int
		a5     int
		a1     int
	}
	type args struct {
		amount int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *MoneyCalculation
	}{
		{
			name: "Normal",
			fields: fields{
				m10000: 1,
				m5000:  2,
				m2000:  3,
				m1000:  4,
			},
			args: args{
				amount: 5000,
			},
			want: &MoneyCalculation{
				m10000: 0,
				m5000:  0,
				m2000:  2,
				m1000:  1,
				m500:   0,
				m100:   0,
				m50:    0,
				m10:    0,
				m5:     0,
				m1:     0,
				s500:   0,
				s100:   0,
				s50:    0,
				s10:    0,
				s5:     0,
				s1:     0,
				a10000: 0,
				a5000:  0,
				a2000:  0,
				a1000:  0,
				a500:   0,
				a100:   0,
				a50:    0,
				a10:    0,
				a5:     0,
				a1:     0,
			},
		},
		{
			name: "10000 priority",
			fields: fields{
				m10000: 1,
				m5000:  2,
				m2000:  3,
				m1000:  3,
			},
			args: args{
				amount: 15000,
			},
			want: &MoneyCalculation{
				m10000: 1,
				m5000:  0,
				m2000:  2,
				m1000:  1,
				m500:   0,
				m100:   0,
				m50:    0,
				m10:    0,
				m5:     0,
				m1:     0,
				s500:   0,
				s100:   0,
				s50:    0,
				s10:    0,
				s5:     0,
				s1:     0,
				a10000: 0,
				a5000:  0,
				a2000:  0,
				a1000:  0,
				a500:   0,
				a100:   0,
				a50:    0,
				a10:    0,
				a5:     0,
				a1:     0,
			},
		},
		{
			name: "priorityLageBill",
			fields: fields{
				m10000: 3,
				m5000:  2,
				m2000:  6,
				m1000:  0,
			},
			args: args{
				amount: 15000,
			},
			want: &MoneyCalculation{
				m10000: 0,
				m5000:  1,
				m2000:  5,
				m1000:  0,
				m500:   0,
				m100:   0,
				m50:    0,
				m10:    0,
				m5:     0,
				m1:     0,
				s500:   0,
				s100:   0,
				s50:    0,
				s10:    0,
				s5:     0,
				s1:     0,
				a10000: 0,
				a5000:  0,
				a2000:  0,
				a1000:  0,
				a500:   0,
				a100:   0,
				a50:    0,
				a10:    0,
				a5:     0,
				a1:     0,
			},
		},
		{
			name: "priority2000V2",
			fields: fields{
				m10000: 3,
				m5000:  3,
				m2000:  3,
				m1000:  3,
			},
			args: args{
				amount: 10000,
			},
			want: &MoneyCalculation{
				m10000: 0,
				m5000:  1,
				m2000:  2,
				m1000:  1,
				m500:   0,
				m100:   0,
				m50:    0,
				m10:    0,
				m5:     0,
				m1:     0,
				s500:   0,
				s100:   0,
				s50:    0,
				s10:    0,
				s5:     0,
				s1:     0,
				a10000: 0,
				a5000:  0,
				a2000:  0,
				a1000:  0,
				a500:   0,
				a100:   0,
				a50:    0,
				a10:    0,
				a5:     0,
				a1:     0,
			},
		},
		{
			name: "priority2000 is 0",
			fields: fields{
				m10000: 3,
				m5000:  3,
				m2000:  0,
				m1000:  3,
			},
			args: args{
				amount: 10000,
			},
			want: &MoneyCalculation{
				m10000: 1,
				m5000:  0,
				m2000:  0,
				m1000:  0,
				m500:   0,
				m100:   0,
				m50:    0,
				m10:    0,
				m5:     0,
				m1:     0,
				s500:   0,
				s100:   0,
				s50:    0,
				s10:    0,
				s5:     0,
				s1:     0,
				a10000: 0,
				a5000:  0,
				a2000:  0,
				a1000:  0,
				a500:   0,
				a100:   0,
				a50:    0,
				a10:    0,
				a5:     0,
				a1:     0,
			},
		},
		{
			name: "priority2000 is 1",
			fields: fields{
				m10000: 3,
				m5000:  3,
				m2000:  1,
				m1000:  0,
			},
			args: args{
				amount: 5000,
			},
			want: &MoneyCalculation{
				m10000: 0,
				m5000:  1,
				m2000:  0,
				m1000:  0,
				m500:   0,
				m100:   0,
				m50:    0,
				m10:    0,
				m5:     0,
				m1:     0,
				s500:   0,
				s100:   0,
				s50:    0,
				s10:    0,
				s5:     0,
				s1:     0,
				a10000: 0,
				a5000:  0,
				a2000:  0,
				a1000:  0,
				a500:   0,
				a100:   0,
				a50:    0,
				a10:    0,
				a5:     0,
				a1:     0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MoneyCalculation{
				m10000: tt.fields.m10000,
				m5000:  tt.fields.m5000,
				m2000:  tt.fields.m2000,
				m1000:  tt.fields.m1000,
				m500:   tt.fields.m500,
				m100:   tt.fields.m100,
				m50:    tt.fields.m50,
				m10:    tt.fields.m10,
				m5:     tt.fields.m5,
				m1:     tt.fields.m1,
				s500:   tt.fields.s500,
				s100:   tt.fields.s100,
				s50:    tt.fields.s50,
				s10:    tt.fields.s10,
				s5:     tt.fields.s5,
				s1:     tt.fields.s1,
				a10000: tt.fields.a10000,
				a5000:  tt.fields.a5000,
				a2000:  tt.fields.a2000,
				a1000:  tt.fields.a1000,
				a500:   tt.fields.a500,
				a100:   tt.fields.a100,
				a50:    tt.fields.a50,
				a10:    tt.fields.a10,
				a5:     tt.fields.a5,
				a1:     tt.fields.a1,
			}
			if got := c.toExchangeBill(tt.args.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toExchangeBill() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoneyCalculation_overflowPriorityExchange(t *testing.T) {
	type fields struct {
		m10000 int
		m5000  int
		m2000  int
		m1000  int
		m500   int
		m100   int
		m50    int
		m10    int
		m5     int
		m1     int
		s500   int
		s100   int
		s50    int
		s10    int
		s5     int
		s1     int
		a10000 int
		a5000  int
		a2000  int
		a1000  int
		a500   int
		a100   int
		a50    int
		a10    int
		a5     int
		a1     int
	}
	type args struct {
		amount int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *MoneyCalculation
	}{
		{
			name: "A priority",
			fields: fields{
				m10000: 3,
				m5000:  3,
				m2000:  3,
				m1000:  3,
				m500:   5,
				m100:   5,
				m50:    5,
				m10:    5,
				m5:     5,
				m1:     5,
				s500:   0,
				s100:   5,
				s50:    0,
				s10:    5,
				s5:     0,
				s1:     5,
				a10000: 0,
				a5000:  0,
				a2000:  0,
				a1000:  0,
				a500:   0,
				a100:   1,
				a50:    0,
				a10:    0,
				a5:     0,
				a1:     0,
			},
			args: args{
				amount: 1100,
			},
			want: &MoneyCalculation{
				m10000: 0,
				m5000:  0,
				m2000:  0,
				m1000:  1,
				m500:   0,
				m100:   0,
				m50:    0,
				m10:    0,
				m5:     0,
				m1:     0,
				s500:   0,
				s100:   0,
				s50:    0,
				s10:    0,
				s5:     0,
				s1:     0,
				a10000: 0,
				a5000:  0,
				a2000:  0,
				a1000:  0,
				a500:   0,
				a100:   1,
				a50:    0,
				a10:    0,
				a5:     0,
				a1:     0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MoneyCalculation{
				m10000: tt.fields.m10000,
				m5000:  tt.fields.m5000,
				m2000:  tt.fields.m2000,
				m1000:  tt.fields.m1000,
				m500:   tt.fields.m500,
				m100:   tt.fields.m100,
				m50:    tt.fields.m50,
				m10:    tt.fields.m10,
				m5:     tt.fields.m5,
				m1:     tt.fields.m1,
				s500:   tt.fields.s500,
				s100:   tt.fields.s100,
				s50:    tt.fields.s50,
				s10:    tt.fields.s10,
				s5:     tt.fields.s5,
				s1:     tt.fields.s1,
				a10000: tt.fields.a10000,
				a5000:  tt.fields.a5000,
				a2000:  tt.fields.a2000,
				a1000:  tt.fields.a1000,
				a500:   tt.fields.a500,
				a100:   tt.fields.a100,
				a50:    tt.fields.a50,
				a10:    tt.fields.a10,
				a5:     tt.fields.a5,
				a1:     tt.fields.a1,
			}
			if got := c.overflowPriorityExchange(tt.args.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("overflowPriorityExchange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoneyCalculation_OverflowPriorityExchange(t *testing.T) {
	type fields struct {
		m10000 int
		m5000  int
		m2000  int
		m1000  int
		m500   int
		m100   int
		m50    int
		m10    int
		m5     int
		m1     int
		s500   int
		s100   int
		s50    int
		s10    int
		s5     int
		s1     int
		a10000 int
		a5000  int
		a2000  int
		a1000  int
		a500   int
		a100   int
		a50    int
		a10    int
		a5     int
		a1     int
	}
	type args struct {
		amount int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   [26]int
	}{
		{
			name: "A priority",
			fields: fields{
				m10000: 3,
				m5000:  3,
				m2000:  3,
				m1000:  3,
				m500:   5,
				m100:   5,
				m50:    5,
				m10:    5,
				m5:     5,
				m1:     5,
				s500:   0,
				s100:   5,
				s50:    0,
				s10:    5,
				s5:     0,
				s1:     5,
				a10000: 0,
				a5000:  0,
				a2000:  0,
				a1000:  0,
				a500:   0,
				a100:   1,
				a50:    0,
				a10:    0,
				a5:     0,
				a1:     0,
			},
			args: args{
				amount: 1100,
			},
			want: [26]int{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MoneyCalculation{
				m10000: tt.fields.m10000,
				m5000:  tt.fields.m5000,
				m2000:  tt.fields.m2000,
				m1000:  tt.fields.m1000,
				m500:   tt.fields.m500,
				m100:   tt.fields.m100,
				m50:    tt.fields.m50,
				m10:    tt.fields.m10,
				m5:     tt.fields.m5,
				m1:     tt.fields.m1,
				s500:   tt.fields.s500,
				s100:   tt.fields.s100,
				s50:    tt.fields.s50,
				s10:    tt.fields.s10,
				s5:     tt.fields.s5,
				s1:     tt.fields.s1,
				a10000: tt.fields.a10000,
				a5000:  tt.fields.a5000,
				a2000:  tt.fields.a2000,
				a1000:  tt.fields.a1000,
				a500:   tt.fields.a500,
				a100:   tt.fields.a100,
				a50:    tt.fields.a50,
				a10:    tt.fields.a10,
				a5:     tt.fields.a5,
				a1:     tt.fields.a1,
			}
			if got := c.OverflowPriorityExchange(tt.args.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OverflowPriorityExchange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoneyCalculation_toIntTbl10(t *testing.T) {
	type fields struct {
		m10000 int
		m5000  int
		m2000  int
		m1000  int
		m500   int
		m100   int
		m50    int
		m10    int
		m5     int
		m1     int
		s500   int
		s100   int
		s50    int
		s10    int
		s5     int
		s1     int
		a10000 int
		a5000  int
		a2000  int
		a1000  int
		a500   int
		a100   int
		a50    int
		a10    int
		a5     int
		a1     int
	}
	type args struct {
		i *MoneyCalculation
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   [10]int
	}{
		{
			name: "test from 26 to 10",
			args: args{
				i: &MoneyCalculation{
					m10000: 1,
					m5000:  2,
					m2000:  3,
					m1000:  4,
					m500:   5,
					m100:   6,
					m50:    7,
					m10:    8,
					m5:     9,
					m1:     10,
					s500:   5,
					s100:   4,
					s50:    3,
					s10:    2,
					s5:     1,
					s1:     0,
					a10000: 9,
					a5000:  8,
					a2000:  7,
					a1000:  6,
					a500:   5,
					a100:   4,
					a50:    3,
					a10:    2,
					a5:     1,
					a1:     10,
				},
			},
			want: [10]int{
				10, 10, 10, 10, 15, 14, 13, 12, 11, 20,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MoneyCalculation{
				m10000: tt.fields.m10000,
				m5000:  tt.fields.m5000,
				m2000:  tt.fields.m2000,
				m1000:  tt.fields.m1000,
				m500:   tt.fields.m500,
				m100:   tt.fields.m100,
				m50:    tt.fields.m50,
				m10:    tt.fields.m10,
				m5:     tt.fields.m5,
				m1:     tt.fields.m1,
				s500:   tt.fields.s500,
				s100:   tt.fields.s100,
				s50:    tt.fields.s50,
				s10:    tt.fields.s10,
				s5:     tt.fields.s5,
				s1:     tt.fields.s1,
				a10000: tt.fields.a10000,
				a5000:  tt.fields.a5000,
				a2000:  tt.fields.a2000,
				a1000:  tt.fields.a1000,
				a500:   tt.fields.a500,
				a100:   tt.fields.a100,
				a50:    tt.fields.a50,
				a10:    tt.fields.a10,
				a5:     tt.fields.a5,
				a1:     tt.fields.a1,
			}
			if got := c.toIntTbl10(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toIntTbl10() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoneyCalculation_GetOutCountTbl(t *testing.T) {
	type args struct {
		amount int
	}
	tests := []struct {
		name         string
		c            *MoneyCalculation
		args         args
		wantCountTbl [16]int
		wantResult   bool
	}{
		{
			name: "払出可能",
			c: &MoneyCalculation{
				m10000: 1,
				m5000:  2,
				m2000:  3,
				m1000:  4,
				m500:   5,
				m100:   6,
				m50:    7,
				m10:    8,
				m5:     9,
				m1:     10,
				s500:   11,
				s100:   12,
				s50:    13,
				s10:    14,
				s5:     15,
				s1:     16,
				a10000: 17,
				a5000:  18,
				a2000:  19,
				a1000:  20,
				a500:   21,
				a100:   22,
				a50:    23,
				a10:    24,
				a5:     25,
				a1:     26,
			},
			args: args{
				amount: 12345,
			},
			wantCountTbl: [16]int{1, 0, 1, 0, 0, 3, 0, 4, 1, 0, 0, 0, 0, 0, 0, 0},
			wantResult:   true,
		},
		{
			name: "置き換えて払出可能",
			c: &MoneyCalculation{
				m10000: 0,
				m5000:  1,
				m2000:  0,
				m1000:  0,
				m500:   0,
				m100:   4,
				m50:    0,
				m10:    5,
				m5:     0,
				m1:     0,
				s500:   0,
				s100:   1,
				s50:    0,
				s10:    0,
				s5:     0,
				s1:     0,
				a10000: 0,
				a5000:  0,
				a2000:  0,
				a1000:  0,
				a500:   0,
				a100:   0,
				a50:    0,
				a10:    0,
				a5:     0,
				a1:     0,
			},
			args: args{
				amount: 550,
			},
			wantCountTbl: [16]int{0, 0, 0, 0, 0, 5, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0},
			wantResult:   true,
		},
		{
			name: "置き換えて払出可能2",
			c: &MoneyCalculation{
				m10000: 4,
				m5000:  3,
				m2000:  0,
				m1000:  3,
				m500:   5,
				m100:   5,
				m50:    5,
				m10:    5,
				m5:     5,
				m1:     5,
				s500:   0,
				s100:   5,
				s50:    0,
				s10:    5,
				s5:     0,
				s1:     5,
				a10000: 0,
				a5000:  0,
				a2000:  0,
				a1000:  0,
				a500:   0,
				a100:   0,
				a50:    0,
				a10:    0,
				a5:     0,
				a1:     0,
			},
			args: args{
				amount: 9900,
			},
			wantCountTbl: [16]int{0, 1, 0, 3, 3, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantResult:   true,
		},
		{
			name: "置き換えて払出可能3",
			c: &MoneyCalculation{
				m10000: 0,
				m5000:  0,
				m2000:  0,
				m1000:  0,
				m500:   2,
				m100:   0,
				m50:    10,
				m10:    0,
				m5:     0,
				m1:     0,
				s500:   0,
				s100:   0,
				s50:    0,
				s10:    0,
				s5:     0,
				s1:     0,
				a10000: 0,
				a5000:  0,
				a2000:  0,
				a1000:  0,
				a500:   0,
				a100:   0,
				a50:    0,
				a10:    0,
				a5:     0,
				a1:     0,
			},
			args: args{
				amount: 1500,
			},
			wantCountTbl: [16]int{0, 0, 0, 0, 2, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantResult:   true,
		},
		{
			name: "払出可能（1万円の場合は下位3金種まで、千円優先）",
			c: &MoneyCalculation{
				m10000: 0,
				m5000:  1,
				m2000:  4,
				m1000:  3,
				m500:   4,
				m100:   0,
				m50:    0,
				m10:    0,
				m5:     0,
				m1:     0,
				s500:   0,
				s100:   0,
				s50:    0,
				s10:    0,
				s5:     0,
				s1:     0,
				a10000: 0,
				a5000:  0,
				a2000:  0,
				a1000:  0,
				a500:   0,
				a100:   0,
				a50:    0,
				a10:    0,
				a5:     0,
				a1:     0,
			},
			args: args{
				amount: 12000,
			},
			wantCountTbl: [16]int{0, 1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantResult:   true,
		},
		{
			name: "払出可能（1万円の場合は下位3金種まで、千円優先2）",
			c: &MoneyCalculation{
				m10000: 0,
				m5000:  0,
				m2000:  4,
				m1000:  15,
				m500:   4,
				m100:   5,
				m50:    0,
				m10:    0,
				m5:     0,
				m1:     0,
				s500:   0,
				s100:   0,
				s50:    0,
				s10:    0,
				s5:     0,
				s1:     0,
				a10000: 0,
				a5000:  0,
				a2000:  0,
				a1000:  0,
				a500:   0,
				a100:   0,
				a50:    0,
				a10:    0,
				a5:     0,
				a1:     0,
			},
			args: args{
				amount: 12100,
			},
			wantCountTbl: [16]int{0, 0, 0, 12, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantResult:   true,
		},
		{
			name: "払出可能（5千円の場合は下位3金種まで、千円優先）",
			c: &MoneyCalculation{
				m10000: 0,
				m5000:  1,
				m2000:  3,
				m1000:  0,
				m500:   4,
				m100:   0,
				m50:    0,
				m10:    0,
				m5:     0,
				m1:     0,
				s500:   0,
				s100:   0,
				s50:    0,
				s10:    0,
				s5:     0,
				s1:     0,
				a10000: 0,
				a5000:  0,
				a2000:  0,
				a1000:  0,
				a500:   0,
				a100:   0,
				a50:    0,
				a10:    0,
				a5:     0,
				a1:     0,
			},
			args: args{
				amount: 6000,
			},
			wantCountTbl: [16]int{0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantResult:   true,
		},
		{
			name: "払出可能（5千円の場合は下位3金種まで、千円優先2）",
			c: &MoneyCalculation{
				m10000: 0,
				m5000:  1,
				m2000:  3,
				m1000:  6,
				m500:   4,
				m100:   0,
				m50:    0,
				m10:    0,
				m5:     0,
				m1:     0,
				s500:   0,
				s100:   0,
				s50:    0,
				s10:    0,
				s5:     0,
				s1:     0,
				a10000: 0,
				a5000:  0,
				a2000:  0,
				a1000:  0,
				a500:   0,
				a100:   0,
				a50:    0,
				a10:    0,
				a5:     0,
				a1:     0,
			},
			args: args{
				amount: 6000,
			},
			wantCountTbl: [16]int{0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantResult:   true,
		},
		{
			name: "払出可能（5千円の場合は下位3金種まで）",
			c: &MoneyCalculation{
				m10000: 0,
				m5000:  0,
				m2000:  0,
				m1000:  0,
				m500:   10,
				m100:   0,
				m50:    0,
				m10:    0,
				m5:     0,
				m1:     0,
				s500:   0,
				s100:   0,
				s50:    0,
				s10:    0,
				s5:     0,
				s1:     0,
				a10000: 0,
				a5000:  0,
				a2000:  0,
				a1000:  0,
				a500:   0,
				a100:   0,
				a50:    0,
				a10:    0,
				a5:     0,
				a1:     0,
			},
			args: args{
				amount: 5000,
			},
			wantCountTbl: [16]int{0, 0, 0, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantResult:   false,
		},
		{
			name: "払出不可",
			c: &MoneyCalculation{
				m10000: 0,
				m5000:  1,
				m2000:  2,
				m1000:  0,
				m500:   4,
				m100:   0,
				m50:    0,
				m10:    0,
				m5:     0,
				m1:     0,
				s500:   0,
				s100:   0,
				s50:    0,
				s10:    0,
				s5:     0,
				s1:     0,
				a10000: 0,
				a5000:  0,
				a2000:  0,
				a1000:  0,
				a500:   0,
				a100:   0,
				a50:    0,
				a10:    0,
				a5:     0,
				a1:     0,
			},
			args: args{
				amount: 12345,
			},
			wantCountTbl: [16]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantResult:   false,
		},
		{
			name: "有高不足2回",
			c: &MoneyCalculation{
				m10000: 0,
				m5000:  1,
				m2000:  2,
				m1000:  0,
				m500:   0,
				m100:   0,
				m50:    0,
				m10:    50,
				m5:     0,
				m1:     0,
				s500:   0,
				s100:   0,
				s50:    0,
				s10:    0,
				s5:     0,
				s1:     0,
				a10000: 0,
				a5000:  0,
				a2000:  0,
				a1000:  0,
				a500:   0,
				a100:   0,
				a50:    0,
				a10:    0,
				a5:     0,
				a1:     0,
			},
			args: args{
				amount: 500,
			},
			wantCountTbl: [16]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantResult:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCountTbl, gotResult := tt.c.GetOutCountTbl(tt.args.amount)
			if !tt.wantResult {
				if gotResult != tt.wantResult {
					t.Errorf("MoneyCalculation.GetOutCountTbl() gotResult = %v, want %v", gotResult, tt.wantResult)
				}
			} else {
				if !reflect.DeepEqual(gotCountTbl, tt.wantCountTbl) {
					t.Errorf("MoneyCalculation.GetOutCountTbl() gotCountTbl = %v, want %v", gotCountTbl, tt.wantCountTbl)
				}
				if gotResult != tt.wantResult {
					t.Errorf("MoneyCalculation.GetOutCountTbl() gotResult = %v, want %v", gotResult, tt.wantResult)
				}
			}
		})
	}
}
