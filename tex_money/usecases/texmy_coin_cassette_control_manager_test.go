package usecases

import (
	"reflect"
	"testing"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
)

func TestCoinCassetteControlUsecases_collection(t *testing.T) {
	type fields struct {
		logger handler.LoggerRepository
		safe   SafeInfoManager
		config config.Configuration
	}
	type args struct {
		cassette int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	texCon := domain.NewTexContext(domain.RegisterTexContext{})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &coinCassetteControlManager{
				logger: tt.fields.logger,
				safe:   tt.fields.safe,
				config: tt.fields.config,
			}
			c.Collection(texCon, tt.args.cassette)
		})
	}
}

func TestCoinCassetteControlUsecases_ClearMainCassette(t *testing.T) {
	type fields struct {
		logger handler.LoggerRepository
		safe   SafeInfoManager
		config config.Configuration
	}
	type args struct {
		i *domain.Cassette
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *domain.Cassette
	}{
		{
			name:   "メインカセットクリア",
			fields: fields{},
			args: args{
				i: &domain.Cassette{
					M10000: 10,
					M5000:  9,
					M2000:  8,
					M1000:  7,
					M500:   6,
					M100:   5,
					M50:    4,
					M10:    3,
					M5:     2,
					M1:     1,
					S500:   6,
					S100:   5,
					S50:    4,
					S10:    3,
					S5:     2,
					S1:     1,
					A10000: -5,
					A5000:  -4,
					A2000:  -3,
					A1000:  -2,
					A500:   -1,
					A100:   -4,
					A50:    -3,
					A10:    -2,
					A5:     -1,
					A1:     0,
				},
			},
			want: &domain.Cassette{
				M10000: 10,
				M5000:  9,
				M2000:  8,
				M1000:  7,
				M500:   0,
				M100:   0,
				M50:    0,
				M10:    0,
				M5:     0,
				M1:     0,
				S500:   6,
				S100:   5,
				S50:    4,
				S10:    3,
				S5:     2,
				S1:     1,
				A10000: -5,
				A5000:  -4,
				A2000:  -3,
				A1000:  -2,
				A500:   -1,
				A100:   -4,
				A50:    -3,
				A10:    -2,
				A5:     -1,
				A1:     0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &coinCassetteControlManager{
				logger: tt.fields.logger,
				safe:   tt.fields.safe,
				config: tt.fields.config,
			}
			if got := c.ClearMainCassette(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClearMainCassette() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoinCassetteControlUsecases_ClearSubCassette(t *testing.T) {
	type fields struct {
		logger handler.LoggerRepository
		safe   SafeInfoManager
		config config.Configuration
	}
	type args struct {
		i *domain.Cassette
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *domain.Cassette
	}{
		{
			name:   "サブカセットクリア",
			fields: fields{},
			args: args{
				i: &domain.Cassette{
					M10000: 10,
					M5000:  9,
					M2000:  8,
					M1000:  7,
					M500:   6,
					M100:   5,
					M50:    4,
					M10:    3,
					M5:     2,
					M1:     1,
					S500:   6,
					S100:   5,
					S50:    4,
					S10:    3,
					S5:     2,
					S1:     1,
					A10000: -5,
					A5000:  -4,
					A2000:  -3,
					A1000:  -2,
					A500:   -1,
					A100:   -4,
					A50:    -3,
					A10:    -2,
					A5:     -1,
					A1:     0,
				},
			},
			want: &domain.Cassette{
				M10000: 10,
				M5000:  9,
				M2000:  8,
				M1000:  7,
				M500:   6,
				M100:   5,
				M50:    4,
				M10:    3,
				M5:     2,
				M1:     1,
				S500:   0,
				S100:   0,
				S50:    0,
				S10:    0,
				S5:     0,
				S1:     0,
				A10000: -5,
				A5000:  -4,
				A2000:  -3,
				A1000:  -2,
				A500:   -1,
				A100:   -4,
				A50:    -3,
				A10:    -2,
				A5:     -1,
				A1:     0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &coinCassetteControlManager{
				logger: tt.fields.logger,
				safe:   tt.fields.safe,
				config: tt.fields.config,
			}
			if got := c.ClearSubCassette(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClearSubCassette() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoinCassetteControlUsecases_Subtract(t *testing.T) {
	type fields struct {
		logger handler.LoggerRepository
		safe   SafeInfoManager
		config config.Configuration
	}
	type args struct {
		after  *domain.Cassette
		before *domain.Cassette
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *domain.Cassette
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &coinCassetteControlManager{
				logger: tt.fields.logger,
				safe:   tt.fields.safe,
				config: tt.fields.config,
			}
			if got := c.subtract(tt.args.after, tt.args.before); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("subtract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoinCassetteControlUsecases_change(t *testing.T) {
	type fields struct {
		logger handler.LoggerRepository
		safe   SafeInfoManager
		config config.Configuration
	}
	type args struct {
		cassette int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	texCon := domain.NewTexContext(domain.RegisterTexContext{})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &coinCassetteControlManager{
				logger: tt.fields.logger,
				safe:   tt.fields.safe,
				config: tt.fields.config,
			}
			c.Exchange(texCon, tt.args.cassette)
		})
	}
}

func TestCoinCassetteControlUsecases_specificationReplenishment(t *testing.T) {
	type fields struct {
		logger handler.LoggerRepository
		safe   SafeInfoManager
		config config.Configuration
	}
	type args struct {
		cassette    int
		amountCount [domain.EXTRA_CASH_TYPE_SHITEI]int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	texCon := domain.NewTexContext(domain.RegisterTexContext{})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &coinCassetteControlManager{
				logger: tt.fields.logger,
				safe:   tt.fields.safe,
				config: tt.fields.config,
			}
			c.SpecificationReplenishment(texCon, tt.args.cassette, tt.args.amountCount)
		})
	}
}

func TestCoinCassetteControlUsecases_toCassette(t *testing.T) {
	type fields struct {
		logger handler.LoggerRepository
		safe   SafeInfoManager
		config config.Configuration
	}
	type args struct {
		i [domain.EXTRA_CASH_TYPE_SHITEI]int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   domain.Cassette
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &coinCassetteControlManager{
				logger: tt.fields.logger,
				safe:   tt.fields.safe,
				config: tt.fields.config,
			}
			if got := c.toCassette(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toCassette() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoinCassetteControlUsecases_totalAmount(t *testing.T) {
	type fields struct {
		logger handler.LoggerRepository
		safe   SafeInfoManager
		config config.Configuration
	}
	type args struct {
		i *domain.Cassette
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name:   "計算テスト",
			fields: fields{},
			args: args{
				i: &domain.Cassette{
					M10000: 10,
					M5000:  9,
					M2000:  8,
					M1000:  7,
					M500:   6,
					M100:   5,
					M50:    4,
					M10:    3,
					M5:     2,
					M1:     1,
					S500:   6,
					S100:   5,
					S50:    4,
					S10:    3,
					S5:     2,
					S1:     1,
					A10000: -5,
					A5000:  -4,
					A2000:  -3,
					A1000:  -2,
					A500:   -1,
					A100:   -4,
					A50:    -3,
					A10:    -2,
					A5:     -1,
					A1:     0,
				},
			},
			want: 96407,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &coinCassetteControlManager{
				logger: tt.fields.logger,
				safe:   tt.fields.safe,
				config: tt.fields.config,
			}
			if got := c.toTalAmount(tt.args.i); got != tt.want {
				t.Errorf("toTalAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoinCassetteControlUsecases_toExchange(t *testing.T) {
	type fields struct {
		logger handler.LoggerRepository
		safe   SafeInfoManager
		config config.Configuration
	}
	type args struct {
		amount   int
		cassette *domain.Cassette
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *domain.Cassette
	}{
		{
			name: "両替全金種対象",
			args: args{
				amount: 41166,
				cassette: &domain.Cassette{
					M10000: 1,
					M5000:  2,
					M2000:  3,
					M1000:  4,
					M500:   5,
					M100:   6,
					M50:    7,
					M10:    8,
					M5:     9,
					M1:     10,
					S500:   11,
					S100:   12,
					S50:    13,
					S10:    14,
					S5:     15,
					S1:     16,
					A10000: -5,
					A5000:  -4,
					A2000:  -3,
					A1000:  -2,
					A500:   -1,
					A100:   -4,
					A50:    -3,
					A10:    -2,
					A5:     -1,
					A1:     0,
				},
			},
			want: &domain.Cassette{
				M10000: 1,
				M5000:  2,
				M2000:  3,
				M1000:  4,
				M500:   5,
				M100:   6,
				M50:    7,
				M10:    8,
				M5:     9,
				M1:     10,
				S500:   11,
				S100:   12,
				S50:    13,
				S10:    14,
				S5:     15,
				S1:     16,
				A10000: 0,
				A5000:  0,
				A2000:  0,
				A1000:  0,
				A500:   0,
				A100:   0,
				A50:    0,
				A10:    0,
				A5:     0,
				A1:     0,
			},
		},
		{
			name: "両替 100円 メイン＆サブ",
			args: args{
				amount: 1000,
				cassette: &domain.Cassette{
					M10000: 0,
					M5000:  0,
					M2000:  0,
					M1000:  0,
					M500:   0,
					M100:   6,
					M50:    10,
					M10:    10,
					M5:     10,
					M1:     10,
					S500:   0,
					S100:   4,
					S50:    10,
					S10:    10,
					S5:     10,
					S1:     10,
					A10000: -5,
					A5000:  -4,
					A2000:  -3,
					A1000:  -2,
					A500:   -1,
					A100:   -4,
					A50:    -3,
					A10:    -2,
					A5:     -1,
					A1:     0,
				},
			},
			want: &domain.Cassette{
				M10000: 0,
				M5000:  0,
				M2000:  0,
				M1000:  0,
				M500:   0,
				M100:   6,
				M50:    0,
				M10:    0,
				M5:     0,
				M1:     0,
				S500:   0,
				S100:   4,
				S50:    0,
				S10:    0,
				S5:     0,
				S1:     0,
				A10000: 0,
				A5000:  0,
				A2000:  0,
				A1000:  0,
				A500:   0,
				A100:   0,
				A50:    0,
				A10:    0,
				A5:     0,
				A1:     0,
			},
		},
		{
			name: "両替 M100&S100+S50",
			args: args{
				amount: 1050,
				cassette: &domain.Cassette{
					M10000: 0,
					M5000:  0,
					M2000:  0,
					M1000:  0,
					M500:   0,
					M100:   6,
					M50:    0,
					M10:    10,
					M5:     10,
					M1:     10,
					S500:   0,
					S100:   4,
					S50:    10,
					S10:    10,
					S5:     10,
					S1:     10,
					A10000: -5,
					A5000:  -4,
					A2000:  -3,
					A1000:  -2,
					A500:   -1,
					A100:   -4,
					A50:    -3,
					A10:    -2,
					A5:     -1,
					A1:     0,
				},
			},
			want: &domain.Cassette{
				M10000: 0,
				M5000:  0,
				M2000:  0,
				M1000:  0,
				M500:   0,
				M100:   6,
				M50:    0,
				M10:    0,
				M5:     0,
				M1:     0,
				S500:   0,
				S100:   4,
				S50:    1,
				S10:    0,
				S5:     0,
				S1:     0,
				A10000: 0,
				A5000:  0,
				A2000:  0,
				A1000:  0,
				A500:   0,
				A100:   0,
				A50:    0,
				A10:    0,
				A5:     0,
				A1:     0,
			},
		},
		{
			name: "両替全金種対象&対象金額がマイナス",
			args: args{
				amount: -41166,
				cassette: &domain.Cassette{
					M10000: 1,
					M5000:  2,
					M2000:  3,
					M1000:  4,
					M500:   5,
					M100:   6,
					M50:    7,
					M10:    8,
					M5:     9,
					M1:     10,
					S500:   11,
					S100:   12,
					S50:    13,
					S10:    14,
					S5:     15,
					S1:     16,
					A10000: -5,
					A5000:  -4,
					A2000:  -3,
					A1000:  -2,
					A500:   -1,
					A100:   -4,
					A50:    -3,
					A10:    -2,
					A5:     -1,
					A1:     0,
				},
			},
			want: &domain.Cassette{
				M10000: -1,
				M5000:  -2,
				M2000:  -3,
				M1000:  -4,
				M500:   -5,
				M100:   -6,
				M50:    -7,
				M10:    -8,
				M5:     -9,
				M1:     -10,
				S500:   -11,
				S100:   -12,
				S50:    -13,
				S10:    -14,
				S5:     -15,
				S1:     -16,
				A10000: 0,
				A5000:  0,
				A2000:  0,
				A1000:  0,
				A500:   0,
				A100:   0,
				A50:    0,
				A10:    0,
				A5:     0,
				A1:     0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &coinCassetteControlManager{
				logger: tt.fields.logger,
				safe:   tt.fields.safe,
				config: tt.fields.config,
			}
			if got := c.toExchange(tt.args.amount, tt.args.cassette); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toExchange() = %v, want %v", got, tt.want)
			}

		})
	}
}
