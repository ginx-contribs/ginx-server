package handler

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/ginx-contribs/ginx-server/internal/conf"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/types"
	"github.com/ginx-contribs/ginx-server/pkg/email"
	"github.com/ginx-contribs/ginx-server/pkg/mq"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
	"github.com/ginx-contribs/str2bytes"
	"github.com/wneessen/go-mail"
	"golang.org/x/net/context"
)

func NewEmailHandler(cfg conf.Email, sender *email.Sender, queue mq.Queue) (EmailHandler, error) {
	handler := EmailHandler{Config: cfg, Sender: sender, Queue: queue}

	// subscribe the Queue
	for _, consumer := range cfg.MQ.Consumers {
		c := &EmailConsumer{
			topic:     cfg.MQ.Topic,
			group:     cfg.MQ.Group,
			name:      consumer,
			batchSize: cfg.MQ.BatchSize,
			sender:    sender,
		}
		if err := queue.Subscribe(c); err != nil {
			return handler, err
		}
	}
	return handler, nil
}

// EmailHandler is responsible for publishing and sending emails
type EmailHandler struct {
	Config conf.Email
	Sender *email.Sender

	Queue mq.Queue
}

// Publish publishes message to Queue
func (e *EmailHandler) Publish(ctx context.Context, msg email.Message) error {
	msg.From = e.Config.Username
	marshal, err := sonic.Marshal(msg)
	if err != nil {
		return statuserr.InternalError(err)
	}
	_, err = e.Queue.Publish(ctx, e.Config.MQ.Topic, map[string]any{"mail": marshal}, 0)
	if err != nil {
		return statuserr.InternalError(err)
	}
	return err
}

func (e *EmailHandler) buildMail(msg types.EmailBody) (*mail.Msg, error) {
	mailMsg := mail.NewMsg()
	steps := []error{
		mailMsg.From(msg.From),
		mailMsg.To(msg.To...),
		mailMsg.Cc(msg.CC...),
		mailMsg.Bcc(msg.Bcc...),
	}
	for _, err := range steps {
		if err != nil {
			return nil, err
		}
	}
	mailMsg.Subject(msg.Subject)
	mailMsg.SetBodyString(msg.ContentType, msg.Body)
	return mailMsg, nil
}

// EmailConsumer is responsible for reading messages from Queue and then send these emails.
type EmailConsumer struct {
	topic     string
	group     string
	name      string
	batchSize int64

	sender *email.Sender
}

func (c *EmailConsumer) Name() string {
	return c.name
}

func (c *EmailConsumer) Topic() string {
	return c.topic
}

func (c *EmailConsumer) Group() string {
	return c.group
}

func (c *EmailConsumer) Size() int64 {
	return c.batchSize
}

func (c *EmailConsumer) Consume(ctx context.Context, id string, value any) error {
	return c.consume(ctx, id, value)
}

func (c *EmailConsumer) consume(ctx context.Context, id string, value any) error {
	val, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("mismatched value type from mq, expected map[string]any, but got %T", value)
	}

	var mqMessage string
	if val["mail"] != nil {
		mqMessage, ok = val["mail"].(string)
		if !ok {
			return fmt.Errorf("mismatched value type from mq, expected string, but got %T", mqMessage)
		}
	}

	var msg email.Message
	err := sonic.Unmarshal(str2bytes.Str2Bytes(mqMessage), &msg)
	if err != nil {
		return err
	}
	return c.sender.SendEmail(ctx, msg)
}
