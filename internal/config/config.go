package config

import (
	"flag"
	"log"
	"os"

	"github.com/Sadere/gophermart/internal/structs"
)

type Config struct {
	Address     structs.NetAddress // Адрес сервера
	PostgresDSN string             // DSN строка для подключения к бд
	SecretKey   string             // Секретный ключ для подписи JWT токенов
}

func NewConfig() Config {
	newConfig := Config{
		Address: structs.NetAddress{
			Host: "",
			Port: 8080,
		},
	}

	flag.Var(&newConfig.Address, "a", "Адрес сервера")
	flag.StringVar(&newConfig.PostgresDSN, "d", "", "DSN для postgresql")
	flag.Parse()

	// Конфиг из переменных окружений

	if envAddr := os.Getenv("RUN_ADDRESS"); len(envAddr) > 0 {
		err := newConfig.Address.Set(envAddr)
		if err != nil {
			log.Fatalf("Invalid server address supplied, RUN_ADDRESS = %s", envAddr)
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

	return newConfig
}
