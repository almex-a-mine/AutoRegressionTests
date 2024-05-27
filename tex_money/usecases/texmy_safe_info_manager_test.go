package usecases

import (
	"reflect"
	"testing"
	"tex_money/config"
	"tex_money/domain"
	"tex_money/domain/handler"
)

func Test_safeInfoManager_updateChangeAvailable(t *testing.T) {
	type fields struct {
		logger     handler.LoggerRepository
		cfg        config.Configuration
		syslogMng  SyslogManager
		errorMng   ErrorManager
		SafeInfo   domain.SafeInfo
		iniService IniServiceRepository
	}
	type args struct {
		s domain.SortInfoTbl
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   domain.SortInfoTbl
	}{
		{
			name:   "test",
			fields: fields{},
			args: args{
				s: domain.SortInfoTbl{
					CountTbl:   [domain.CASH_TYPE_SHITEI]int{36, 34, 32, 30, 44, 41, 38, 35, 32, 29},
					ExCountTbl: [domain.EXTRA_CASH_TYPE_SHITEI]int{26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
				},
			},

			want: domain.SortInfoTbl{
				SortType:   domain.CHANGE_AVAILABLE,
				Amount:     480798,
				CountTbl:   [domain.CASH_TYPE_SHITEI]int{26, 25, 24, 23, 38, 36, 34, 32, 30, 28},
				ExCountTbl: [domain.EXTRA_CASH_TYPE_SHITEI]int{26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &safeInfoManager{
				logger:     tt.fields.logger,
				cfg:        tt.fields.cfg,
				syslogMng:  tt.fields.syslogMng,
				errorMng:   tt.fields.errorMng,
				SafeInfo:   tt.fields.SafeInfo,
				iniService: tt.fields.iniService,
			}
			if got := c.updateChangeAvailable(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updateChangeAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}
