package worker

import (
	"context"

	"github.com/FrostJ143/simplebank/internal/email"
	"github.com/FrostJ143/simplebank/internal/query"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
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
	server     *asynq.Server
	store      query.Store
	mailSender email.EmailSender
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store query.Store, mailSender email.EmailSender) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).Str("type", task.Type()).
				Bytes("payload", task.Payload()).Msg("processed task failed")
		}),
		Logger: NewLogger(),
	})

	return &RedisTaskProcessor{
		server:     server,
		store:      store,
		mailSender: mailSender,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	serveMux := asynq.NewServeMux()

	serveMux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Run(serveMux)
}
