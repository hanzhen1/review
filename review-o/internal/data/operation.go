package data

import (
	"context"
	v1 "review-o/api/review/v1"

	"review-o/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type OperationRepo struct {
	data *Data
	log  *log.Helper
}

// NewOperationRepo .
func NewOperationRepo(data *Data, logger log.Logger) biz.OperationRepo {
	return &OperationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// AuditReview 审核评价 （运营对用户评价进行审核）
func (r *OperationRepo) AuditReview(ctx context.Context, param *biz.AuditReviewParam) error {
	r.log.WithContext(ctx).Debugf("[data] AuditReview param:%v", param)
	//之前我们都是写操作数据库
	//而现在我们需要的是通过RPC调用其他的服务
	ret, err := r.data.rc.AuditReview(ctx, &v1.AuditReviewRequest{
		ReviewID:  param.ReviewID,
		Status:    param.Status,
		OpUser:    param.OpUser,
		OpReason:  param.OpReason,
		OpRemarks: &param.OpRemarks,
	})
	r.log.WithContext(ctx).Debugf("rpc AuditReview return ,ret:%v err:%v", ret, err)
	if err != nil {
		return err
	}
	return nil
}

// AuditAppeal 申诉审核 （运营对商家申诉进行审核）
func (r *OperationRepo) AuditAppeal(ctx context.Context, param *biz.AuditAppealParam) error {
	r.log.WithContext(ctx).Debugf("[data] AuditAppeal param:%v", param)
	//之前我们都是写操作数据库
	//而现在我们需要的是通过RPC调用其他的服务
	ret, err := r.data.rc.AuditAppeal(ctx, &v1.AuditAppealRequest{
		AppealID:  param.AppealID,
		ReviewID:  param.ReviewID,
		Status:    param.Status,
		OpUser:    param.OpUser,
		OpRemarks: &param.OpRemarks,
	})
	r.log.WithContext(ctx).Debugf("rpc AuditAppeal return ,ret:%v err:%v", ret, err)
	if err != nil {
		return err
	}
	return nil
}
