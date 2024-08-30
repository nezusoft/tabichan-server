package app

import (
	"github.com/tabichanorg/tabichan-server/internal/config"
	"github.com/tabichanorg/tabichan-server/internal/server"
)

func InitializeApp() (*server.Server, error) {
	cfg := config.LoadConfig()

	// postgresDB, err := db.NewPostgresDB(cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort)
	// if err != nil {
	// 	return nil, err
	// }
	// defer postgresDB.Close()

	// userRepo := user.NewRepository(postgresDB.Conn)
	// userService := user.NewService(userRepo)

	srv := server.NewServer(cfg.HTTPPort)

	return srv, nil
}
