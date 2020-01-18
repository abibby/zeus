package zeus

import (
	"github.com/asdine/storm/v3"
	"github.com/pkg/errors"
)

type Zeus struct {
	db     *storm.DB
	tables map[string]interface{}
}

func Open(tables map[string]interface{}) (*Zeus, error) {
	db, err := storm.Open("./test.db")
	if err != nil {
		return nil, errors.Wrap(err, "error opening db")
	}
	return &Zeus{
		db:     db,
		tables: tables,
	}, nil
}

func (z *Zeus) Close() error {
	return errors.Wrap(z.db.Close(), "error closing db")
}
