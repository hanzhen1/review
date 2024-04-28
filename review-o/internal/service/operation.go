package service

import (
	"context"

	pb "review-o/api/operation/v1"
	"review-o/internal/biz"
)

type OperationService struct {
	pb.UnimplementedOperationServer
	uc *biz.OperationUsecase
}

func NewOperationService(uc *biz.OperationUsecase) *OperationService {
	return &OperationService{uc: uc}
}

// AuditReview 审核评价 （运营对用户评价进行审核）
func (s *OperationService) AuditReview(ctx context.Context, req *pb.AuditReviewRequest) (*pb.AuditReviewReply, error) {
	err := s.uc.AuditReview(ctx, &biz.AuditReviewParam{
		ReviewID:  req.GetReviewID(),
		Status:    req.GetStatus(),
		OpUser:    req.GetOpUser(),
		OpReason:  req.GetOpReason(),
		OpRemarks: req.GetOpRemarks(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.AuditReviewReply{ReviewID: req.ReviewID, Status: req.Status}, nil
}

// AuditAppeal 申诉审核 （运营对商家申诉进行审核）
func (s *OperationService) AuditAppeal(ctx context.Context, req *pb.AuditAppealRequest) (*pb.AuditAppealReply, error) {
	err := s.uc.AuditAppeal(ctx, &biz.AuditAppealParam{
		AppealID:  req.GetAppealID(),
		ReviewID:  req.GetReviewID(),
		Status:    req.GetStatus(),
		OpUser:    req.GetOpUser(),
		OpRemarks: req.GetOpRemarks(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.AuditAppealReply{AppealID: req.AppealID, Status: req.Status}, nil
}
