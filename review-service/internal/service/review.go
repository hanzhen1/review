package service

import (
	"context"
	"fmt"
	"review-service/internal/biz"
	"review-service/internal/data/model"

	pb "review-service/api/review/v1"
)

type ReviewService struct {
	pb.UnimplementedReviewServer
	uc *biz.ReviewUsecase
}

func NewReviewService(uc *biz.ReviewUsecase) *ReviewService {
	return &ReviewService{uc: uc}
}

// CreateReview C端创建评价
func (s *ReviewService) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.CreateReviewReply, error) {
	//fmt.Printf("[service] CreateReview req:%#v\n", req)
	//参数转换
	//调用biz层
	var anonymous int32
	if req.Anonymous {
		anonymous = 1
	}
	review, err := s.uc.CreateReview(ctx, &model.ReviewInfo{
		UserID:       req.UserID,
		OrderID:      req.OrderID,
		StoreID:      req.StoreID,
		Score:        req.Score,
		ServiceScore: req.ServiceScore,
		ExpressScore: req.ExpressScore,
		Content:      req.Content,
		PicInfo:      req.PicInfo,
		VideoInfo:    req.VideoInfo,
		Anonymous:    anonymous,
		Status:       0,
	})
	if err != nil {
		return nil, err
	}
	//拼装返回结果
	return &pb.CreateReviewReply{ReviewID: review.ReviewID}, nil
}

// GetReview C端根据评价ID获取评价
func (s *ReviewService) GetReview(ctx context.Context, req *pb.GetReviewRequest) (*pb.GetReviewReply, error) {
	fmt.Printf("[service] GetReview req:%#v\n", req)
	//调用biz层
	review, err := s.uc.GetReview(ctx, req.ReviewID)
	if err != nil {
		return nil, err
	}
	var anonymous bool
	if review.Anonymous == 1 {
		anonymous = true
	}
	//拼装返回结果
	return &pb.GetReviewReply{
		Data: &pb.ReviewInfo{
			ReviewID:     review.ReviewID,
			UserID:       review.UserID,
			OrderID:      review.OrderID,
			StoreID:      review.StoreID,
			Score:        review.Score,
			ServiceScore: review.ServiceScore,
			ExpressScore: review.ExpressScore,
			Status:       review.Status,
			Content:      review.Content,
			PicInfo:      review.PicInfo,
			VideoInfo:    review.VideoInfo,
			Anonymous:    anonymous,
		},
	}, nil
}

// AuditReview O端审核评价
func (s *ReviewService) AuditReview(ctx context.Context, req *pb.AuditReviewRequest) (*pb.AuditReviewReply, error) {
	fmt.Printf("[service] AuditReview req:%#v\n", req)
	//调用biz层
	err := s.uc.AuditReview(ctx, &biz.AuditParam{
		ReviewID:  req.GetReviewID(),
		OpUser:    req.GetOpUser(),
		OpReason:  req.GetOpReason(),
		OpRemarks: req.GetOpRemarks(),
		Status:    req.GetStatus(),
	})
	if err != nil {
		return nil, err
	}
	//拼装返回结果
	return &pb.AuditReviewReply{ReviewID: req.ReviewID, Status: req.Status}, nil
}

// ReplyReview B端回复评价
func (s *ReviewService) ReplyReview(ctx context.Context, req *pb.ReplyReviewRequest) (*pb.ReplyReviewReply, error) {
	fmt.Printf("[service] ReplyReview req:%#v\n", req)
	//调用biz层
	reply, err := s.uc.CreateReply(ctx, &biz.ReplyParam{
		ReviewID:  req.GetReviewID(),
		StoreID:   req.GetStoreID(),
		Content:   req.GetContent(),
		PicInfo:   req.GetPicInfo(),
		VideoInfo: req.GetVideoInfo(),
	})
	if err != nil {
		return nil, err
	}
	//拼装返回结果
	return &pb.ReplyReviewReply{ReplyID: reply.ReplyID}, nil
}

// GetReplyReview B端获取评价回复
func (s *ReviewService) GetReplyReview(ctx context.Context, req *pb.GetReplyReviewRequest) (*pb.GetReplyReviewReply, error) {
	fmt.Printf("[service] GetReplyReview req:%#v\n", req)
	reply, err := s.uc.GetReplyReview(ctx, req.ReviewID)
	if err != nil {
		return nil, err
	}
	return &pb.GetReplyReviewReply{Data: &pb.ReplyInfo{
		ReplyID:   reply.ReplyID,
		ReviewID:  reply.ReviewID,
		StoreID:   reply.StoreID,
		Content:   reply.Content,
		PicInfo:   reply.PicInfo,
		VideoInfo: reply.VideoInfo,
	}}, nil
}

