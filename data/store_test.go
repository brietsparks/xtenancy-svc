package data

import (
	"database/sql"
	"errors"
	"flag"
	"github.com/gocraft/dbr/v2"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"gopkg.in/testfixtures.v2"
	"log"
	"testing"
)

type StoreTestSuite struct {
	suite.Suite
	Store    *Store
	fixtures *testfixtures.Context
}

var envPath string

func init() {
	flag.StringVar(&envPath, "env", "", "")
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}

func (s *StoreTestSuite) SetupSuite() {
	d := connect(s)

	// load fixtures
	fixtures, err := testfixtures.NewFolder(d, &testfixtures.PostgreSQL{}, "./fixtures")
	if err != nil {
		log.Fatal(err)
	}

	s.fixtures = fixtures

	// create store
	store, err := NewStore(d, 10)

	if err != nil {
		s.T().Fatalf("failed to create store: %s", err)
	}

	s.Store = store

	// clear tables
	err = clearTables(d)

	if err != nil {
		s.T().Fatalf("failed to clear table: %s", err)
	}
}

func (s *StoreTestSuite) TearDownSuite() {
	d := connect(s)

	err := clearTables(d)

	if err != nil {
		s.T().Fatalf("failed to clear table: %s", err)
	}
}

func connect(s *StoreTestSuite) *sql.DB {
	if envPath == "" {
		s.T().Fatal("missing variable --env <path to .env file>")
	}

	vars, err := LoadEnvVars(envPath)

	if err != nil {
		s.T().Fatalf("failed to load environment variables: %s", err)
	}

	url := MakeUrl(vars)
	d, err := sql.Open("postgres", url)

	if err != nil {
		s.T().Fatalf("failed to connect to database: %s", err)
	}

	return d
}

func clearTables(db *sql.DB) error {
	_, err := db.Query(`
		truncate table "user" cascade;
		truncate table tenant cascade;
		truncate table joinrequest cascade;
	`)

	return err
}

func (s *StoreTestSuite) SetupTest() {
	testfixtures.ResetSequencesTo(1)

	if err := s.fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}

func (s *StoreTestSuite) TestGetUser() {
	id := "00000000-0000-0000-0000-000000000000"
	u, _ := s.Store.GetUser(id)

	expected := &User{
		Id:        id,
		AuthId:    "00000000-0000-0000-0000-000000000000",
		Email:     "a@a.a",
		FirstName: "firstName0",
		LastName:  "lastName0",
	}
	s.Assert().EqualValues(expected, u)

	u, _ = s.Store.GetUser("00000000-0000-0000-7777-000000000000")
	s.Assert().Nil(u)
}
func (s *StoreTestSuite) TestUpdateUser() {
	id := "00000000-0000-0000-0000-000000000001"

	// doesn't update an unspecified field field
	_ = s.Store.UpdateUser(id, &User{FirstName: "abc",})
	u, _ := s.Store.GetUser(id)
	expected := &User{
		Id:        id,
		AuthId:    "00000000-0000-0000-0000-000000000001",
		Email:     "b@b.b",
		FirstName: "firstName1",
		LastName:  "lastName1",
	}
	s.Assert().EqualValues(expected, u)

	// updates a specified field
	_ = s.Store.UpdateUser(id, &User{FirstName: "abc",}, "FirstName")
	u, _ = s.Store.GetUser(id)
	expected = &User{
		Id:        id,
		AuthId:    "00000000-0000-0000-0000-000000000001",
		Email:     "b@b.b",
		FirstName: "abc",
		LastName:  "lastName1",
	}
	s.Assert().EqualValues(expected, u)

	//
	err := s.Store.UpdateUser("00000000-0000-0000-7777-000000000001", &User{FirstName: "abc"}, "FirstName")
	s.Assert().Equal(ErrResourceDNE, err.Error())
	s.Assert().Nil(errors.Unwrap(err))
}

func (s *StoreTestSuite) TestCreateUser() {
	created, _ := s.Store.CreateUser(&User{
		AuthId:    "00000000-0000-0000-0000-000000000002",
		Email:     "z@z.z",
		FirstName: "foo",
		LastName:  "bar",
	})

	retrieved, _ := s.Store.GetUser(created.Id)
	s.Assert().EqualValues(created, retrieved)
}

func (s *StoreTestSuite) TestDeleteUser() {
	id := "00000000-0000-0000-0000-000000000003"

	retrieved, _ := s.Store.GetUser(id)
	s.Assert().NotNil(retrieved)

	_ = s.Store.DeleteUser(id)
	retrieved, _ = s.Store.GetUser(id)
	s.Assert().Nil(retrieved)

	err := s.Store.DeleteUser("00000000-0000-0000-7777-000000000003")
	s.Assert().Equal(ErrResourceDNE, err.Error())
	s.Assert().Nil(errors.Unwrap(err))
}

func (s *StoreTestSuite) TestGetTenant() {
	id := "00000000-0000-0000-0000-000000000000"
	u, _ := s.Store.GetTenant(id)

	expected := &Tenant{
		Id:   id,
		Name: "name0",
	}
	s.Assert().EqualValues(expected, u)

	u, _ = s.Store.GetTenant("00000000-0000-0000-7777-000000000000")
	s.Assert().Nil(u)
}

