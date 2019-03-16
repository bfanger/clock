package pubsub

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

// TheThingsNetwork connects to both MQTT and rest api
type TheThingsNetwork struct {
	urlPrefix string
	key       string
	mqtt      mqtt.Client
}

// NewTheThingsNetwork creates a new TTN connection
func NewTheThingsNetwork() (*TheThingsNetwork, error) {
	u, err := url.Parse(os.Getenv("TTN_URL"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	urlPrefix := "https://" + u.User.Username() + ".data.thethingsnetwork.org/api/v2/query?last="
	key, _ := u.User.Password()
	opts := mqtt.NewClientOptions()
	opts.SetUsername(u.User.Username())
	if p, ok := u.User.Password(); ok {
		opts.SetPassword(p)
	}
	u.User = nil
	opts.AddBroker(u.String())
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, errors.WithStack(token.Error())
	}

	return &TheThingsNetwork{urlPrefix: urlPrefix, key: key, mqtt: c}, nil
}

// LatLng a measurement from the GPS node
type LatLng struct {
	Latitude  float64   `json:"lat"`
	Longitude float64   `json:"lng"`
	Precision int       `json:"p"`
	Time      time.Time `json:"time"`
}

// History loads historic data from the Data Storage intergration API.
func (t *TheThingsNetwork) History(last string) ([]LatLng, error) {
	if last == "" {
		last = "1h"
	}
	r, err := http.NewRequest("GET", t.urlPrefix+last, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	r.Header.Set("Authorization", "key "+t.key)
	response, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var records []LatLng
	json.Unmarshal(data, &records)
	sort.Slice(records, func(i, j int) bool {
		return records[i].Time.After(records[j].Time)
	})
	return records, nil
}

// Updates from nodes
func (t *TheThingsNetwork) Updates() chan *LatLng {
	messages := make(chan *LatLng)
	type mqttMessage struct {
		LatLng   *LatLng `json:"payload_fields"`
		Metadata struct {
			Time time.Time
		}
	}
	t.mqtt.Subscribe("+/devices/+/up", 0, func(c mqtt.Client, m mqtt.Message) {
		var decoded mqttMessage
		if err := json.Unmarshal(m.Payload(), &decoded); err != nil {
			panic(err)
		}
		decoded.LatLng.Time = decoded.Metadata.Time
		messages <- decoded.LatLng
	})
	return messages
}
