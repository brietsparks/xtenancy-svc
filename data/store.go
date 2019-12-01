package data

import (
	"database/sql"
	"errors"
	"github.com/gocraft/dbr/v2"
	"github.com/gocraft/dbr/v2/dialect"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gopkg.in/go-playground/validator.v9"
	"time"
)

type Store struct {
	db        *dbr.Session
	validator *validator.Validate
}

func NewStore(d *sql.DB, maxConn int) (*Store, error) {
	conn := &dbr.Connection{
		DB:            d,
		EventReceiver: &dbr.NullEventReceiver{},
		Dialect:       dialect.PostgreSQL,
	}

	conn.SetMaxOpenConns(maxConn)
	sess := conn.NewSession(nil)
	_, err := sess.Begin()

	if err != nil {
		return nil, errors.New("unable to create data store")
	}

	v := validator.New()

	return &Store{
		db:        sess,
		validator: v,
	}, nil
}

// CreateUser creates a new user
func (s *Store) CreateUser(u *User) (*User, error) {
	u.Id = uuid.New().String()

	if err := s.validate(u); err != nil {
		return nil, err
	}

	columns := []string{"id", "auth_id", "email", "first_name", "last_name",}
	err := s.create("user", u, columns)

	if err != nil {
		return nil, NewDbError(err)
	}

	return u, nil
}

// UpdateUser updates an existing user.
// The variadic "fields" arg should contain the field names that should be updated
func (s *Store) UpdateUser(id string, u *User, fields ...string) error {
	if err := s.validatePartial(u, fields...); err != nil {
		return err
	}

	err := s.update("user", id, fields,
		set{"AuthId", "auth_id", u.AuthId},
		set{"Email", "email", u.Email},
		set{"FirstName", "first_name", u.FirstName},
		set{"LastName", "last_name", u.LastName},
	)

	return NewDbError(err)
}

// GetUser gets a user by id
func (s *Store) GetUser(id string) (*User, error) {
	u := &User{}
	retrieved, count, err := s.getById("user", id, u)

	if err != nil {
		return nil, NewDbError(err)
	}

	if count == 0 {
		return nil, nil
	}

	return retrieved.(*User), nil
}

// DeleteUser deletes a user
func (s *Store) DeleteUser(id string) error {
	err := s.delete("user", id)
	return NewDbError(err)
}

// CreateTenant creates a new tenant
func (s *Store) CreateTenant(t *Tenant) (*Tenant, error) {
	t.Id = uuid.New().String()

	if err := s.validate(t); err != nil {
		return nil, err
	}

	columns := []string{"id", "name", "owner_id"}
	err := s.create("tenant", t, columns)

	if err != nil {
		return nil, NewDbError(err)
	}

	return t, nil
}

// UpdateTenant updates an existing tenant.
// The variadic "fields" arg should contain the field names that should be updated
func (s *Store) UpdateTenant(id string, t *Tenant, fields ...string) error {
	if err := s.validatePartial(t, fields...); err != nil {
		return err
	}

	err := s.update("tenant", id, fields,
		set{"Name", "name", t.Name},
		set{"OwnerId", "owner_id", t.OwnerId},
	)

	return NewDbError(err)
}

// GetTenant gets a tenant by id
func (s *Store) GetTenant(id string) (*Tenant, error) {
	t := &Tenant{}
	retrieved, count, err := s.getById("tenant", id, t)

	if err != nil {
		return nil, NewDbError(err)
	}

	if count == 0 {
		return nil, nil
	}

	return retrieved.(*Tenant), nil
}

// DeleteTenant deletes a tenant
func (s *Store) DeleteTenant(id string) error {
	err := s.delete("tenant", id)
	return NewDbError(err)
}

// CreateJoinrequest creates a new joinrequest
func (s *Store) CreateJoinrequest(jr *Joinrequest) (*Joinrequest, error) {
	jr.Id = uuid.New().String()
	jr.CreatedAt = time.Now()

	if err := s.validate(jr); err != nil {
		return nil, err
	}

	columns := []string{
		"id",
		"tenant_id",
		"user_id",
		"anon_email",
		"is_from_user",
		"created_at",
		"expires_at",
	}

	err := s.create("joinrequest", jr, columns)

	if err != nil {
		return nil, NewDbError(err)
	}

	return jr, nil
}

// UpdateJoinrequest updates an existing joinrequest.
// The variadic "fields" arg should contain the field names that should be updated
func (s *Store) UpdateJoinrequest(id string, jr *Joinrequest, fields ...string) error {
	if err := s.validatePartial(jr, fields...); err != nil {
		return err
	}

	err := s.update("joinrequest", id, fields,
		set{"TenantId", "tenant_id", jr.TenantId},
		set{"UserId", "user_id", jr.UserId},
		set{"AnonEmail", "anon_email", jr.AnonEmail},
		set{"IsAccepted", "is_accepted", jr.IsAccepted},
		set{"IsFromUser", "is_from_user", jr.IsFromUser},
		set{"ExpiresAt", "expires_at", jr.ExpiresAt},
	)

	return NewDbError(err)
}

