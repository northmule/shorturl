package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/db"
	"github.com/northmule/shorturl/internal/app/handlers"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/url"
	appStorage "github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/workers"
	"github.com/northmule/shorturl/internal/grpc/contract"
	grpcHandlers "github.com/northmule/shorturl/internal/grpc/handlers"
	"github.com/northmule/shorturl/internal/grpc/handlers/interceptors"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	printLabel()
	appCtx, appStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer appStop()
	if err := run(appCtx); err != nil {
		log.Fatal(err)
	}
}

// run преднастройка
func run(ctx context.Context) error {
	var err error
	err = logger.InitLogger("info")
	if err != nil {
		return err
	}
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	storage, finderStats, err := getStorage(ctx, cfg)
	if err != nil {
		return err
	}
	sessionStorage := appStorage.NewSessionStorage()
	shortURLService := url.NewShortURLService(storage, storage)
	stop := make(chan struct{})
	worker := workers.NewWorker(storage, stop)

	lc := net.ListenConfig{}
	listen, err := lc.Listen(ctx, "tcp", cfg.ServerURL)
	if err != nil {
		return err
	}

	logger.LogSugar.Info("создаём gRPC-сервер")
	authInterceptor := interceptors.NewCheckAuth(storage, sessionStorage)
	trustedInterceptor := interceptors.NewCheckTrustedSubnet(cfg)
	loggerInterceptor := interceptors.NewLogger(logger.LogSugar)

	s := grpc.NewServer(grpc.ChainUnaryInterceptor([]grpc.UnaryServerInterceptor{
		loggerInterceptor.LogStart,
		authInterceptor.AuthEveryone,
		authInterceptor.AccessVerificationUserUrls,
		trustedInterceptor.GrantAccess,
		loggerInterceptor.LogEnd,
	}...))

	logger.LogSugar.Info("Подготовка сервисов")
	contract.RegisterPingHandlerServer(s, grpcHandlers.NewPingHandler(storage))
	contract.RegisterRedirectHandlerServer(s, grpcHandlers.NewRedirectHandler(shortURLService))
	contract.RegisterShortenerHandlerServer(s, grpcHandlers.NewShortenerHandler(shortURLService, storage, storage))
	contract.RegisterStatsHandlerServer(s, grpcHandlers.NewStatsHandler(finderStats))
	contract.RegisterUserUrlsHandlerServer(s, grpcHandlers.NewUserURLsHandler(storage, sessionStorage, worker))

	logger.LogSugar.Infof("Running server on - %s", cfg.ServerURL)
	go func() {
		<-ctx.Done()
		stop <- struct{}{}
		logger.LogSugar.Info("Получин сигнал. Останавливаю сервер...")
		s.GracefulStop()
	}()
	if err = s.Serve(listen); err != nil {
		return err
	}

	return nil
}

func getStorage(ctx context.Context, cfg *config.Config) (appStorage.StorageQuery, handlers.FinderStats, error) {

	if cfg.DataBaseDsn != "" {
		s, err := appStorage.NewPostgresStorage(cfg.DataBaseDsn)
		if err != nil {
			logger.LogSugar.Errorf("Failed NewPostgresStorage dsn: %s, %s", cfg.DataBaseDsn, err)
			return nil, nil, err
		}

		logger.LogSugar.Info("Инициализация миграций")
		migrations := db.NewMigrations(s.RawDB)
		err = migrations.Up(ctx)
		if err != nil {
			return nil, nil, err
		}

		return s, s, nil
	}

	if cfg.FileStoragePath != "" {
		file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			logger.LogSugar.Errorf("Failed to open file %s: error: %s", cfg.FileStoragePath, err)
			return nil, nil, err
		}
		s := appStorage.NewFileStorage(file)
		return s, s, nil
	}

	s := appStorage.NewMemoryStorage()
	return s, s, nil
}

func printLabel() {
	template := `
	Build version: <buildVersion>
	Build date: <buildDate>
	Build commit: <buildCommit>
`
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}

	if buildCommit == "" {
		buildCommit = "N/A"
	}
	template = strings.Replace(template, "<buildVersion>", buildVersion, 1)
	template = strings.Replace(template, "<buildDate>", buildDate, 1)
	template = strings.Replace(template, "<buildCommit>", buildCommit, 1)

	fmt.Println(template)
}
