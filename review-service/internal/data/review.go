package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"review-service/internal/data/model"
	"review-service/internal/data/query"
	"review-service/pkg/snowflake"
	"strconv"
	"strings"
	"time"

	"review-service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type reviewRepo struct {
	data *Data
	log  *log.Helper
}

// NewReviewRepo .
func NewReviewRepo(data *Data, logger log.Logger) biz.ReviewRepo {
	return &reviewRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *reviewRepo) SaveReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	err := r.data.query.ReviewInfo.WithContext(ctx).Save(review)
	return review, err
}

// GetReviewByOrderID 根据订单服务查询评价
func (r *reviewRepo) GetReviewByOrderID(ctx context.Context, orderID int64) ([]*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewInfo.OrderID.Eq(orderID)).
		Find()
}

// GetReview 根据评价ID获取评价
func (r *reviewRepo) GetReview(ctx context.Context, reviewID int64) (*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewInfo.ReviewID.Eq(reviewID)).
		First()
}

// SaveReply 保存评价回复
func (r *reviewRepo) SaveReply(ctx context.Context, reply *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error) {
	//1.数据校验
	//1.1 数据合法性校验(已回复的评价不允许商家再次回复)
	//先用评价ID查库，看是否已回复
	review, err := r.GetReview(ctx, reply.ReviewID)
	if err != nil {
		return nil, err
	}
	if review.HasReply == 1 {
		return nil, errors.New("该评价已回复")
	}
	//1.2.水平越权校验(A商家只能回复自己的不能回复B商家的)
	//举例子：用户A删除订单，userID + orderID 去查询订单然后删除
	if review.StoreID != reply.StoreID {
		return nil, errors.New("水平越权")
	}
	//2.更新数据库中的数据(评价回复表和评价表要同时更新，涉及到事务操作)
	//事务参数
	err = r.data.query.Transaction(func(tx *query.Query) error {
		//回复表插入一条数据
		if err := tx.ReviewReplyInfo.WithContext(ctx).
			Save(reply); err != nil {
			r.log.WithContext(ctx).Errorf("SaveReply create reply fail,err:%v\n", err)
			return err
		}
		//评价表更新HasReply字段
		if _, err := tx.ReviewInfo.WithContext(ctx).
			Where(tx.ReviewInfo.ReviewID.Eq(reply.ReviewID)).
			Update(tx.ReviewInfo.HasReply, 1); err != nil {
			r.log.WithContext(ctx).Errorf("SaveReply update review fail,err:%v\n", err)
			return err
		}
		return nil
	})
	//3.返回
	return reply, err
}

// GetReviewReply 根据reviewID获取评价回复表
func (r *reviewRepo) GetReviewReply(ctx context.Context, reviewID int64) (*model.ReviewReplyInfo, error) {
	return r.data.query.ReviewReplyInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewReplyInfo.ReviewID.Eq(reviewID)).
		First()
}

// AuditReview 审核评价（运营对用户的评价进行审核）
func (r *reviewRepo) AuditReview(ctx context.Context, param *biz.AuditParam) error {
	// 数据合法性校验
	//用评价ID查库看是否已审核
	review, err := r.GetReview(ctx, param.ReviewID)
	if err != nil {
		return err
	}
	if review.Status > 10 {
		return errors.New("该评价已有审核过的记录")
	}
	_, err = r.data.query.ReviewInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewInfo.ReviewID.Eq(param.ReviewID)).
		Updates(map[string]interface{}{
			"status":     param.Status,
			"op_user":    param.OpUser,
			"op_reason":  param.OpReason,
			"op_remarks": param.OpRemarks,
		})
	return err
}

