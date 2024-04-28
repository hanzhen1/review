package data

import (
	"context"
	v1 "review-b/api/review/v1"
	"review-b/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type businessRepo struct {
	data *Data
	log  *log.Helper
}

// NewBusinessRepo .
func NewBusinessRepo(data *Data, logger log.Logger) biz.BusinessRepo {
	return &businessRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// ReplyReview 回复评价
func (r *businessRepo) ReplyReview(ctx context.Context, param *biz.ReplyParam) (int64, error) {
	r.log.WithContext(ctx).Debugf("[data] ReplyReview param:%v", param)
	//之前我们都是写操作数据库
	//而现在我们需要的是通过RPC调用其他的服务
	reply, err := r.data.rc.ReplyReview(ctx, &v1.ReplyReviewRequest{
		ReviewID:  param.ReviewID,
		StoreID:   param.StoreID,
		Content:   param.Content,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	})
	r.log.WithContext(ctx).Debugf("RPC ReplyReview return,reply:%v err:%v", reply, err)
	if err != nil {
		return 0, err
	}
	return reply.GetReplyID(), nil
}

// AppealReview 申诉评估价(商家对用户评价进行申诉)
func (r *businessRepo) AppealReview(ctx context.Context, param *biz.AppealParam) (int64, error) {
	r.log.WithContext(ctx).Debugf("[data] AppealReview param:%v", param)
	//之前我们都是写操作数据库
	//而现在我们需要的是通过RPC调用其他的服务
	appeal, err := r.data.rc.AppealReview(ctx, &v1.AppealReviewRequest{
		ReviewID:  param.ReviewID,
		StoreID:   param.StoreID,
		Reason:    param.Reason,
		Content:   param.Content,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	})
	r.log.WithContext(ctx).Debugf("RPC AppealReview return,appeal:%v err:%v", appeal, err)
	if err != nil {
		return 0, err
	}
	return appeal.GetAppealID(), err
}