// GetJoinrequest gets a joinrequest by id
func (s *Store) GetJoinrequest(id string) (*Joinrequest, error) {
	jr := &Joinrequest{}
	retrieved, count, err := s.getById("joinrequest", id, jr)

	if err != nil {
		return nil, NewDbError(err)
	}

	if count == 0 {
		return nil, nil
	}

	return retrieved.(*Joinrequest), nil
}

// DeleteJoinrequest deletes a joinrequest
func (s *Store) DeleteJoinrequest(id string) error {
	err := s.delete("joinrequest", id)
	return NewDbError(err)
}

// CreateMember creates a new member
func (s *Store) CreateMember(m *Member) (*Member, error) {
	m.Id = uuid.New().String()

	if err := s.validate(m); err != nil {
		return nil, err
	}

	columns := []string{
		"id",
		"tenant_id",
		"user_id",
		"alias",
		"is_admin",
		"is_inactive",
	}

	err := s.create("member", m, columns)

	if err != nil {
		return nil, NewDbError(err)
	}

	return m, nil
}

// UpdateMember updates an existing member.
// The variadic "fields" arg should contain the field names that should be updated
func (s *Store) UpdateMember(id string, m *Member, fields ...string) error {
	if err := s.validatePartial(m, fields...); err != nil {
		return err
	}

	err := s.update("member", id, fields,
		set{"TenantId", "tenant_id", m.TenantId},
		set{"UserId", "user_id", m.UserId},
		set{"Alias", "alias", m.Alias},
		set{"IsAdmin", "is_admin", m.IsAdmin},
		set{"IsInactive", "is_inactive", m.IsInactive},
	)

	return NewDbError(err)
}

// GetMember gets a member by id
func (s *Store) GetMember(id string) (*Member, error) {
	m := &Member{}
	retrieved, count, err := s.getById("member", id, m)

	if err != nil {
		return nil, NewDbError(err)
	}

	if count == 0 {
		return nil, nil
	}

	return retrieved.(*Member), nil
}

// DeleteMember deletes a member
func (s *Store) DeleteMember(id string) error {
	err := s.delete("member", id)
	return NewDbError(err)
}

func (s *Store) GetUsers(ids []string) ([]*User, error) {
	var u []*User

	stmt := s.db.SelectBySql(`select * from "user" where id = any(?)`, pq.Array(ids))

	_, err := stmt.Load(&u)

	if err != nil {
		return nil, NewDbError(err)
	}

	return u, nil
}

func (s *Store) GetUserByEmail(email string) (*User, error) {
	return nil, nil
}

func (s *Store) GetMemberByUserId(userId string) (*Member, error) {
	return nil, nil
}

func (s *Store) GetMembersByTenantId(tenantId string) ([]*Member, error) {
	return nil, nil
}

func (s *Store) GetMembersByUserId(userId string) ([]*Member, error) {
	return nil, nil
}

func (s *Store) GetJoinrequestsByUserId(userId string) ([]*Joinrequest, error) {
	return nil, nil
}

func (s *Store) GetJoinrequestsByTenantId(tenantId string) ([]*Joinrequest, error) {
	return nil, nil
}

func (s *Store) GetJoinrequestsByAnonEmail(email string) ([]*Joinrequest, error) {
	return nil, nil
}

func (s *Store) CheckTenantMember(tenantId string, memberId string) (bool, error) {
	return false, nil
}

func (s *Store) InviteByEmail(tenantId string, email string) (*Joinrequest, error) {
	u, err := s.GetUserByEmail(email)

	if err != nil {
		return nil, err
	}

	if u == nil {
		// send invite to anon
		j, err := s.CreateJoinrequest(&Joinrequest{
			TenantId:  tenantId,
			AnonEmail: dbr.NewNullString(email),
		})

		return j, err
	}

	m, err := s.GetMemberByUserId(u.Id)

	if err != nil {
		return nil, err
	}

	alreadyMember, err := s.CheckTenantMember(tenantId, m.Id)

	if alreadyMember == false {
		// send invite to user
		j, err := s.CreateJoinrequest(&Joinrequest{
			TenantId: tenantId,
			UserId:   dbr.NewNullString(u.Id),
		})

		return j, err
	}

	return nil, errors.New(ErrAlreadyMember)
}

func (s *Store) AcceptInvitation(joinrequestId string) error {
	_, err := s.GetJoinrequest(joinrequestId)

	if err != nil {
		return err
	}



    return nil
}

func (s *Store) PromoteMember(id string) error {
	return nil
}

func (s *Store) DemoteMember(id string) error {
	return nil
}

func (s *Store) ActivateMember(id string) error {
	return nil
}

func (s *Store) DeactivateMember(id string) error {
	return nil
}
