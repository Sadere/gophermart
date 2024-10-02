package config

import (
	"bytes"
	"flag"
	"log"
	"os"

	"github.com/Sadere/gophermart/internal/structs"
)

type Config struct {
	Address      structs.NetAddress // Адрес сервера
	PostgresDSN  string             // DSN строка для подключения к бд
	SecretKey    string             // Секретный ключ для подписи JWT токенов
	AccrualAddr  structs.NetAddress // Адрес сервиса accrual
	PullInterval int                // Интервал опроса accrual в секундах
}

const DefaultPullInterval = 10

func NewConfig(args []string) (Config, error) {
	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	var buf bytes.Buffer

	flags.SetOutput(&buf)

	newConfig := Config{
		Address: structs.NetAddress{
			Host: "",
			Port: 8080,
		},
	}

	flags.Var(&newConfig.Address, "a", "Адрес сервера")
	flags.StringVar(&newConfig.PostgresDSN, "d", "", "DSN для postgresql")
	flags.Var(&newConfig.AccrualAddr, "r", "Адрес сервиса accrual")
	flags.IntVar(&newConfig.PullInterval, "i", DefaultPullInterval, "Интервал опроса accrual в секундах")
	err := flags.Parse(args)
	if err != nil {
		return newConfig, err
	}

	// Конфиг из переменных окружений

	if envAddr := os.Getenv("RUN_ADDRESS"); len(envAddr) > 0 {
		err := newConfig.Address.Set(envAddr)
		if err != nil {
			log.Fatalf("Invalid server address supplied, RUN_ADDRESS = %s", envAddr)
		}
	}

	if envAccAddr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); len(envAccAddr) > 0 {
		err := newConfig.AccrualAddr.Set(envAccAddr)
		if err != nil {
			log.Fatalf("Invalid accrual address supplied, ACCRUAL_SYSTEM_ADDRESS = %s", envAccAddr)
		}
	}

	if envDSN := os.Getenv("DATABASE_URI"); len(envDSN) > 0 {
		newConfig.PostgresDSN = envDSN
	}

	envSecret, ok := os.LookupEnv("SECRET_KEY")
	if !ok || len(envSecret) == 0 {
		log.Fatal("no SECRET_KEY is set!")
	}

	newConfig.SecretKey = envSecret

	log.Println("server address: ", newConfig.Address)
	log.Println("accrual address: ", newConfig.AccrualAddr)

	return newConfig, nil
}
