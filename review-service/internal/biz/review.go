package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	v1 "review-service/api/review/v1"
	"review-service/internal/data/model"
	"review-service/pkg/snowflake"
	"strings"
	"time"
)

// ReviewRepo is a Review repo.
type ReviewRepo interface {
	SaveReview(context.Context, *model.ReviewInfo) (*model.ReviewInfo, error)
	GetReviewByOrderID(context.Context, int64) ([]*model.ReviewInfo, error)
	GetReview(context.Context, int64) (*model.ReviewInfo, error)
	SaveReply(context.Context, *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error)
	GetReviewReply(context.Context, int64) (*model.ReviewReplyInfo, error)
	AuditReview(context.Context, *AuditParam) error
	AppealReview(context.Context, *AppealParam) (*model.ReviewAppealInfo, error)
	AuditAppeal(context.Context, *AuditAppealParam) error
	ListReviewByUserID(ctx context.Context, userID int64, offset, limit int) ([]*model.ReviewInfo, error)
	ListReviewByStoreID(ctx context.Context, storeID int64, offset, limit int) ([]*MyReviewInfo, error)
}

// ReviewUsecase is a Review usecase.
type ReviewUsecase struct {
	repo ReviewRepo
	log  *log.Helper
}

// NewReviewUsecase new a Review usecase.
func NewReviewUsecase(repo ReviewRepo, logger log.Logger) *ReviewUsecase {
	return &ReviewUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateReview 创建评价 实现业务逻辑的地方
// service调用该方法
func (uc *ReviewUsecase) CreateReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	uc.log.WithContext(ctx).Debugf("[biz] CreateReview req: %#v", review)
	//1.数据校验
	//1.1参数基础校验：正常来说不应该放在这一层，你应该在上一层或者框架层都应该能拦住（validate参数校验）
	//1.2业务逻辑校验：带业务逻辑的参数校验，比如已经评价过的订单不能再创建评价
	reviews, err := uc.repo.GetReviewByOrderID(ctx, review.OrderID)
	if err != nil {
		return nil, v1.ErrorDbFailed("查询数据库失败")
	}
	if len(reviews) > 0 {
		//已经评价过
		return nil, v1.ErrorOrderReviewed("订单%v已被评价", review.OrderID)
	}
	fmt.Printf("%s", reviews)
	//2.生成reviewID(雪花算法或者直接接入公司内部的分布式ID生成服务，前提公司内部有这种服务)
	review.ReviewID = snowflake.GenID()
	//3.查询订单和商品快照信息
	//实际业务场景下就需要查询订单服务和商家服务（比如说通过RPC调用订单服务和商家服务）
	//4.拼装数据入库
	return uc.repo.SaveReview(ctx, review)
}

// GetReview 根据评价ID获取评价
func (uc *ReviewUsecase) GetReview(ctx context.Context, reviewID int64) (*model.ReviewInfo, error) {
	uc.log.WithContext(ctx).Debugf("[biz] GetReview reviewID: %#v", reviewID)
	return uc.repo.GetReview(ctx, reviewID)
}

// CreateReply 创建评价回复
func (uc *ReviewUsecase) CreateReply(ctx context.Context, param *ReplyParam) (*model.ReviewReplyInfo, error) {
	//调用data层创建一个评价的回复
	uc.log.WithContext(ctx).Debugf("[biz] CreateReply: %#v", param)
	reply := &model.ReviewReplyInfo{
		ReplyID:   snowflake.GenID(),
		ReviewID:  param.ReviewID,
		StoreID:   param.StoreID,
		Content:   param.Content,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	}
	return uc.repo.SaveReply(ctx, reply)
}

// GetReplyReview 获取评价回复
func (uc *ReviewUsecase) GetReplyReview(ctx context.Context, reviewID int64) (*model.ReviewReplyInfo, error) {
	uc.log.WithContext(ctx).Debugf("[biz] GetReplyReview reviewID: %#v", reviewID)
	return uc.repo.GetReviewReply(ctx, reviewID)
}

// AuditReview 审核评价
func (uc *ReviewUsecase) AuditReview(ctx context.Context, param *AuditParam) error {
	uc.log.WithContext(ctx).Debugf("[biz] AuditReview param: %#v", param)
	return uc.repo.AuditReview(ctx, param)
}

// AppealReview 申诉评价
func (uc *ReviewUsecase) AppealReview(ctx context.Context, param *AppealParam) (*model.ReviewAppealInfo, error) {
	uc.log.WithContext(ctx).Debugf("[biz] AppealReview param: %#v", param)
	return uc.repo.AppealReview(ctx, param)
}

// AuditAppeal 对评价进行申诉审核
func (uc *ReviewUsecase) AuditAppeal(ctx context.Context, param *AuditAppealParam) error {
	uc.log.WithContext(ctx).Debugf("[biz] AuditAppeal param: %#v", param)
	return uc.repo.AuditAppeal(ctx, param)
}

// ListReviewByUserID 根据userID分页查询所有评价
func (uc *ReviewUsecase) ListReviewByUserID(ctx context.Context, userID int64, page, size int) ([]*model.ReviewInfo, error) {
	uc.log.WithContext(ctx).Debugf("[biz] ListReviewByUserID userID:%#v page:%#v size:%#v", userID, page, size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 10
	}
	offset := (page - 1) * size
	limit := size
	return uc.repo.ListReviewByUserID(ctx, userID, offset, limit)
}

// ListReviewByStoreID 根据storeID分页查询所有评价
func (uc *ReviewUsecase) ListReviewByStoreID(ctx context.Context, storeID int64, page, size int) ([]*MyReviewInfo, error) {
	uc.log.WithContext(ctx).Debugf("[biz] ListReviewByStoreID storeID:%#v page:%#v size:%#v", storeID, page, size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 10
	}
	offset := (page - 1) * size
	limit := size
	return uc.repo.ListReviewByStoreID(ctx, storeID, offset, limit)
}

type MyReviewInfo struct {
	*model.ReviewInfo
	CreateAt     MyTime `json:"create_at"`
	UpdateAt     MyTime `json:"update_at"`
	ID           int64  `json:"id,string"`
	Version      int32  `json:"version,string"`
	ReviewID     int64  `json:"review_id,string"`
	Score        int32  `json:"score,string"`
	ServiceScore int32  `json:"service_score,string"`
	ExpressScore int32  `json:"express_score,string"`
	HasMedia     int32  `json:"has_media,string"`
	OrderID      int64  `json:"order_id,string"`
	SkuID        int64  `json:"sku_id,string"`
	SpuID        int64  `json:"spu_id,string"`
	StoreID      int64  `json:"store_id,string"`
	UserID       int64  `json:"user_id,string"`
	Anonymous    int32  `json:"anonymous,string"`
	Status       int32  `json:"status,string"`
	IsDefault    int32  `json:"is_default,string"`
	HasReply     int32  `json:"has_reply,string"`
}
type MyTime time.Time

// UnmarshalJSON json.Unmarshal的时候会自动调用该方法
func (t *MyTime) UnmarshalJSON(data []byte) error {
	//data = "\"2023-12-17 14:20:18\""
	s := strings.Trim(string(data), `"`)
	temp, err := time.Parse(time.DateTime, s)
	if err != nil {
		return err
	}
	*t = MyTime(temp)
	return nil
}
