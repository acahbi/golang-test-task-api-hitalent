package postgresql

import (
	"context"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewConnection(ctx context.Context, dsn string) (*gorm.DB, error) {
	/*
		CONFIG VALUE:
			Поле конфига        | Значение по умолчанию | Что делает
			Logger 				 logger.Default 		 Выводит в консоль только ошибки (Error)
			DryRun 				 false 					 Запросы реально выполняются в БД
			PrepareStmt 		 false 					 Не кэширует подготовленные выражения (PREPARE)
			DisableAutomaticPing false 				     Автоматически проверяет соединение перед запросами
			NamingStrategy 		 schema.NamingStrategy{} Таблицы в snake_case, имена во множественном числе (users)
			QueryFields 		 false 					 SELECT * вместо перечисления колонок
	*/
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: 100 * time.Millisecond,
				LogLevel:      logger.Info,
				Colorful:      false,
			}),
		PrepareStmt: true})

	if err != nil {
		return nil, err
	}
	db.WithContext(ctx)

	return db, nil
}
