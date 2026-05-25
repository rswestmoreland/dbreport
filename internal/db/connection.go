package db

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/rswestmoreland/dbreport/internal/config"
)

const driverName = "mysql"

func Open(ctx context.Context, cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn, timeout, err := BuildDSNFromEnv(cfg)
	if err != nil {
		return nil, err
	}

	handle, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("open database handle: %w", err)
	}

	handle.SetMaxOpenConns(1)
	handle.SetMaxIdleConns(1)
	handle.SetConnMaxLifetime(3 * time.Minute)
	handle.SetConnMaxIdleTime(1 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := handle.PingContext(pingCtx); err != nil {
		_ = handle.Close()
		return nil, fmt.Errorf("connect to database %s: %w", SafeTarget(cfg), err)
	}

	return handle, nil
}

func BuildDSNFromEnv(cfg config.DatabaseConfig) (string, time.Duration, error) {
	user, password, err := CredentialsFromEnv(cfg)
	if err != nil {
		return "", 0, err
	}
	return BuildDSN(cfg, user, password)
}

func BuildDSN(cfg config.DatabaseConfig, user string, password string) (string, time.Duration, error) {
	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		return "", 0, fmt.Errorf("database.timeout_seconds must be greater than zero")
	}

	driverCfg := mysql.NewConfig()
	driverCfg.User = user
	driverCfg.Passwd = password
	driverCfg.Net = "tcp"
	driverCfg.Addr = net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	driverCfg.DBName = cfg.Name
	driverCfg.ParseTime = true
	driverCfg.Timeout = timeout
	driverCfg.ReadTimeout = timeout
	driverCfg.WriteTimeout = timeout
	driverCfg.Params = map[string]string{
		"charset": "utf8mb4",
	}
	if cfg.TLS {
		driverCfg.TLSConfig = "true"
	}

	return driverCfg.FormatDSN(), timeout, nil
}

func CredentialsFromEnv(cfg config.DatabaseConfig) (string, string, error) {
	userEnv := strings.TrimSpace(cfg.UserEnv)
	passwordEnv := strings.TrimSpace(cfg.PasswordEnv)
	if userEnv == "" {
		return "", "", fmt.Errorf("database.user_env is required")
	}
	if passwordEnv == "" {
		return "", "", fmt.Errorf("database.password_env is required")
	}

	user := os.Getenv(userEnv)
	if user == "" {
		return "", "", fmt.Errorf("database username environment variable is not set or is empty: %s", userEnv)
	}
	password := os.Getenv(passwordEnv)
	if password == "" {
		return "", "", fmt.Errorf("database password environment variable is not set or is empty: %s", passwordEnv)
	}

	return user, password, nil
}

func SafeTarget(cfg config.DatabaseConfig) string {
	return fmt.Sprintf("%s/%s", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)), cfg.Name)
}