// AppealReview B端申诉评价
func (s *ReviewService) AppealReview(ctx context.Context, req *pb.AppealReviewRequest) (*pb.AppealReviewReply, error) {
	fmt.Printf("[service] AppealReview req:%#v\n", req)
	//调用biz层
	appeal, err := s.uc.AppealReview(ctx, &biz.AppealParam{
		ReviewID:  req.GetReviewID(),
		StoreID:   int64(req.GetStoreID()),
		Reason:    req.GetReason(),
		Content:   req.GetContent(),
		PicInfo:   req.GetPicInfo(),
		VideoInfo: req.GetVideoInfo(),
	})
	if err != nil {
		return nil, err
	}
	//拼装返回结果
	return &pb.AppealReviewReply{AppealID: appeal.AppealID}, nil
}

// AuditAppeal O端评价审核申诉
func (s *ReviewService) AuditAppeal(ctx context.Context, req *pb.AuditAppealRequest) (*pb.AuditAppealReply, error) {
	fmt.Printf("[service] AuditAppeal req:%#v\n", req)
	//调用biz层
	err := s.uc.AuditAppeal(ctx, &biz.AuditAppealParam{
		ReviewID: req.GetReviewID(),
		AppealID: req.GetAppealID(),
		OpUser:   req.GetOpUser(),
		Status:   req.GetStatus(),
	})
	if err != nil {
		return nil, err
	}

	//拼装返回结果
	return &pb.AuditAppealReply{AppealID: req.AppealID, Status: req.Status}, nil
}

// ListReviewByUserID C端查看userID下所有评价
func (s *ReviewService) ListReviewByUserID(ctx context.Context, req *pb.ListReviewByUserIDRequest) (*pb.ListReviewByUserIDReply, error) {
	fmt.Printf("[service] ListReviewByUserID req:%#v\n", req)
	//调用biz层
	reviewList, err := s.uc.ListReviewByUserID(ctx, req.UserID, int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}
	//format
	list := make([]*pb.ReviewInfo, 0, len(reviewList))
	for _, r := range reviewList {
		var anonymous bool
		if r.Anonymous == 1 {
			anonymous = true
		}
		list = append(list, &pb.ReviewInfo{
			ReviewID:     r.ReviewID,
			UserID:       r.UserID,
			OrderID:      r.OrderID,
			StoreID:      r.StoreID,
			Score:        r.Score,
			ServiceScore: r.ServiceScore,
			ExpressScore: r.ExpressScore,
			Status:       r.Status,
			Content:      r.Content,
			PicInfo:      r.PicInfo,
			VideoInfo:    r.VideoInfo,
			Anonymous:    anonymous,
		})
	}
	//拼装返回结果
	return &pb.ListReviewByUserIDReply{List: list}, nil
}

// ListReviewByStoreID 商家根据storeID查询评价列表
func (s *ReviewService) ListReviewByStoreID(ctx context.Context, req *pb.ListReviewByStoreIDRequest) (*pb.ListReviewByStoreIDReply, error) {
	fmt.Printf("[service] ListReviewByStoreID req:%#v\n", req)
	//调用biz层
	reviewList, err := s.uc.ListReviewByStoreID(ctx, req.StoreID, int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}
	//format
	list := make([]*pb.ReviewInfo, 0, len(reviewList))
	for _, r := range reviewList {
		var anonymous bool
		if r.Anonymous == 1 {
			anonymous = true
		}
		list = append(list, &pb.ReviewInfo{
			ReviewID:     r.ReviewID,
			UserID:       r.UserID,
			OrderID:      r.OrderID,
			StoreID:      r.StoreID,
			Score:        r.Score,
			ServiceScore: r.ServiceScore,
			ExpressScore: r.ExpressScore,
			Status:       r.Status,
			Content:      r.Content,
			PicInfo:      r.PicInfo,
			VideoInfo:    r.VideoInfo,
			Anonymous:    anonymous,
		})
	}
	//拼装返回结果
	return &pb.ListReviewByStoreIDReply{List: list}, nil
}