// AppealReview 申诉评价（商家对用户评价进行申诉）
func (r *reviewRepo) AppealReview(ctx context.Context, param *biz.AppealParam) (*model.ReviewAppealInfo, error) {
	//1.1 数据合法性校验
	//先用评价ID查库，看用户是否已评价
	review, err := r.GetReview(ctx, param.ReviewID)
	if err != nil { //查数据库错误
		return nil, err
	}
	//1.2.水平越权校验(A商家只能申诉自己店铺的评价不能回复B商家的)
	//举例子：商家A申诉评价，storeID + reviewID 去查询评价id然后进行申诉
	if review.StoreID != param.StoreID {
		return nil, errors.New("水平越权")
	}
	//2.先查询有没有申诉
	ret, err := r.data.query.ReviewAppealInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewAppealInfo.ReviewID.Eq(param.ReviewID),
			r.data.query.ReviewAppealInfo.StoreID.Eq(param.StoreID),
		).First()
	r.log.Debugf("AppealReview query, ret:%v err:%v", ret, err)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err //其他查询错误
	}
	if err == nil && ret.Status > 10 {
		return nil, errors.New("该评价已有审核过的申诉记录")
	}
	// 查询不到审核过的申诉记录
	// 2.1. 有申诉记录但是处于待审核状态，需要更新
	// if ret != nil{
	// 	// update
	// }else{
	// 	// insert
	// }
	// 2.2. 没有申诉记录，需要创建
	appeal := &model.ReviewAppealInfo{
		ReviewID:  param.ReviewID,
		StoreID:   param.StoreID,
		Status:    10,
		Reason:    param.Reason,
		Content:   param.Content,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	}
	if ret != nil {
		appeal.AppealID = ret.AppealID
	} else {
		appeal.AppealID = snowflake.GenID()
	}
	//如果存在则更新否则新增数据
	//INSERT INTO `table` *** ON DUPLICATE KEY UPDATE ***;
	err = r.data.query.ReviewAppealInfo.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "review_id"}, // ON DUPLICATE KEY
			},
			DoUpdates: clause.Assignments(map[string]interface{}{ //update
				"status":     appeal.Status,
				"reason":     appeal.Reason,
				"content":    appeal.Content,
				"pic_info":   appeal.PicInfo,
				"video_info": appeal.VideoInfo,
			}),
		}).
		Create(appeal) // INSERT
	r.log.Debugf("AppealReview, err:%v", err)
	if err != nil {
		return nil, err
	}
	return appeal, nil
}

