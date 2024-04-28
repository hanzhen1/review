package job

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/segmentio/kafka-go"
	"review-job/internal/conf"
)

// 评价数据流处理任务

// JobWorker 自定义执行job的结构体，实现transport.Server
type JobWorker struct {
	kafkaReader *kafka.Reader //kafka reader
	esClient    *ESClient     //ES Client
	log         *log.Helper
}
type ESClient struct {
	*elasticsearch.TypedClient
	index string
}

func NewJobWorker(kafkaReader *kafka.Reader, esClient *ESClient, logger log.Logger) *JobWorker {
	return &JobWorker{
		kafkaReader: kafkaReader,
		esClient:    esClient,
		log:         log.NewHelper(logger),
	}
}
func NewKafkaReader(conf *conf.Kafka) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: conf.Brokers,
		GroupID: conf.GroupId, // 指定消费者组id
		Topic:   conf.Topic,
		// MaxBytes: 10e6, // 10MB
	})
}
func NewESClient(conf *conf.Elasticsearch) (*ESClient, error) {
	// ES 配置
	cfg := elasticsearch.Config{
		Addresses: conf.Addresses,
	}
	// 创建客户端连接
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ESClient{
		TypedClient: client,
		index:       conf.Index,
	}, nil
}

// Msg 定义Kafka中接收到的数据
type Msg struct {
	Type     string `json:"type"`
	Database string `json:"database"`
	Table    string `json:"table"`
	IsDdl    bool   `json:"isDdl"`
	Data     []map[string]interface{}
}

// Start kratos程序启动之后会调用的方法
// ctx 是kratos框架启动的时候传入的ctx,是带有退出取消的
func (jw JobWorker) Start(ctx context.Context) error {
	jw.log.Debug("JobWorker Start....")
	// 1.从Kafka中获取MySQL中的数据变更消息
	// 接收消息
	for {
		m, err := jw.kafkaReader.ReadMessage(ctx)
		if errors.Is(err, context.Canceled) { //错误是退出取消
			return nil
		}
		if err != nil {
			jw.log.Errorf("ReadMessage from kafka failed,err:%v", err)
			break
		}
		jw.log.Debugf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		// 2.将完整的评价数据写入ES
		msg := new(Msg)
		if err := json.Unmarshal(m.Value, msg); err != nil {
			jw.log.Errorf("Unmarshal msg from kafka failed,err:%v", err)
			continue
		}
		//补充
		//实际的业务场景可能需要在这增加一个步骤：对数据做业务处理
		//例如：把两张表的数据合成一个文档写入ES
		if msg.Type == "INSERT" {
			//往ES中新增文档
			for idx := range msg.Data {
				jw.indexDocument(msg.Data[idx])
			}
		} else {
			//往ES中更新文档
			for idx := range msg.Data {
				jw.updateDocument(msg.Data[idx])
			}
		}
	}
	return nil
}

// Stop kratos 程序结束之后调用的
func (jw JobWorker) Stop(ctx context.Context) error {
	jw.log.Debug("JobWorker Stop....")
	// 程序退出前关闭Reader
	return jw.kafkaReader.Close()
}

// indexDocument 创建索引文档
func (jw JobWorker) indexDocument(d map[string]interface{}) {
	// 添加文档
	reviewID := d["review_id"].(string)
	resp, err := jw.esClient.Index(jw.esClient.index).
		Id(reviewID).
		Document(d).
		Do(context.Background())
	if err != nil {
		jw.log.Errorf("indexing document failed, err:%v\n", err)
		return
	}
	jw.log.Debugf("result:%#v\n", resp.Result)
}

// updateDocument 更新文档
func (jw JobWorker) updateDocument(d map[string]interface{}) {
	reviewID := d["review_id"].(string)
	resp, err := jw.esClient.Update(jw.esClient.index, reviewID).
		Doc(d). // 使用结构体变量更新
		Do(context.Background())
	if err != nil {
		jw.log.Errorf("update document failed, err:%v\n", err)
		return
	}
	jw.log.Debugf("result:%v\n", resp.Result)
}
