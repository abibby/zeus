package zeus

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"github.com/pkg/errors"
)

type Query struct {
	Table  string       `json:"table"`
	Insert *InsertQuery `json:"insert"`
	Select *SelectQuery `json:"select"`
}

type InsertQuery struct {
	Value json.RawMessage `json:"value"`
}
type SelectQuery struct {
}

func (z *Zeus) Handle() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		queries := []*Query{}

		err := json.NewDecoder(r.Body).Decode(&queries)
		if err != nil {
			log.Print(err)
			return
		}

		for _, query := range queries {
			if model, ok := z.tables[query.Table]; ok {
				if query.Insert != nil {
					err = z.insert(w, model, query.Insert)
					if err != nil {
						log.Print(err)
						return
					}
				}
				if query.Select != nil {
					err = z.selectMoldes(w, model, query.Select)
					if err != nil {
						log.Print(err)
						return
					}
				}
			}
		}
	}
}

// newSlice returns a pointer to a slice of the given type
func newSlice(model interface{}) interface{} {
	v := reflect.ValueOf(model)
	return reflect.New(reflect.MakeSlice(reflect.SliceOf(v.Type()), 0, 0).Type()).Interface()
}

func newStruct(model interface{}) interface{} {
	v := reflect.ValueOf(model)
	return reflect.New(v.Elem().Type()).Interface()
}

func (z *Zeus) insert(w http.ResponseWriter, model interface{}, insert *InsertQuery) error {
	s := newStruct(model)
	err := json.Unmarshal([]byte(insert.Value), s)
	if err != nil {
		return errors.Wrap(err, "failed to load model to insert")
	}
	err = z.db.Save(s)
	if err != nil {
		return errors.Wrap(err, "could not save the given model")
	}
	err = json.NewEncoder(w).Encode(s)
	return errors.Wrap(err, "failed to print inserted model")
}

func (z *Zeus) selectMoldes(w http.ResponseWriter, model interface{}, selectQuery *SelectQuery) error {
	s := newSlice(model)
	err := z.db.All(s)
	if err != nil {
		return errors.Wrap(err, "unable to load selected models")
	}
	err = json.NewEncoder(w).Encode(s)
	return errors.Wrap(err, "unable to write selected models")
}
