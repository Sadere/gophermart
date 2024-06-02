package service

import (
	"testing"

	"github.com/Sadere/gophermart/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestSaveOrderForUser(t *testing.T) {
	repo := &repository.TestOrderRepository{}
	orderService := NewOrderService(repo)

	type want struct {
		exists bool
		err    bool
	}

	tests := []struct {
		name   string
		userID uint64
		number string
		want   want
	}{
		{
			name:   "success save",
			userID: 111,
			number: "89920",
			want: want{
				exists: false,
				err:    false,
			},
		},
		{
			name:   "invalid order number",
			userID: 111,
			number: "111",
			want: want{
				exists: false,
				err:    true,
			},
		},
		{
			name:   "error on create",
			userID: 111,
			number: "27078",
			want: want{
				exists: false,
				err:    true,
			},
		},
		{
			name:   "error on getting order",
			userID: 111,
			number: "43513",
			want: want{
				exists: false,
				err:    true,
			},
		},
		{
			name:   "order already exists",
			userID: 111,
			number: "56317",
			want: want{
				exists: true,
				err:    false,
			},
		},
		{
			name:   "order loaded by another user",
			userID: 222,
			number: "56317",
			want: want{
				exists: true,
				err:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := orderService.SaveOrderForUser(tt.userID, tt.number)

			if tt.want.exists {
				assert.True(t, exists)
			} else {
				assert.False(t, exists)
			}

			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetOrdersByUser(t *testing.T) {
	repo := &repository.TestOrderRepository{}
	orderService := NewOrderService(repo)

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
		{
			name:   "empty list",
			userID: 333,
			want: want{
				len: 0,
				err: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders, err := orderService.GetOrdersByUser(tt.userID)

			assert.Len(t, orders, tt.want.len)

			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
