package datastore

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"strings"

	pbM "github.com/Investliftng/ocm-api/merchant/proto/merchant"
	"github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes"
	"github.com/micro/go-micro/v2/errors"
	"golang.org/x/crypto/bcrypt"
)

//MerchantRepo communicates with database and squirrel
type MerchantRepo interface {
	Create(context.Context, *pbM.CreateRequest) error
	Read(context.Context, *pbM.GetIdRequest) (*pbM.Merchant, error)
	Update(context.Context, *pbM.UpdateRequest) error
	UpdataPassword(context.Context, *pbM.UpdatePasswordRequest) error
	Get(context.Context) ([]*pbM.Merchant, error)
	Delete(context.Context, *pbM.GetIdRequest) error
}

var (
	alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

func random(i int) string {
	bytes := make([]byte, i)
	for {
		rand.Read(bytes)
		for i, b := range bytes {
			bytes[i] = alphanum[b%byte(len(alphanum))]
		}
		return string(bytes)
	}
}

//Create a merchant in the db
func (c *Connection) Create(ctx context.Context, req *pbM.CreateRequest) error {
	salt := random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(salt+req.Password), 10)
	if err != nil {
		c.Logger.Errorf("failed to create hash password %s\n", err.Error())
		return err
	}
	req.Password = base32.StdEncoding.EncodeToString(h)
	_, err = c.SQLBuilder.Insert(
		"merchants",
	).SetMap(map[string]interface{}{
		"number_of_product": req.Merchant.GetNumberOfProduct(),
		"email":             strings.ToLower(req.Merchant.GetEmail()),
		"phone":             req.Merchant.Phone.GetNumber(),
		"phone_type":        (phoneTypeWrapper)(req.Merchant.Phone.GetType()),
		"user_id":           req.Merchant.GetUserId(),
		"number_of_outlet":  req.Merchant.GetNumberOfOutlet(),
		"business_name":     req.Merchant.GetBusinessName(),
		"role_type":         (roleTypeWrapper)(req.Merchant.GetRole()),
		"salt":              salt,
		"passwd":            req.Password,
	}).ExecContext(ctx)
	return err
}

//Update a merchant
func (c *Connection) Update(ctx context.Context, req *pbM.UpdateRequest) error {
	now := ptypes.TimestampNow()
	_, err := c.SQLBuilder.Update(
		"merchants",
	).SetMap(map[string]interface{}{
		"number_of_product": req.Merchant.GetNumberOfProduct(),
		"email":             strings.ToLower(req.Merchant.GetEmail()),
		"phone":             req.Merchant.Phone.GetNumber(),
		"phone_type":        (phoneTypeWrapper)(req.Merchant.Phone.GetType()),
		"number_of_outlet":  req.Merchant.GetNumberOfOutlet(),
		"business_name":     req.Merchant.GetBusinessName(),
		"role_type":         (roleTypeWrapper)(req.Merchant.GetRole()),
		"updated_at":        (*timeWrapper)(now),
	}).Where(squirrel.Eq{"id": req.GetId()}).ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

//UpdataPassword of a merchant
func (c *Connection) UpdataPassword(ctx context.Context, req *pbM.UpdatePasswordRequest) error {
	return nil
}

//Get a list merchants
func (c *Connection) Get(ctx context.Context) (ms []*pbM.Merchant, err error) {
	rows, err := c.SQLBuilder.Select(
		"id",
		"number_of_product",
		"role_type",
		"email",
		"phone",
		"phone_type",
		"user_id",
		"number_of_outlet",
		"business_name",
		"is_suspended",
		"is_email_verified",
		"last_login",
		"created_at",
		"updated_at",
	).From(
		"merchants",
	).QueryContext(ctx)
	if err != nil {
		return
	}
	defer func() {
		cerr := rows.Close()
		if err == nil && cerr != nil {
			err = errors.InternalServerError("501", "service.Merchant.Get", err)
		}
	}()

	for rows.Next() {
		merchant, err := scanMerchant(rows)
		if err != nil {
			return nil, err
		}
		ms = append(ms, merchant)
	}
	return
}

//Delete a merchant
func (c *Connection) Delete(ctx context.Context, req *pbM.GetIdRequest) error {
	_, err := c.SQLBuilder.Delete(
		"merchants",
	).Where(squirrel.Eq{"id": req.GetId()}).ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

//Read return a merchant
func (c *Connection) Read(ctx context.Context, req *pbM.GetIdRequest) (*pbM.Merchant, error) {
	q, args, err := c.SQLBuilder.Select(
		"id",
		"number_of_product",
		"role_type",
		"email",
		"phone",
		"phone_type",
		"user_id",
		"number_of_outlet",
		"business_name",
		"is_suspended",
		"is_email_verified",
		"last_login",
		"created_at",
		"updated_at",
	).From(
		"merchants",
	).Where(squirrel.Eq{"id": req.GetId()}).ToSql()
	if err != nil {
		return nil, err
	}
	return scanMerchant(c.DB.QueryRowContext(ctx, q, args...))
}

// test for update merchant
// {
// 	"merchant": {
// 		"number_of_product":343,
// 		"email" :"merchanst@gmail.com",
// 		"role" : 1,
// 		"phone": {
// 			"type" : 1,
// 			"number": "040328844423"
// 		},
// 		"number_of_outlet": 343,
// 		"business_name": "mercshant"
// 	},
// 	"id": 2
// }
