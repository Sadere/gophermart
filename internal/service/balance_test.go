package service

import (
	"testing"

	"github.com/Sadere/gophermart/internal/model"
	"github.com/Sadere/gophermart/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestRegisterWithdraw(t *testing.T) {
	repo := &repository.TestBalanceRepository{}
	balanceService := NewBalanceService(repo)

	tests := []struct {
		name    string
		userID  uint64
		number  string
		sum     float64
		wantErr bool
	}{
		{
			name:    "success withdrawal register",
			userID:  111,
			number:  "89920",
			sum:     100,
			wantErr: false,
		},
		{
			name:    "invalid withdraw order number",
			userID:  111,
			number:  "111",
			sum:     100,
			wantErr: true,
		},
		{
			name:    "no funds",
			userID:  222,
			number:  "89920",
			sum:     100,
			wantErr: true,
		},
		{
			name:    "error on getting balance",
			userID:  333,
			number:  "89920",
			sum:     100,
			wantErr: true,
		},
		{
			name:    "no funds on withdraw",
			userID:  444,
			number:  "89920",
			sum:     50,
			wantErr: true,
		},
		{
			name:    "error withdraw",
			userID:  555,
			number:  "89920",
			sum:     50,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := balanceService.RegisterWithdraw(tt.userID, tt.number, tt.sum)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListUserWithdrawals(t *testing.T) {
	repo := &repository.TestBalanceRepository{}
	balanceService := NewBalanceService(repo)

	type want struct {
		len int
		err bool
	}
	tests := []struct {
		name   string
		userID uint64
		want   want
	}{
		{
			name:   "success list",
			userID: 111,
			want: want{
				len: 1,
				err: false,
			},
		},
		{
			name:   "error",
			userID: 222,
			want: want{
				len: 0,
				err: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withdrawals, err := balanceService.ListUserWithdrawals(tt.userID)

			assert.Len(t, withdrawals, tt.want.len)

			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUserBalance(t *testing.T) {
	repo := &repository.TestBalanceRepository{}
	balanceService := NewBalanceService(repo)

	type want struct {
		balance model.UserBalance
		err     bool
	}
	tests := []struct {
		name   string
		userID uint64
		want   want
	}{
		{
			name:   "success get user balance",
			userID: 111,
			want: want{
				balance: model.UserBalance{
					Balance:   200,
					Withdrawn: 200,
				},
				err: false,
			},
		},
		{
			name:   "error",
			userID: 333,
			want: want{
				err: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			balance, err := balanceService.GetUserBalance(tt.userID)

			if tt.want.err {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if assert.NotNil(t, balance) {
				assert.Equal(t, tt.want.balance.Balance, (*balance).Balance)
				assert.Equal(t, tt.want.balance.Withdrawn, (*balance).Withdrawn)
			}
		})
	}
}
