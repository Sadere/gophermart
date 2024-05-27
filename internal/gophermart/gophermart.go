package gophermart

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sadere/gophermart/internal/config"
	"github.com/Sadere/gophermart/internal/database"
	"github.com/Sadere/gophermart/internal/repository"
	"github.com/Sadere/gophermart/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type GopherMart struct {
	config         config.Config
	userRepo       repository.UserRepository
	userService    *service.UserService
	orderService   *service.OrderService
	balanceService *service.BalanceService
	accService     *service.AccrualService
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

	// Подключаем сервисы
	g.InitServices(db)

	// Подключаем пути
	g.SetupRoutes(r, db)

	// Запускаем сервер
	srv := &http.Server{
		Addr:    g.config.Address.String(),
		Handler: r,
	}

	// Запускаем сервис опроса accrual
	go g.accService.Pull()

	// Запускаем сервер в фоне
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Ловим сигналы отключения сервера
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("graceful server shutdown ...")
}

func (g *GopherMart) InitServices(db *sqlx.DB) {
	userRepo := repository.NewPgUserRepository(db)
	g.userRepo = userRepo
	g.userService = service.NewUserService(userRepo)

	orderRepo := repository.NewPgOrderRepository(db)
	g.orderService = service.NewOrderService(orderRepo)

	balanceRepo := repository.NewPgBalanceRepository(db)
	g.balanceService = service.NewBalanceService(balanceRepo)

	g.accService = service.NewAccrualService(orderRepo,
		balanceRepo,
		g.config.AccrualAddr,
		time.Second*time.Duration(g.config.PullInterval),
	)
}

func Run() {
	app := &GopherMart{}
	app.config = config.NewConfig()

	app.Start()
}
