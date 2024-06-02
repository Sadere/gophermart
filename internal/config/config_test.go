package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/Sadere/gophermart/internal/structs"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	var tests = []struct {
		name string
		args []string
		env  map[string]string
		conf Config
	}{
		{
			name: "address",
			args: []string{"-a", "localhost:1337"},
			env: map[string]string{
				"SECRET_KEY": "test",
			},
			conf: Config{
				Address: structs.NetAddress{
					Host: "localhost",
					Port: 1337,
				},
				SecretKey:    "test",
				PullInterval: DefaultPullInterval,
			},
		},
		{
			name: "address env priority",
			args: []string{"-a", "localhost:1337"},
			env: map[string]string{
				"SECRET_KEY":  "test",
				"RUN_ADDRESS": "testhost:2222",
			},
			conf: Config{
				Address: structs.NetAddress{
					Host: "testhost",
					Port: 2222,
				},
				SecretKey:    "test",
				PullInterval: DefaultPullInterval,
			},
		},
		{
			name: "accrual address",
			args: []string{"-a", "localhost:1337", "-r", "accrual:8888"},
			env: map[string]string{
				"SECRET_KEY": "test",
			},
			conf: Config{
				Address: structs.NetAddress{
					Host: "localhost",
					Port: 1337,
				},
				AccrualAddr: structs.NetAddress{
					Host: "accrual",
					Port: 8888,
				},
				SecretKey:    "test",
				PullInterval: DefaultPullInterval,
			},
		},
		{
			name: "accrual address env",
			args: []string{"-a", "localhost:1337", "-r", "accrual:8888"},
			env: map[string]string{
				"SECRET_KEY":             "test",
				"ACCRUAL_SYSTEM_ADDRESS": "accrual22:9999",
			},
			conf: Config{
				Address: structs.NetAddress{
					Host: "localhost",
					Port: 1337,
				},
				AccrualAddr: structs.NetAddress{
					Host: "accrual22",
					Port: 9999,
				},
				SecretKey:    "test",
				PullInterval: DefaultPullInterval,
			},
		},
		{
			name: "dsn from arg",
			args: []string{"-a", "localhost:1337", "-d", "dsn_test"},
			env: map[string]string{
				"SECRET_KEY": "test",
			},
			conf: Config{
				Address: structs.NetAddress{
					Host: "localhost",
					Port: 1337,
				},
				PostgresDSN:  "dsn_test",
				SecretKey:    "test",
				PullInterval: DefaultPullInterval,
			},
		},
		{
			name: "dsn from env",
			args: []string{"-a", "localhost:1337", "-d", "dsn_test"},
			env: map[string]string{
				"SECRET_KEY":   "test",
				"DATABASE_URI": "000dsn_test000",
			},
			conf: Config{
				Address: structs.NetAddress{
					Host: "localhost",
					Port: 1337,
				},
				PostgresDSN:  "000dsn_test000",
				SecretKey:    "test",
				PullInterval: DefaultPullInterval,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.env {
				os.Setenv(key, value)
			}

			conf, err := NewConfig(tt.args)

			assert.NoError(t, err)

			if !reflect.DeepEqual(conf, tt.conf) {
				t.Errorf("conf got %+v, want %+v", conf, tt.conf)
			}

			// Remove envs for next tests
			for key := range tt.env {
				os.Unsetenv(key)
			}
		})
	}
}
