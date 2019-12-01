package data

import (
	"github.com/gocraft/dbr/v2"
	"time"
)

type User struct {
	Id        string `db:"id" json:"id" validate:"uuid,required"`
	AuthId    string `db:"auth_id" json:"authId" validate:"uuid,required"`
	Email     string `db:"email" validate:"email,required"`
	FirstName string `db:"first_name" json:"firstName,required"`
	LastName  string `db:"last_name" json:"lastName,required"`
}

type Tenant struct {
	Id   string `db:"id" json:"id" validate:"uuid,required"`
	Name string `db:"name" json:"required"`
	OwnerId string `db:"owner_id" json:"ownerId" validate:"uuid,required"`
}

type Joinrequest struct {
	Id         string         `db:"id" json:"id" validate:"uuid,required"`
	TenantId   string         `db:"tenant_id" json:"id" validate:"uuid,required"`
	UserId     dbr.NullString `db:"user_id" json:"userId" validate:"uuid"`
	AnonEmail  dbr.NullString `db:"anon_email" json:"anonEmail" validate:"email"`
	IsAccepted dbr.NullBool   `db:"is_accepted" json:"isAccepted"`
	IsFromUser dbr.NullBool   `db:"is_from_user" json:"isFromUser"`
	CreatedAt  time.Time      `db:"created_at" json:"createdAt,required"`
	ExpiresAt  dbr.NullTime   `db:"expires_at" json:"expiresAt"`
}

type Member struct {
	Id         string         `db:"id" json:"id" validate:"uuid,required"`
	TenantId   string         `db:"tenant_id" json:"id" validate:"uuid,required"`
	UserId     string         `db:"user_id" json:"userId" validate:"uuid,required"`
	Alias      dbr.NullString `db:"alias" json:"alias"`
	IsAdmin    bool           `db:"is_admin" json:"isAdmin,required"`
	IsInactive bool           `db:"is_inactive" json:"isInactive"`
}

func (j *Joinrequest) comparable() *Joinrequest {
	return &Joinrequest{
		Id:         j.Id,
		TenantId:   j.TenantId,
		AnonEmail:  j.AnonEmail,
		IsAccepted: j.IsAccepted,
		IsFromUser: j.IsFromUser,
		ExpiresAt:  j.ExpiresAt,
	}
}
