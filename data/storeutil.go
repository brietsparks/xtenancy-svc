package data

import (
	"errors"
	"fmt"
	"github.com/gocraft/dbr/v2"
	"github.com/gocraft/dbr/v2/dialect"
	"github.com/lib/pq"
)

type junction struct {
	table1        string
	table2        string
	junctionTable string
	table1Pk      string
	table2Pk      string
	junctionFk1   string
	junctionFk2   string
}

func (s *Store) validatePartial(resource interface{}, fields ...string) error {
	if len(fields) == 0 {
		return errors.New(ErrEmptyFieldMask)
	}

    return s.validator.StructPartial(resource, fields...)
}

func (s *Store) validate(resource interface{}) error {
    return s.validator.Struct(resource)
}

func (s *Store) selectJunction(db *dbr.Session, lookupId interface{}, j junction) *dbr.SelectStmt {
	if j.table1Pk == "" {
		j.table1Pk = "id"
	}

	if j.table2Pk == "" {
		j.table2Pk = "id"
	}

	// wrapping table names in quotes prevents errors when tables/columns are named after reserved words
	return db.
		Select(quotes(j.table1)+".*").
		From(quotes(j.table1)).
		Join(
			j.junctionTable,
			fmt.Sprintf("%s.%s = %s.%s", quotes(j.table1), j.table1Pk, j.junctionTable, j.junctionFk1),
		).
		Join(
			j.table2,
			fmt.Sprintf("%s.%s = %s.%s", quotes(j.table2), j.table2Pk, j.junctionTable, j.junctionFk2),
		).
		Where(fmt.Sprintf("%s.%s = ?", quotes(j.table2), j.table2Pk), lookupId)
}

func quotes(s string) string {
	return fmt.Sprintf("\"%s\"", s)
}

func (s *Store) create(table string, record interface{}, columns []string) error {
	_, err := s.db.
		InsertInto(table).
		Columns(columns...).
		Record(record).
		Exec()

	return err
}

func (s *Store) update(table string, id interface{}, fields []string, updateSets ...set) error {
	setMap := makeSetMap(fields, updateSets...)

	if len(setMap) == 0 {
		return errors.New("update setMap contains zero fields")
	}

	result, err := s.db.
		Update(table).
		SetMap(setMap).
		Where("id = ?", id).
		Exec()

	if err != nil {
		return err
	}

	count, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New(ErrResourceDNE)
	}

	return err
}

func (s *Store) getById(table string, id interface{}, resource interface{}) (interface{}, int, error) {
	count, err := s.db.
		Select("*").
		From(quotes(table)).
		Where("id = ?", id).
		Load(resource)

	if err != nil {
		return nil, 0, err
	}

	return resource, count, nil
}

func (s *Store) getManyByIds(table string, ids interface{}, resources interface{}) (interface{}, error) {
	_, err := s.db.
		Select("*").
		From(quotes(table)).
		Where("id = any(?)", pq.Array(ids)).
		Load(&resources)

	if err != nil {
		return nil, err
	}

	return resources, nil
}

func (s *Store) delete(table string, id interface{}) error {
	result, err := s.db.DeleteFrom(table).Where("id = ?", id).Exec()

	if err != nil {
		return err
	}

	count, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New(ErrResourceDNE)
	}

	return err
}

func (s *Store) unlink(junctionTable string, pk1 string, id1 interface{}, pk2 string, id2 interface{}) error {
	result, err := s.db.
		InsertInto(junctionTable).
		Pair(pk1, id1).
		Pair(pk2, id2).
		Exec()

	if err != nil {
		return err
	}

	count, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	return nil
}

func includes(strings []string, val string) bool {
	for _, v := range strings {
		if v == val {
			return true
		}
	}

	return false
}

type set struct {
	Field string
	Col   string
	Val   interface{}
}

func makeSetMap(fields []string, sets ...set) map[string]interface{} {
	setMap := map[string]interface{}{}

	for _, set := range sets {
		if includes(fields, set.Field) {
			setMap[set.Col] = set.Val
		}
	}

	return setMap
}

type stmt interface {
	Build(d dbr.Dialect, buf dbr.Buffer) error
}

func dumpStmt(stmt stmt) string {
	buf := dbr.NewBuffer()
	_ = stmt.Build(dialect.PostgreSQL, buf)
	return buf.String()
}
