package handler

import (
	"context"

	"github.com/Investliftng/ocm-api/merchant/datastore"
	pbM "github.com/Investliftng/ocm-api/merchant/proto/merchant"
	"github.com/micro/go-micro/v2/errors"
)

//Service _
type Service struct {
	Repository         datastore.MerchantRepo
	OutletRepository   datastore.OutletRepository
	TerminalRepository datastore.TerminalRepository
}

//CreateMerchant _
func (s *Service) CreateMerchant(ctx context.Context, req *pbM.CreateRequest, res *pbM.Response) error {
	if err := s.Repository.Create(ctx, req); err != nil {
		return errors.InternalServerError("500", "failed to create a new merchant", err.Error())
	}
	res.Created = true
	return nil
}

//UpdateMerchant _
func (s *Service) UpdateMerchant(ctx context.Context, req *pbM.UpdateRequest, res *pbM.Response) error {
	if err := s.Repository.Update(ctx, req); err != nil {
		return errors.InternalServerError("500", "failed to update merchant", err.Error())
	}
	res.Updated = true
	return nil
}

//GetMerchantByID _
func (s *Service) GetMerchantByID(ctx context.Context, req *pbM.GetIdRequest, res *pbM.Response) error {
	m, err := s.Repository.Read(ctx, req)
	if err != nil {
		return errors.InternalServerError("500", "failed to get a merchant", err.Error())
	}
	res.Merchant = m
	return nil
}

//GetMerchants _
func (s *Service) GetMerchants(ctx context.Context, req *pbM.Request, res *pbM.Response) error {
	ms, err := s.Repository.Get(ctx)
	if err != nil {
		return errors.InternalServerError("500", "failed to Get all merchant", err.Error())
	}
	res.Merchants = ms
	return nil
}

//DeleteMerchantByID delete a merchant by id
func (s *Service) DeleteMerchantByID(ctx context.Context, req *pbM.GetIdRequest, res *pbM.Response) error {
	if err := s.Repository.Delete(ctx, req); err != nil {
		return errors.InternalServerError("500", "failed to Delete merchant", err.Error())
	}
	res.Deleted = true
	return nil
}

//UpdatePassword _
func (s *Service) UpdatePassword(ctx context.Context, req *pbM.UpdatePasswordRequest, res *pbM.UpdatePasswordResponse) error {
	return nil
}

// Outlets methods

//CreateMerchantOutlet _
func (s *Service) CreateMerchantOutlet(ctx context.Context, req *pbM.OutletRequest, res *pbM.Response) error {
	if err := s.OutletRepository.CreateOutlet(ctx, req); err != nil {
		return errors.InternalServerError("500", "failed to create a new outlet", err.Error())
	}
	res.Created = true
	return nil
}

//DeleteMerchantOutlet _
func (s *Service) DeleteMerchantOutlet(ctx context.Context, req *pbM.GetOutletIdRequest, res *pbM.Response) error {
	if err := s.OutletRepository.DeleteOutlet(ctx, req); err != nil {
		return errors.InternalServerError("500", "failed to delete outlet", err.Error())
	}
	res.Deleted = true
	return nil
}

//UpdateMerchantOutlet _
func (s *Service) UpdateMerchantOutlet(ctx context.Context, req *pbM.UpdateOutletRequest, res *pbM.Response) error {
	if err := s.OutletRepository.UpdateOutlet(ctx, req); err != nil {
		return errors.InternalServerError("500", "failed to update outlet", err.Error())
	}
	res.Updated = true
	return nil
}

//GetMerchantOutlets _
func (s *Service) GetMerchantOutlets(ctx context.Context, req *pbM.GetOutletIdRequest, res *pbM.Response) error {
	outlet, err := s.OutletRepository.GetOutlet(ctx, req)
	if err != nil {
		return errors.InternalServerError("500", "failed to get by id outlets", err.Error())
	}
	res.Outlets = outlet
	return nil
}

//GetMerchantOutletByID _
func (s *Service) GetMerchantOutletByID(ctx context.Context, req *pbM.GetOutletIdRequest, res *pbM.Response) error {
	outlets, err := s.OutletRepository.ReadOutlet(ctx, req)
	if err != nil {
		return errors.InternalServerError("500", "failed to get all outlets", err.Error())
	}
	res.Outlet = outlets
	return nil
}

// Terminal methods

//CreateMerchantTerminal _
func (s *Service) CreateMerchantTerminal(ctx context.Context, req *pbM.TerminalRequest, res *pbM.Response) error {
	if err := s.TerminalRepository.CreateTerminal(ctx, req); err != nil {
		return errors.InternalServerError("500", "failed to create a new Terminal", err.Error())
	}
	res.Created = true
	return nil
}

//UpdateMerchantTerminal _
func (s *Service) UpdateMerchantTerminal(ctx context.Context, req *pbM.UpdateTerminalRequest, res *pbM.Response) error {
	if err := s.TerminalRepository.UpdateTerminal(ctx, req); err != nil {
		return errors.InternalServerError("500", "failed to update Terminal", err.Error())
	}
	res.Updated = true
	return nil
}

//GetMerchantTerminals _
func (s *Service) GetMerchantTerminals(ctx context.Context, req *pbM.GetTerminalIdRequest, res *pbM.Response) error {
	terminal, err := s.TerminalRepository.GetTerminals(ctx, req)
	if err != nil {
		return errors.InternalServerError("500", "failed to get by id Terminals", err.Error())
	}
	res.Terminals = terminal
	return nil
}

//GetMerchantTerminalByID _
func (s *Service) GetMerchantTerminalByID(ctx context.Context, req *pbM.GetTerminalIdRequest, res *pbM.Response) error {
	terminal, err := s.TerminalRepository.ReadTerminal(ctx, req)
	if err != nil {
		return errors.InternalServerError("500", "failed to get by id Terminals", err.Error())
	}
	res.Terminal = terminal
	return nil
}

//DeleteMerchantTerminal _
func (s *Service) DeleteMerchantTerminal(ctx context.Context, req *pbM.GetTerminalIdRequest, res *pbM.Response) error {
	if err := s.TerminalRepository.DeleteTerminal(ctx, req); err != nil {
		return errors.InternalServerError("500", "failed to delete Terminal", err.Error())
	}
	res.Deleted = true
	return nil
}
