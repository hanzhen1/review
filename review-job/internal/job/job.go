package job

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewJobWorker, NewKafkaReader, NewESClient)
