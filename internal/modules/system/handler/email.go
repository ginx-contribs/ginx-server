package handler

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/ginx-contribs/ginx-server/internal/conf"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/types"
	"github.com/ginx-contribs/ginx-server/pkg/mq"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
	"github.com/ginx-contribs/str2bytes"
	"github.com/matcornic/hermes/v2"
	"github.com/wneessen/go-mail"
	"golang.org/x/net/context"
)

func NewEmailHandler(cfg conf.Email, client *mail.Client, queue mq.Queue) (EmailHandler, error) {
	handler := EmailHandler{Config: cfg, Email: client, Queue: queue}

	// subscribe the Queue
	for _, consumer := range cfg.MQ.Consumers {
		c := &EmailConsumer{
			topic:     cfg.MQ.Topic,
			group:     cfg.MQ.Group,
			name:      consumer,
			batchSize: cfg.MQ.BatchSize,
			h:         handler,
		}
		if err := queue.Subscribe(c); err != nil {
			return handler, err
		}
	}

	handler.product = hermes.Hermes{
		Product: hermes.Product{
			Name:      "ginx-contribs",
			Copyright: "Copyright © ginx-contribs",
		},
	}
	return handler, nil
}

// EmailHandler is responsible for publishing and sending emails
type EmailHandler struct {
	Config conf.Email
	Email  *mail.Client

	Queue   mq.Queue
	product hermes.Hermes
}

// Send send email to smtp server
func (e *EmailHandler) Send(ctx context.Context, message types.EmailBody) error {
	msg, err := e.buildMail(message)
	if err != nil {
		return err
	}
	return e.Email.DialAndSendWithContext(ctx, msg)
}

// PublishEmail publish email into Queue
func (e *EmailHandler) PublishEmail(ctx context.Context, subject string, to []string, email hermes.Email) error {
	html, err := e.product.GenerateHTML(email)
	if err != nil {
		return err
	}

	msg := types.EmailBody{
		ContentType: mail.TypeTextHTML,
		From:        e.Config.Username,
		To:          to,
		Subject:     subject,
		Body:        html,
	}

	if err := e.Publish(ctx, msg); err != nil {
		return err
	}
	return nil
}

// Publish publishes normal message to Queue
func (e *EmailHandler) Publish(ctx context.Context, msg types.EmailBody) error {
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

	h EmailHandler
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

	var mailMsg string
	if val["mail"] != nil {
		mailMsg = val["mail"].(string)
	}

	var msg types.EmailBody
	err := sonic.Unmarshal(str2bytes.Str2Bytes(mailMsg), &msg)
	if err != nil {
		return err
	}
	return c.h.Send(ctx, msg)
}
