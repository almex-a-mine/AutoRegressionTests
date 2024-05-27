package usecases

import (
	"reflect"
	"testing"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
)

func Test_reverseExchangeCalculationManager_baseExchange(t *testing.T) {
	type fields struct {
		logger handler.LoggerRepository
		safe   SafeInfoManager
		config config.Configuration
	}
	type args struct {
		exchangeType int
		setup        [domain.EXTRA_CASH_TYPE_SHITEI]int
		safeOne      [domain.EXTRA_CASH_TYPE_SHITEI]int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  [domain.EXTRA_CASH_TYPE_SHITEI]int
		want2  [domain.EXTRA_CASH_TYPE_SHITEI]int
	}{
		{
			name:   "11_MainCassette_return Coin",
			fields: fields{},
			args: args{
				exchangeType: 11,
				setup:        [domain.EXTRA_CASH_TYPE_SHITEI]int{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10},
				safeOne:      [domain.EXTRA_CASH_TYPE_SHITEI]int{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
			},
			want:  666,
			want1: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want2: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:   "11_MainCassette_return Bill_Coin",
			fields: fields{},
			args: args{
				exchangeType: 11,
				setup:        [domain.EXTRA_CASH_TYPE_SHITEI]int{10, 10, 10, 10, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100},
				safeOne:      [domain.EXTRA_CASH_TYPE_SHITEI]int{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
			},
			want:  60606,
			want1: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 91, 91, 91, 91, 91, 91, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want2: [domain.EXTRA_CASH_TYPE_SHITEI]int{4, 0, 9, 2, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:   "12_MainCassette_return Coin",
			fields: fields{},
			args: args{
				exchangeType: 12,
				setup:        [domain.EXTRA_CASH_TYPE_SHITEI]int{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10},
				safeOne:      [domain.EXTRA_CASH_TYPE_SHITEI]int{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
			},
			want:  666,
			want1: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want2: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:   "12_MainCassette_return Bill_Coin",
			fields: fields{},
			args: args{
				exchangeType: 12,
				setup:        [domain.EXTRA_CASH_TYPE_SHITEI]int{10, 10, 10, 10, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100},
				safeOne:      [domain.EXTRA_CASH_TYPE_SHITEI]int{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
			},
			want:  60606,
			want1: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 91, 91, 91, 91, 91, 91, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want2: [domain.EXTRA_CASH_TYPE_SHITEI]int{4, 0, 9, 2, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:   "13_MainCassette_return Coin",
			fields: fields{},
			args: args{
				exchangeType: 13,
				setup:        [domain.EXTRA_CASH_TYPE_SHITEI]int{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10},
				safeOne:      [domain.EXTRA_CASH_TYPE_SHITEI]int{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
			},
			want:  1332,
			want1: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want2: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 1, 0, 3, 0, 3, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:   "13_MainCassette_return Bill_Coin",
			fields: fields{},
			args: args{
				exchangeType: 13,
				setup:        [domain.EXTRA_CASH_TYPE_SHITEI]int{10, 10, 10, 10, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100},
				safeOne:      [domain.EXTRA_CASH_TYPE_SHITEI]int{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
			},
			want:  121212,
			want1: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 91, 91, 91, 91, 91, 91, 91, 91, 91, 91, 91, 91, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want2: [domain.EXTRA_CASH_TYPE_SHITEI]int{9, 2, 9, 3, 0, 2, 0, 1, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:   "16_bill_and_coin 釣銭準備金超過（プラス値）",
			fields: fields{},
			args: args{
				exchangeType: 16,
				setup:        [domain.EXTRA_CASH_TYPE_SHITEI]int{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				safeOne:      [domain.EXTRA_CASH_TYPE_SHITEI]int{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 0, 0, 0, 0, 0, 0, 9, 9, 9, 9, 0, 0, 0, 0, 0, 0},
			},
			want:  -143334,
			want1: [domain.EXTRA_CASH_TYPE_SHITEI]int{-8, -8, -8, -8, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want2: [domain.EXTRA_CASH_TYPE_SHITEI]int{-10, -1, -18, -2, 0, -3, 0, -3, 0, -4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:   "16_bill_and_coin 釣銭準備金より不足（マイナス値）",
			fields: fields{},
			args: args{
				exchangeType: 16,
				setup:        [domain.EXTRA_CASH_TYPE_SHITEI]int{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				safeOne:      [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 9, 9, 9, 9, 9, 9, 0, 0, 0, 0, 0, 0, 9, 9, 9, 9, 0, 0, 0, 0, 0, 0},
			},
			want:  18666,
			want1: [domain.EXTRA_CASH_TYPE_SHITEI]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want2: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 9, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &reverseExchangeCalculationManager{
				logger: tt.fields.logger,
				safe:   tt.fields.safe,
				config: tt.fields.config,
			}
			got, got1, got2 := c.baseExchange(tt.args.exchangeType, tt.args.setup, tt.args.safeOne)
			if got != tt.want {
				t.Errorf("baseExchange() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("baseExchange() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("baseExchange() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_reverseExchangeCalculationManager_specifyExchange(t *testing.T) {
	type fields struct {
		logger handler.LoggerRepository
		safe   SafeInfoManager
		config config.Configuration
	}
	type args struct {
		amount  int
		safeOne [domain.EXTRA_CASH_TYPE_SHITEI]int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   [domain.EXTRA_CASH_TYPE_SHITEI]int
	}{
		{
			name: "30_666",
			args: args{
				amount:  666,
				safeOne: [domain.EXTRA_CASH_TYPE_SHITEI]int{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
			},
			want: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "30_121212",
			args: args{
				amount:  121212,
				safeOne: [domain.EXTRA_CASH_TYPE_SHITEI]int{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
			},
			want: [domain.EXTRA_CASH_TYPE_SHITEI]int{9, 2, 9, 3, 0, 2, 0, 1, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &reverseExchangeCalculationManager{
				logger: tt.fields.logger,
				safe:   tt.fields.safe,
				config: tt.fields.config,
			}
			if got := c.specifyExchange(tt.args.amount, tt.args.safeOne); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("specifyExchange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reverseExchangeCalculationManager_salesMoneyExchange(t *testing.T) {
	type fields struct {
		logger handler.LoggerRepository
		safe   SafeInfoManager
		config config.Configuration
	}
	type args struct {
		salesAmount int
		overflowBox bool
		safeOne     [domain.EXTRA_CASH_TYPE_SHITEI]int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   [domain.EXTRA_CASH_TYPE_SHITEI]int
	}{
		{
			name: "Sales-true-167994",
			args: args{
				salesAmount: 167994,
				overflowBox: true,
				safeOne:     [domain.EXTRA_CASH_TYPE_SHITEI]int{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
			},
			want: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
		},
		{
			name: "Sales-true-1100",
			args: args{
				salesAmount: 1100,
				overflowBox: true,
				safeOne:     [domain.EXTRA_CASH_TYPE_SHITEI]int{9, 9, 9, 9, 5, 5, 5, 5, 5, 5, 0, 5, 0, 5, 0, 5, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
			},
			want: [domain.EXTRA_CASH_TYPE_SHITEI]int{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
		},
		{
			name: "Sales-false-167994",
			args: args{
				salesAmount: 167994,
				overflowBox: false,
				safeOne:     [domain.EXTRA_CASH_TYPE_SHITEI]int{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
			},
			want: [domain.EXTRA_CASH_TYPE_SHITEI]int{9, 9, 9, 9, 9, 4, 1, 4, 0, 4, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &reverseExchangeCalculationManager{
				logger: tt.fields.logger,
				safe:   tt.fields.safe,
				config: tt.fields.config,
			}
			if got := c.salesMoneyExchange(tt.args.salesAmount, tt.args.overflowBox, tt.args.safeOne); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("salesMoneyExchange() = %v, want %v", got, tt.want)
			}
		})
	}
}
