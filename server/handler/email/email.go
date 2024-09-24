package email

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/ginx-contribs/ginx-server/server/conf"
	"github.com/ginx-contribs/ginx-server/server/data/mq"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
	"github.com/ginx-contribs/str2bytes"
	"github.com/matcornic/hermes/v2"
	"github.com/wneessen/go-mail"
	"golang.org/x/net/context"
)

func NewEmailHandler(cfg conf.Email, client *mail.Client, queue mq.Queue) (*Handler, error) {
	handler := &Handler{Cfg: cfg, client: client, queue: queue}

	// subscribe the queue
	for _, consumer := range cfg.MQ.Consumers {
		c := &Consumer{
			topic:     cfg.MQ.Topic,
			group:     cfg.MQ.Group,
			name:      consumer,
			batchSize: cfg.MQ.BatchSize,
			h:         handler,
		}
		if err := queue.Subscribe(c); err != nil {
			return nil, err
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

// Handler is responsible for publishing and sending emails
type Handler struct {
	Cfg conf.Email

	client  *mail.Client
	queue   mq.Queue
	product hermes.Hermes
}

// SendEmail send email to smtp server
func (h *Handler) SendEmail(ctx context.Context, message Message) error {
	msg, err := BuildMail(message)
	if err != nil {
		return err
	}
	return h.client.DialAndSendWithContext(ctx, msg)
}

// PublishHermesEmail publish hermes email into queue
func (h *Handler) PublishHermesEmail(ctx context.Context, subject string, to []string, email hermes.Email) error {
	html, err := h.product.GenerateHTML(email)
	if err != nil {
		return err
	}

	msg := Message{
		ContentType: mail.TypeTextHTML,
		From:        h.Cfg.Username,
		To:          to,
		Subject:     subject,
		Body:        html,
	}

	if err := h.Publish(ctx, msg); err != nil {
		return err
	}
	return nil
}

// Publish publishes normal message to queue
func (h *Handler) Publish(ctx context.Context, msg Message) error {
	msg.From = h.Cfg.Username
	marshal, err := sonic.Marshal(msg)
	if err != nil {
		return statuserr.InternalError(err)
	}
	_, err = h.queue.Publish(ctx, h.Cfg.MQ.Topic, map[string]any{"mail": marshal}, 0)
	if err != nil {
		return statuserr.InternalError(err)
	}
	return err
}

// Consumer is responsible for reading messages from queue and then send these emails.
type Consumer struct {
	topic     string
	group     string
	name      string
	batchSize int64

	h *Handler
}

func (c *Consumer) Name() string {
	return c.name
}

func (c *Consumer) Topic() string {
	return c.topic
}

func (c *Consumer) Group() string {
	return c.group
}

func (c *Consumer) Size() int64 {
	return c.batchSize
}

func (c *Consumer) Consume(ctx context.Context, id string, value any) error {
	return c.consume(ctx, id, value)
}

func (c *Consumer) consume(ctx context.Context, id string, value any) error {
	val, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("mismatched value type from mq, expected map[string]any, but got %T", value)
	}

	var mailMsg string
	if val["mail"] != nil {
		mailMsg = val["mail"].(string)
	}

	var msg Message
	err := sonic.Unmarshal(str2bytes.Str2Bytes(mailMsg), &msg)
	if err != nil {
		return err
	}
	return c.h.SendEmail(ctx, msg)
}
