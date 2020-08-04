package datastore

import (
	"context"

	pbM "github.com/Investliftng/ocm-api/merchant/proto/merchant"
	"github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes"
	pbError "github.com/micro/go-micro/v2/errors"
)

//OutletRepository communicates with the outlet database
type OutletRepository interface {
	CreateOutlet(context.Context, *pbM.OutletRequest) error
	UpdateOutlet(context.Context, *pbM.UpdateOutletRequest) error
	DeleteOutlet(context.Context, *pbM.GetOutletIdRequest) error
	ReadOutlet(context.Context, *pbM.GetOutletIdRequest) (*pbM.Outlet, error)
	GetOutlet(context.Context, *pbM.GetOutletIdRequest) ([]*pbM.Outlet, error)
}

//DeleteOutlet _
func (c *Connection) DeleteOutlet(ctx context.Context, req *pbM.GetOutletIdRequest) error {
	rows, err := c.SQLBuilder.Select(
		"merchant_id",
	).From("outlets").Where(
		squirrel.Eq{"id": req.GetId()},
	).QueryContext(ctx)
	if err != nil {
		return err
	}
	merchantID, err := scanOutletMerchantID(rows)
	if err != nil {
		return err
	}

	if merchantID != req.MerchantId {
		return pbError.BadRequest("404", "you do not have permission to Delete this outlet", err.Error())
	}

	_, err = c.SQLBuilder.Delete(
		"outlets",
	).Where(squirrel.Eq{"id": req.GetId()}).ExecContext(ctx)
	return err
}

//CreateOutlet _
func (c *Connection) CreateOutlet(ctx context.Context, req *pbM.OutletRequest) error {
	_, err := c.SQLBuilder.Insert(
		"outlets",
	).SetMap(map[string]interface{}{
		"merchant_id": req.GetMerchantId(),
		"phone":       req.Phone.GetNumber(),
		"phone_type":  (phoneTypeWrapper)(req.Phone.GetType()),
		"latitude":    req.GetLatitude(),
		"longitude":   req.GetLongitude(),
		"position":    req.GetPosition(),
		"city_id":     req.GetCityId(),
		"country_id":  req.GetCountryId(),
		"addr":        req.GetAddress(),
		"available":   req.GetAvailable(),
	}).ExecContext(ctx)
	return err
}

//ReadOutlet _
func (c *Connection) ReadOutlet(ctx context.Context, req *pbM.GetOutletIdRequest) (*pbM.Outlet, error) {
	rows, err := c.SQLBuilder.Select(
		"id",
		"merchant_id",
		"phone",
		"phone_type",
		"latitude",
		"longitude",
		"position",
		"city_id",
		"country_id",
		"addr",
		"available",
		"created_at",
		"updated_at",
	).From("outlets").Where(
		squirrel.Eq{"id": req.GetId()},
	).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	outlet, err := scanOutlet(rows)
	if err != nil {
		return nil, err
	}

	if outlet.MerchantId != req.MerchantId {
		return nil, pbError.BadRequest("404", "you do not have permission to view this outlet", err.Error())
	}
	return outlet, nil
}

//UpdateOutlet _
func (c *Connection) UpdateOutlet(ctx context.Context, req *pbM.UpdateOutletRequest) error {
	rows, err := c.SQLBuilder.Select(
		"merchant_id",
	).From("outlets").Where(
		squirrel.Eq{"id": req.GetId()},
	).QueryContext(ctx)
	if err != nil {
		return err
	}
	merchantID, err := scanOutletMerchantID(rows)
	if err != nil {
		return err
	}

	if merchantID != req.MerchantId {
		return pbError.BadRequest("404", "you do not have permission to update this outlet", err.Error())
	}
	//lets update db
	now := ptypes.TimestampNow()
	_, err = c.SQLBuilder.Update(
		"outlets",
	).SetMap(map[string]interface{}{
		"phone":      req.Outlet.GetPhone().GetNumber(),
		"phone_type": (phoneTypeWrapper)(req.Outlet.GetPhone().Type),
		"position":   req.Outlet.GetPosition(),
		"city_id":    req.Outlet.GetCityId(),
		"country_id": req.Outlet.GetCountryId(),
		"addr":       req.Outlet.GetAddress(),
		"available":  req.Outlet.GetAvailable(),
		"updated_at": (*timeWrapper)(now),
	}).Where(squirrel.Eq{"id": req.GetId()}).ExecContext(ctx)

	return err
}

//GetOutlet _
func (c *Connection) GetOutlet(ctx context.Context, req *pbM.GetOutletIdRequest) (outlets []*pbM.Outlet, err error) {
	rows, err := c.SQLBuilder.Select(
		"id",
		"merchant_id",
		"phone",
		"phone_type",
		"latitude",
		"longitude",
		"position",
		"city_id",
		"country_id",
		"addr",
		"available",
		"created_at",
		"updated_at",
	).From("outlets").Where(
		squirrel.Eq{"merchant_id": req.GetMerchantId()},
	).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := rows.Close()
		if err == nil && cerr != nil {
			err = pbError.InternalServerError("501", "service.Merchant.GetOutlet", err.Error())
		}
	}()

	for rows.Next() {
		outlet, err := scanOutlet(rows)
		if err != nil {
			return nil, err
		}
		outlets = append(outlets, outlet)
	}

	return
}
