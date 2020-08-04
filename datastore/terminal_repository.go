package datastore

import (
	"context"

	pbM "github.com/Investliftng/ocm-api/merchant/proto/merchant"
	"github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes"
	"github.com/micro/go-micro/v2/errors"
)

//TerminalRepository _
type TerminalRepository interface {
	CreateTerminal(context.Context, *pbM.TerminalRequest) error
	DeleteTerminal(context.Context, *pbM.GetTerminalIdRequest) error
	UpdateTerminal(context.Context, *pbM.UpdateTerminalRequest) error
	GetTerminals(context.Context, *pbM.GetTerminalIdRequest) ([]*pbM.Terminal, error)
	ReadTerminal(context.Context, *pbM.GetTerminalIdRequest) (*pbM.Terminal, error)
}

//CreateTerminal _
func (c *Connection) CreateTerminal(ctx context.Context, req *pbM.TerminalRequest) error {
	_, err := c.SQLBuilder.Insert(
		"terminals",
	).SetMap(map[string]interface{}{
		"merchant_Id":   req.GetMerchantId(),
		"user_id":       req.GetUserId(),
		"outlet_Id":     req.GetOutletId(),
		"terminal_name": req.GetName(),
	}).ExecContext(ctx)

	return err
}

//DeleteTerminal _
func (c *Connection) DeleteTerminal(ctx context.Context, req *pbM.GetTerminalIdRequest) error {
	var id uint64
	row, err := c.SQLBuilder.Select(
		"merchant_id",
	).From(
		"terminals",
	).Where(
		squirrel.Eq{"id": req.GetId()},
	).QueryContext(ctx)

	if err != nil {
		return err
	}
	if err := row.Scan(&id); err != nil {
		return err
	}
	if id != req.GetMerchantId() {
		return errors.BadRequest("404", "you do not have permission to Delete this Terminal")
	}

	_, err = c.SQLBuilder.Delete(
		"terminals",
	).Where(squirrel.Eq{"id": req.GetId()}).ExecContext(ctx)
	return err
}

//UpdateTerminal _
func (c *Connection) UpdateTerminal(ctx context.Context, req *pbM.UpdateTerminalRequest) error {
	var id uint64
	row, err := c.SQLBuilder.Select(
		"merchant_id",
	).From(
		"terminals",
	).Where(
		squirrel.Eq{"id": req.GetId()},
	).QueryContext(ctx)

	if err != nil {
		return err
	}
	if err := row.Scan(&id); err != nil {
		return err
	}
	if id != req.GetMerchantId() {
		return errors.BadRequest("404", "you do not have permission to Delete this Terminal")
	}

	now := ptypes.TimestampNow()
	_, err = c.SQLBuilder.Update(
		"terminals",
	).SetMap(
		map[string]interface{}{
			"user_id":       req.Terminal.GetUserId(),
			"outlet_Id":     req.Terminal.GetOutletId(),
			"terminal_name": req.Terminal.GetName(),
			"updated_at":    (*timeWrapper)(now),
		},
	).Where(
		squirrel.Eq{"id": req.GetId()},
	).ExecContext(ctx)
	return err
}

//GetTerminals _
func (c *Connection) GetTerminals(ctx context.Context, req *pbM.GetTerminalIdRequest) (terminals []*pbM.Terminal, err error) {
	rows, err := c.SQLBuilder.Select(
		"id",
		"merchant_Id",
		"user_id",
		"outlet_id",
		"terminal_name",
		"created_at",
		"updated_at",
	).From("terminals").Where(
		squirrel.Eq{"merchant_id": req.GetMerchantId()},
	).QueryContext(ctx)
	if err != nil {
		return
	}

	defer func() {
		cerr := rows.Close()
		if err == nil && cerr != nil {
			err = errors.InternalServerError("501", "service.Merchant.GetOutlet", err.Error())
		}
	}()

	for rows.Next() {
		terminal, err := scanTerminal(rows)
		if err != nil {
			return nil, err
		}
		terminals = append(terminals, terminal)
	}

	return
}

//ReadTerminal _
func (c *Connection) ReadTerminal(ctx context.Context, req *pbM.GetTerminalIdRequest) (terminal *pbM.Terminal, err error) {
	row, err := c.SQLBuilder.Select(
		"id",
		"merchant_Id",
		"user_id",
		"outlet_id",
		"terminal_name",
		"created_at",
		"updated_at",
	).From("terminals").Where(
		squirrel.Eq{"id": req.GetId()},
	).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	return scanTerminal(row)
}
