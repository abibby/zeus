package zeus

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"github.com/asdine/storm/v3"
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
	Where *Where `json:"where"`
}

type Where struct {
	Key      string      `json:"key"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type Response struct {
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

func (z *Zeus) Handle() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		queries := []*Query{}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			showError(w, err)
			return
		}
		err = json.Unmarshal(b, &queries)
		if err != nil {
			showError(w, err)
			return
		}

		for _, query := range queries {
			if model, ok := z.tables[query.Table]; ok {
				if query.Insert != nil {
					err = z.insert(w, model, query.Insert)
					if err != nil {
						showError(w, err)
						return
					}
				}
				if query.Select != nil {
					err = z.selectMoldes(w, model, query.Select)
					if err != nil {
						showError(w, err)
						return
					}
				}
			}
		}
	}
}

func showError(w http.ResponseWriter, err error) {
	log.Printf("%+v", err)
	err = json.NewEncoder(w).Encode(&Response{
		Error: err.Error(),
	})
	if err != nil {
		log.Printf("%+v", err)
	}
}

func showResponse(w http.ResponseWriter, data interface{}) error {
	err := json.NewEncoder(w).Encode(&Response{
		Data: data,
	})
	return errors.Wrap(err, "failed to encode response")
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
	return showResponse(w, s)
}

func (z *Zeus) selectMoldes(w http.ResponseWriter, model interface{}, selectQuery *SelectQuery) error {
	s := newSlice(model)
	if selectQuery.Where != nil {
		log.Print(selectQuery.Where.Key, selectQuery.Where.Value)
		err := z.db.Find(selectQuery.Where.Key, selectQuery.Where.Value, s)
		if err == storm.ErrNotFound {
			s = []struct{}{}
		} else if err != nil {
			return errors.Wrap(err, "unable to load selected models")
		}
	} else {
		err := z.db.All(s)
		if err != nil {
			return errors.Wrap(err, "unable to load selected models")
		}
	}
	return showResponse(w, s)
}
