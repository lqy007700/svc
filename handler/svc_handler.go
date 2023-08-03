package handler

import (
	"context"
	"errors"
	"github.com/zxnlx/common"
	"github.com/zxnlx/svc/domain/model"
	"github.com/zxnlx/svc/domain/service"
	"github.com/zxnlx/svc/proto/svc"
)

type SvcHandler struct {
	SvcDataService service.ISvcDataService
}

func (s *SvcHandler) AddSvc(ctx context.Context, info *svc.SvcInfo, resp *svc.Response) error {
	common.Info("Add Svc")
	svcModel := &model.Svc{}
	err := common.SwapTo(info, svcModel)
	if err != nil {
		common.Error(err)
		return err
	}

	err = s.SvcDataService.CreateSvcToK8s(info)
	if err != nil {
		common.Error(err)
		resp.Msg = err.Error()
		return err
	}

	addSvc, err := s.SvcDataService.AddSvc(svcModel)
	if err != nil {
		common.Error(err)
		resp.Msg = err.Error()
		return err
	}

	common.Info(addSvc)
	resp.Msg = "Add success"
	return nil
}

func (s *SvcHandler) DeleteSvc(ctx context.Context, req *svc.SvcId, resp *svc.Response) error {
	common.Info("Del Svc")
	info, err := s.SvcDataService.FindSvcById(req.Id)
	if err != nil {
		common.Error(err)
		resp.Msg = err.Error()
		return err
	}

	if info == nil {
		resp.Msg = "Svc Not Exist"
		common.Error(resp.Msg)
		return errors.New(resp.Msg)
	}

	err = s.SvcDataService.DeleteFromK8s(info)
	if err != nil {
		resp.Msg = "Svc Not Exist"
		common.Error(resp.Msg)
		return err
	}
	resp.Msg = "Del success"
	return nil
}

func (s *SvcHandler) UpdateSvc(ctx context.Context, info *svc.SvcInfo, resp *svc.Response) error {
	err := s.SvcDataService.UpdateSvcToK8s(info)
	if err != nil {
		return err
	}

	svcInfo, err := s.SvcDataService.FindSvcById(info.Id)
	if err != nil {
		return err
	}

	err = common.SwapTo(info, svcInfo)
	if err != nil {
		return err
	}

	err = s.SvcDataService.UpdateSvc(svcInfo)
	if err != nil {
		return err
	}
	return nil
}

func (s *SvcHandler) FindSvcByID(ctx context.Context, req *svc.SvcId, info *svc.SvcInfo) error {
	svcInfo, err := s.SvcDataService.FindSvcById(req.Id)
	if err != nil {
		return err
	}
	err = common.SwapTo(svcInfo, info)
	if err != nil {
		return err
	}
	return nil
}

func (s *SvcHandler) FindAllSvc(ctx context.Context, all *svc.FindAll, allSvc *svc.AllSvc) error {
	allSvcs, err := s.SvcDataService.FindAllSvc()
	if err != nil {
		return err
	}

	for _, v := range allSvcs {
		svcInfo := &svc.SvcInfo{}
		err = common.SwapTo(v, svcInfo)
		if err != nil {
			return err
		}

		allSvc.SvcInfo = append(allSvc.SvcInfo, svcInfo)
	}
	return nil
}
