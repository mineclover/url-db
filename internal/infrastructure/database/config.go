package database

import "time"

type Config struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	WALMode         bool
	ForeignKeys     bool
	JournalMode     string
	Synchronous     string
}

func DefaultConfig() *Config {
	return &Config{
		URL:             "file:./app.db",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
		WALMode:         true,
		ForeignKeys:     true,
		JournalMode:     "WAL",
		Synchronous:     "NORMAL",
	}
}

func TestConfig() *Config {
	return &Config{
		URL:             ":memory:",
		MaxOpenConns:    1,
		MaxIdleConns:    1,
		ConnMaxLifetime: time.Hour,
		WALMode:         false,
		ForeignKeys:     true,
		JournalMode:     "DELETE",
		Synchronous:     "OFF",
	}
}

func ProductionConfig(dbPath string) *Config {
	return &Config{
		URL:             "file:" + dbPath,
		MaxOpenConns:    100,
		MaxIdleConns:    50,
		ConnMaxLifetime: time.Hour,
		WALMode:         true,
		ForeignKeys:     true,
		JournalMode:     "WAL",
		Synchronous:     "FULL",
	}
}
