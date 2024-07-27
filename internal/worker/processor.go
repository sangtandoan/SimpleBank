package worker

import (
	"context"

	"github.com/FrostJ143/simplebank/internal/query"
	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
	Start() error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  query.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store query.Store) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
	})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	serveMux := asynq.NewServeMux()

	serveMux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Run(serveMux)
}
