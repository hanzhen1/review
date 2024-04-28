package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type ReplyParam struct {
	ReviewID  int64
	StoreID   int64
	Content   string
	PicInfo   string
	VideoInfo string
}
type AppealParam struct {
	ReviewID  int64
	StoreID   int64
	Reason    string
	Content   string
	PicInfo   string
	VideoInfo string
}

// BusinessRepo is a Business repo.
type BusinessRepo interface {
	ReplyReview(context.Context, *ReplyParam) (int64, error)
	AppealReview(context.Context, *AppealParam) (int64, error)
}

// BusinessUsecase is a Business usecase.
type BusinessUsecase struct {
	repo BusinessRepo
	log  *log.Helper
}

// NewBusinessUsecase new a Business usecase.
func NewBusinessUsecase(repo BusinessRepo, logger log.Logger) *BusinessUsecase {
	return &BusinessUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateReplyReview  创建评价回复
// service层调用此方法
func (uc *BusinessUsecase) CreateReplyReview(ctx context.Context, param *ReplyParam) (int64, error) {
	uc.log.WithContext(ctx).Debugf("[biz] CreateReplyReview param:%v", param)
	return uc.repo.ReplyReview(ctx, param)
}

// CreateAppealReview 申诉评价(商家对用户评价进行申诉)
// service层调用此方法
func (uc *BusinessUsecase) CreateAppealReview(ctx context.Context, param *AppealParam) (int64, error) {
	uc.log.WithContext(ctx).Debugf("[biz] CreateAppealReview param:%v", param)
	return uc.repo.AppealReview(ctx, param)
}
