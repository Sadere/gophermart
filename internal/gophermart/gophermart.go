package gophermart

import (
	"log"

	"github.com/Sadere/gophermart/internal/config"
	"github.com/Sadere/gophermart/internal/database"
	"github.com/gin-gonic/gin"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type GopherMart struct {
	config config.Config
}

func (g *GopherMart) Start() {
	r := gin.Default()

	// Миграции
	if err := database.MigrateUp(g.config.PostgresDSN); err != nil {
		log.Fatal("failed to run migrations: ", err)
	}

	// Подключаемся к БД
	db, err := database.NewConnection("pgx", g.config.PostgresDSN)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	// Подключаем пути
	g.SetupRoutes(r, db)

	// Запускаем сервер
	err = r.Run(g.config.Address.String())
	if err != nil {
		log.Fatal("failed to run server: ", err)
	}
}

func Run() {
	app := &GopherMart{}
	app.config = config.NewConfig()

	app.Start()
}