// AuditAppeal 评价审核申诉(运营对商家的申诉进行审核，审核通过会隐藏该评价)
func (r *reviewRepo) AuditAppeal(ctx context.Context, param *biz.AuditAppealParam) error {
	//1.数据合法性校验
	//查询AppealID,ReviewID是否存在
	ret, err := r.data.query.ReviewAppealInfo.WithContext(ctx).
		Where(r.data.query.ReviewAppealInfo.ReviewID.Eq(param.ReviewID)).
		Where(r.data.query.ReviewAppealInfo.AppealID.Eq(param.AppealID)).
		First()
	if err != nil {
		return err
	}
	if ret.Status > 10 {
		return errors.New("该申诉已有审核过的申诉记录")
	}
	//2.更新数据库中的数据(申诉表和评价表要同时更新，涉及到事务操作)
	err = r.data.query.Transaction(func(tx *query.Query) error {
		// 更新申诉表状态和运营者标识
		if _, err := tx.ReviewAppealInfo.
			WithContext(ctx).
			Where(tx.ReviewAppealInfo.AppealID.Eq(param.AppealID)).
			Updates(map[string]interface{}{
				"status":  param.Status,
				"op_user": param.OpUser,
			}); err != nil {
			return err
		}
		//更新评价表状态
		if param.Status == 20 { // 申诉通过则需要隐藏评价
			if _, err := tx.ReviewInfo.WithContext(ctx).
				Where(tx.ReviewInfo.ReviewID.Eq(param.ReviewID)).
				Update(tx.ReviewInfo.Status, 40); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// ListReviewByUserID 根据userID查询所有评价
func (r *reviewRepo) ListReviewByUserID(ctx context.Context, userID int64, offset, limit int) ([]*model.ReviewInfo, error) {
	// 数据合法性校验
	//查库判断userid是否存在
	if _, err := r.data.query.ReviewInfo.WithContext(ctx).
		Where(r.data.query.ReviewInfo.UserID.Eq(userID)).
		First(); err != nil {
		return nil, err
	}
	return r.data.query.ReviewInfo.
		WithContext(ctx).
		Where(r.data.query.ReviewInfo.UserID.Eq(userID)).
		Order(r.data.query.ReviewInfo.ID.Desc()).
		Offset(offset).
		Limit(limit).
		Find()
}

var g singleflight.Group

// ListReviewByStoreID 根据storeID查询所有评价
func (r *reviewRepo) ListReviewByStoreID(ctx context.Context, storeID int64, offset, limit int) ([]*biz.MyReviewInfo, error) {
	//return r.getData1(ctx, storeID, offset, limit) //第一版直接查ES
	return r.getData2(ctx, storeID, offset, limit) //第二版增加缓存和singleflight
}

// getData1 直接查ES
func (r *reviewRepo) getData1(ctx context.Context, storeID int64, offset, limit int) ([]*biz.MyReviewInfo, error) {
	//去ES里面查询评价
	resp, err := r.data.es.Search().
		Index("review").
		From(offset).
		Size(limit).Query(&types.Query{
		Bool: &types.BoolQuery{
			Filter: []types.Query{
				{
					Term: map[string]types.TermQuery{
						"store_id": {Value: storeID},
					},
				},
			},
		},
	}).Do(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Printf("es result total: %d\n", resp.Hits.Total.Value)

	//反序列化数据
	//resp.Hits.Hits[0].Source_(json.RawMessage) ==>model.ReviewInfo
	// 遍历所有结果
	list := make([]*biz.MyReviewInfo, 0, resp.Hits.Total.Value) //知道切片的长度后续不用动态扩容
	for _, hit := range resp.Hits.Hits {
		temp := &biz.MyReviewInfo{}
		if err := json.Unmarshal(hit.Source_, temp); err != nil {
			r.log.Errorf("json.Unmarshal(hit.Source_, temp) failed,err:%v", err)
			continue
		}
		list = append(list, temp)
	}
	return list, nil
}

// getData2 升级版带缓存版本的查询函数
func (r *reviewRepo) getData2(ctx context.Context, storeID int64, offset, limit int) ([]*biz.MyReviewInfo, error) {
	//取数据
	//1.先查询redis缓存
	//2.缓存没有则查询ES
	//3.通过singleflight合并短时间内大量的并发查询
	key := fmt.Sprintf("review:%d:%d:%d", storeID, offset, limit)
	data, err := r.getDataBySingleflight(ctx, key)
	if err != nil {
		return nil, err
	}
	hm := new(types.HitsMetadata)
	if err := json.Unmarshal(data, hm); err != nil {
		return nil, err
	}
	//反序列化
	//反序列化数据
	//resp.Hits.Hits[0].Source_(json.RawMessage) ==>model.ReviewInfo
	// 遍历所有结果
	list := make([]*biz.MyReviewInfo, 0, hm.Total.Value) //知道切片的长度后续不用动态扩容
	for _, hit := range hm.Hits {
		temp := &biz.MyReviewInfo{}
		if err := json.Unmarshal(hit.Source_, temp); err != nil {
			r.log.Errorf("json.Unmarshal(hit.Source_, temp) failed,err:%v", err)
			continue
		}
		list = append(list, temp)
	}
	return list, nil
}

// key review:storeID:page:size -->"[{},{},{}]"
// json.Unmarshal([]byte, any)

func (r *reviewRepo) getDataBySingleflight(ctx context.Context, key string) ([]byte, error) {
	v, err, shared := g.Do(key, func() (interface{}, error) {
		//查缓存
		data, err := r.getDataFromCache(ctx, key)
		r.log.Debugf("r.getDataFromCache(ctx, key) data:%s err:%v\n", data, err)
		if err == nil {
			return data, nil
		}
		//只有在缓存中没有这个key的错误时才查ES
		if errors.Is(err, redis.Nil) {
			//缓存中没有这个key,说明缓存失效了，需要查ES
			data, err := r.getDataFromES(ctx, key)
			if err == nil {
				//设置缓存
				return data, r.setCache(ctx, key, data)
			}
			return nil, err
		}
		//查缓存失败了，直接返回错误，不继续向下传导压力
		return nil, err
	})
	r.log.Debugf("singleflight ret: v:%s  err:%v shared:%v\n", v, err, shared)
	if err != nil {
		return nil, err
	}
	return v.([]byte), nil
}

// getDataFromCache 读缓存
func (r *reviewRepo) getDataFromCache(ctx context.Context, key string) ([]byte, error) {
	r.log.Debugf("getDataFromCache key:%v\n", key)
	return r.data.rdb.Get(ctx, key).Bytes() //之所以返回Bytes()类型是为了和ES中一致
}

// setCache 设置缓存
func (r *reviewRepo) setCache(ctx context.Context, key string, data []byte) error {
	return r.data.rdb.Set(ctx, key, data, 10*time.Second).Err()
}

// getDataFromES 从ES查询
func (r *reviewRepo) getDataFromES(ctx context.Context, key string) ([]byte, error) {
	values := strings.Split(key, ":")
	if len(values) < 4 {
		return nil, errors.New("invalid key")
	}
	index, storeID, offsetStr, limitStr := values[0], values[1], values[2], values[3]
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return nil, err
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil, err
	}
	resp, err := r.data.es.Search().
		Index(index).
		From(offset).
		Size(limit).Query(&types.Query{
		Bool: &types.BoolQuery{
			Filter: []types.Query{
				{
					Term: map[string]types.TermQuery{
						"store_id": {Value: storeID},
					},
				},
			},
		},
	}).Do(ctx)
	if err != nil {
		return nil, err
	}
	return json.Marshal(resp.Hits)
}
