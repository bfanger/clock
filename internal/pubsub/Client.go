package pubsub

import (
	"log"
	"net/url"
	"regexp"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

const device = "clock"

// Connection to pubsub
type Connection struct {
	subscriptions []chan mqtt.Message
	mqtt          mqtt.Client
}

// NewConnection to the MQTT server
func NewConnection(uri string) (*Connection, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	opts := mqtt.NewClientOptions()
	opts.SetUsername(u.User.Username())
	if p, ok := u.User.Password(); ok {
		opts.SetPassword(p)
	}
	u.User = nil
	opts.AddBroker(u.String())
	opts.SetWill("uptime/"+device, "", 0, true)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	t := c.Publish("uptime/"+device, 0, true, strconv.Itoa(int(time.Now().Unix())))
	if t.Wait() && t.Error() != nil {
		return nil, t.Error()
	}

	return &Connection{mqtt: c}, nil
}

// Close the connection
func (c *Connection) Close() error {
	for _, s := range c.subscriptions {
		close(s)
	}
	c.subscriptions = nil
	t := c.mqtt.Publish("uptime/"+device, 0, true, "")
	t.Wait()
	err := t.Error()
	if err != nil {
		return err
	}
	c.mqtt.Disconnect(250)
	return nil
}

// Subscribe to a topic
func (c *Connection) Subscribe(topic string) chan mqtt.Message {
	messages := make(chan mqtt.Message)
	c.mqtt.Subscribe(topic, 0, func(c mqtt.Client, m mqtt.Message) {
		log.Println(m.Topic(), string(m.Payload()))
		messages <- m
	})
	c.subscriptions = append(c.subscriptions, messages)
	return messages
}

// Publish to a topic
func (c *Connection) Publish(topic string, payload interface{}) error {
	t := c.mqtt.Publish(topic, 0, false, payload)
	t.Wait()
	return errors.WithStack(t.Error())
}

var extractID = regexp.MustCompile("[^\\/]+$")

// HandleRPC respond to a rpc method
func (c *Connection) HandleRPC(method string, handler func(options []byte) []byte) {
	for m := range c.Subscribe("rpc/request/" + method + "/+") {
		id := extractID.FindString((m.Topic()))
		if id == "" {
			continue
		}
		payload := handler(m.Payload())
		if err := c.Publish("rpc/response/"+method+"/"+id, payload); err != nil {
			panic(err)
		}
	}
}
