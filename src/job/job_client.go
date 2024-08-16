package job

import "github.com/hibiken/asynq"

type Client interface {
	Enqueue(task *asynq.Task) (*asynq.TaskInfo, error)
	Close() error
}

type client struct {
	asynqClient *asynq.Client
}

func NewClient(r asynq.RedisConnOpt) Client {
	return &client{
		asynqClient: asynq.NewClient(r),
	}
}

func (c *client) Enqueue(task *asynq.Task) (*asynq.TaskInfo, error) {
	return c.asynqClient.Enqueue(task)
}

func (c *client) Close() error {
	return c.asynqClient.Close()
}
