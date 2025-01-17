package db

import (
	"slices"
	"strings"

	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog/log"
	"github.com/thomiceli/opengist/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Setup(dbPath string, sharedCache bool) error {
	var err error
	journalMode := strings.ToUpper(config.C.SqliteJournalMode)

	if !slices.Contains([]string{"DELETE", "TRUNCATE", "PERSIST", "MEMORY", "WAL", "OFF"}, journalMode) {
		log.Warn().Msg("Invalid SQLite journal mode: " + journalMode)
	}

	sharedCacheStr := ""
	if sharedCache {
		sharedCacheStr = "&cache=shared"
	}

	if config.C.DBUrl != "" {
		if db, err = gorm.Open(postgres.Open(config.C.DBUrl), &gorm.Config{
			TranslateError: true,
			Logger:         logger.Default.LogMode(logger.Silent),
		}); err != nil {
			return err
		}
	} else {
		if db, err = gorm.Open(sqlite.Open(dbPath+"?_fk=true&_journal_mode="+journalMode+sharedCacheStr), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}); err != nil {
			return err
		}
	}

	if err = db.SetupJoinTable(&Gist{}, "Likes", &Like{}); err != nil {
		return err
	}

	if err = db.SetupJoinTable(&User{}, "Liked", &Like{}); err != nil {
		return err
	}

	if err = db.AutoMigrate(&User{}, &Gist{}, &SSHKey{}, &AdminSetting{}); err != nil {
		return err
	}

	if err = ApplyMigrations(db); err != nil {
		return err
	}

	// Default admin setting values
	return initAdminSettings(map[string]string{
		SettingDisableSignup:    "0",
		SettingRequireLogin:     "0",
		SettingDisableLoginForm: "0",
		SettingDisableGravatar:  "0",
	})
}

func Close() error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func CountAll(table interface{}) (int64, error) {
	var count int64
	err := db.Model(table).Count(&count).Error
	return count, err
}

func IsUniqueConstraintViolation(err error) bool {
	return err == gorm.ErrDuplicatedKey
}

func Ping() error {
	sql, err := db.DB()
	if err != nil {
		return err
	}

	return sql.Ping()
}
