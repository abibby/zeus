package zeus

import (
	"io"
	"net/http"

	"github.com/asdine/storm"
	"github.com/pkg/errors"
)

type Zeus struct {
	db        *storm.DB
	changes   chan *Change
	listeners map[int64]chan *Change
}

type Change struct{}

func Open(...interface{}) (*Zeus, error) {
	db, err := storm.Open("./test.db")
	if err != nil {
		return nil, errors.Wrap(err, "error opening db")
	}
	return &Zeus{
		db: db,
	}, nil
}

func (z *Zeus) Close() error {
	return errors.Wrap(z.db.Close(), "error closing db")
}

func (z *Zeus) Handle() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		io.Copy(w, r.Body)
		// fmt.Fprintf(w, "Hello, World!")
	}
}

func (z *Zeus) OnChange() chan *Change {
	cc := make(chan *Change)
	z.listeners[0] = cc
	return cc
}

func (z *Zeus) change(cc chan *Change) {

}
