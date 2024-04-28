package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type AuditReviewParam struct {
	ReviewID  int64
	Status    int32
	OpUser    string
	OpReason  string
	OpRemarks string
}
type AuditAppealParam struct {
	AppealID  int64
	ReviewID  int64
	Status    int32
	OpUser    string
	OpRemarks string
}
type OperationRepo interface {
	AuditReview(context.Context, *AuditReviewParam) error
	AuditAppeal(context.Context, *AuditAppealParam) error
}

type OperationUsecase struct {
	repo OperationRepo
	log  *log.Helper
}

// NewOperationUsecase 构造函数
func NewOperationUsecase(repo OperationRepo, logger log.Logger) *OperationUsecase {
	return &OperationUsecase{repo: repo, log: log.NewHelper(logger)}
}

// AuditReview 审核评价
// service层调用此方法
func (uc *OperationUsecase) AuditReview(ctx context.Context, param *AuditReviewParam) error {
	uc.log.WithContext(ctx).Debugf("AuditReview param: %v", param)
	return uc.repo.AuditReview(ctx, param)
}

// AuditAppeal 申诉审核
// service层调用此方法
func (uc *OperationUsecase) AuditAppeal(ctx context.Context, param *AuditAppealParam) error {
	uc.log.WithContext(ctx).Debugf("AuditAppeal param: %v", param)
	return uc.repo.AuditAppeal(ctx, param)
}
