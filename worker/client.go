package worker

import "github.com/hibiken/asynq"

type WorkerClient struct {
	Client *asynq.Client
}

func NewWorkerClient(redisAddr string) WorkerClient {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})
	return WorkerClient{Client: client}
}