func (s *StoreTestSuite) TestUpdateTenant() {
	id := "00000000-0000-0000-0000-000000000001"

	_ = s.Store.UpdateTenant(id, &Tenant{Name: "abc",})

	u, _ := s.Store.GetTenant(id)
	expected := &Tenant{
		Id:   id,
		Name: "name1",
	}
	s.Assert().EqualValues(expected, u)

	_ = s.Store.UpdateTenant(id, &Tenant{Name: "abc",}, "Name")
	u, _ = s.Store.GetTenant(id)
	expected = &Tenant{
		Id:   id,
		Name: "abc",
	}
	s.Assert().EqualValues(expected, u)

	err := s.Store.UpdateTenant("00000000-0000-0000-7777-000000000001", &Tenant{Name: "abc"}, "Name")
	s.Assert().Equal(ErrResourceDNE, err.Error())
	s.Assert().Nil(errors.Unwrap(err))
}

func (s *StoreTestSuite) TestCreateTenant() {
	created, _ := s.Store.CreateTenant(&Tenant{
		Name: "foo",
	})

	retrieved, _ := s.Store.GetTenant(created.Id)

	s.Assert().EqualValues(created, retrieved)
}

func (s *StoreTestSuite) TestDeleteTenant() {
	id := "00000000-0000-0000-0000-000000000003"

	retrieved, _ := s.Store.GetTenant(id)
	s.Assert().NotNil(retrieved)

	_ = s.Store.DeleteTenant(id)
	retrieved, _ = s.Store.GetTenant(id)
	s.Assert().Nil(retrieved)

	err := s.Store.DeleteTenant("00000000-0000-0000-7777-000000000003")
	s.Assert().Equal(ErrResourceDNE, err.Error())
	s.Assert().Nil(errors.Unwrap(err))
}

func (s *StoreTestSuite) TestGetJoinrequest() {
	id := "00000000-0000-0000-0000-000000000000"
	jr, _ := s.Store.GetJoinrequest(id)

	expected := &Joinrequest{
		Id:         id,
		TenantId:   "00000000-0000-0000-0000-000000000000",
		UserId:     dbr.NewNullString(nil),
		AnonEmail:  dbr.NewNullString(nil),
		IsAccepted: dbr.NewNullBool(nil),
		IsFromUser: dbr.NewNullBool(nil),
		ExpiresAt:  dbr.NewNullTime(nil),
	}

	s.Assert().Equal(expected.comparable(), jr.comparable())

	jr, _ = s.Store.GetJoinrequest("00000000-0000-0000-7777-000000000000")
	s.Assert().Nil(jr)
}

func (s *StoreTestSuite) TestUpdateJoinrequest() {
	id := "00000000-0000-0000-0000-000000000001"

	// doesn't update without specifying field
	_ = s.Store.UpdateJoinrequest(id, &Joinrequest{IsAccepted: dbr.NewNullBool(true)})
	u, _ := s.Store.GetJoinrequest(id)
	expected := &Joinrequest{
		Id:         id,
		TenantId:   "00000000-0000-0000-0000-000000000001",
		IsAccepted: dbr.NewNullBool(nil),
	}
	s.Assert().EqualValues(expected.comparable(), u.comparable())

	// updates specified field
	_ = s.Store.UpdateJoinrequest(id, &Joinrequest{IsAccepted: dbr.NewNullBool(true)}, "IsAccepted")
	u, _ = s.Store.GetJoinrequest(id)
	expected = &Joinrequest{
		Id:         id,
		TenantId:   "00000000-0000-0000-0000-000000000001",
		IsAccepted: dbr.NewNullBool(true),
	}
	s.Assert().EqualValues(expected.comparable(), u.comparable())

	// error on non-existent resource
	err := s.Store.UpdateJoinrequest(
		"00000000-0000-0000-7777-000000000001",
		&Joinrequest{IsAccepted: dbr.NewNullBool(true)},
		"IsAccepted",
	)
	s.Assert().Equal(ErrResourceDNE, err.Error())
	s.Assert().Nil(errors.Unwrap(err))
}

func (s *StoreTestSuite) TestCreateJoinrequest() {
	created, _ := s.Store.CreateJoinrequest(&Joinrequest{
		TenantId: "00000000-0000-0000-0000-000000000000",
	})

	retrieved, _ := s.Store.GetJoinrequest(created.Id)

	s.Assert().EqualValues(created.comparable(), retrieved.comparable())
}

func (s *StoreTestSuite) TestDeleteJoinrequest() {
	id := "00000000-0000-0000-0000-000000000003"
	retrieved, _ := s.Store.GetJoinrequest(id)
	s.Assert().NotNil(retrieved)

	_ = s.Store.DeleteJoinrequest(id)
	retrieved, _ = s.Store.GetJoinrequest(id)
	s.Assert().Nil(retrieved)

	err := s.Store.DeleteJoinrequest("00000000-0000-0000-7777-000000000003")
	s.Assert().Equal(ErrResourceDNE, err.Error())
	s.Assert().Nil(errors.Unwrap(err))
}
