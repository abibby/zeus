package zeus

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type Zeus struct {
	db     *gorm.DB
	tables map[string]interface{}
}

func Open(tables map[string]interface{}) (*Zeus, error) {
	db, err := gorm.Open("sqlite3", "./foo.db")
	if err != nil {
		return nil, errors.Wrap(err, "error opening db")
	}

	for _, table := range tables {
		db.AutoMigrate(table)
	}

	return &Zeus{
		db:     db,
		tables: tables,
	}, nil
}

func (z *Zeus) Close() error {
	return errors.Wrap(z.db.Close(), "error closing db")
}
