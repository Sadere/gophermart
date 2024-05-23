package structs

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Адрес в формате <хост>:<порт>
type NetAddress struct {
	Host string
	Port int
}

func (addr *NetAddress) String() string {
	return fmt.Sprintf("%s:%d", addr.Host, addr.Port)
}

func (addr *NetAddress) Set(flagValue string) error {
	addrParts := strings.Split(flagValue, ":")

	if len(addrParts) == 2 {
		addr.Host = addrParts[0]
		optPort, err := strconv.Atoi(addrParts[1])
		if err != nil {
			return err
		}

		addr.Port = optPort
	}

	return nil
}

// RFCTime - дата и время в формате time.RFC3339
type RFCTime struct {
	time.Time
}

func (t RFCTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", t.Format(time.RFC3339))
	return []byte(stamp), nil
}

func (t *RFCTime) Scan(src interface{}) error {
	if value, ok := src.(time.Time); ok {
		t.Time = value
	}
	return nil
}

func (t RFCTime) Value() (driver.Value, error) {
	return t.Time, nil
}
