package mapps

import (
	"errors"
	"github.com/valyala/fasthttp"
	"log"
	"net/url"
	"strconv"
	"time"
)

type PushConfig struct {
	Subscriber string
	Service    string
	Document   string
	Scenario   string
	Protocol   string
	ResourceId string
	Delay      int64
}

func (c *PushConfig) Values() url.Values {
	values := url.Values{}
	values.Set("protocol", c.Protocol)
	values.Set("service", c.Service)
	values.Set("scenario", func() string {
		if c.Scenario == "" {
			return XmlPush
		}
		return c.Scenario
	}())
	values.Set("subscriber", c.Subscriber)
	values.Set("document", c.Document)
	return values
}

const (
	apiScheme = "https"
	apiPath   = "push"
	apiHost   = "api.miniapps.run"
)

func Push(config PushConfig) error {
	u := url.URL{
		Scheme:   apiScheme,
		Path:     apiPath,
		Host:     apiHost,
		RawQuery: config.Values().Encode(),
	}

	code, _, err := fasthttp.Get(nil, u.String())
	if err != nil {
		return err
	}

	if code != fasthttp.StatusOK {
		return errors.New("Status code: " + strconv.Itoa(code))
	}

	return nil
}

type Request struct {
	Ctx        *fasthttp.RequestCtx `json:"-"`
	RequestUrl string               `json:"request_url,omitempty"`
	Data       url.Values           `json:"data"`
	Page       string               `json:"page,omitempty"`
	Protocol   string               `json:"protocol,omitempty"`
	Subscriber string               `json:"subscriber,omitempty"`
	UserId     string               `json:"user_id,omitempty"`
	Service    string               `json:"service,omitempty"`
	Lang       string               `json:"lang,omitempty"`
	Abonent    string               `json:"abonent,omitempty"`
	ServiceId  string               `json:"service_id,omitempty"`
	Scenario   string               `json:"scenario,omitempty"`
	BadCommand string               `json:"bad_command,omitempty"`
	Date       int64                `json:"date,omitempty"`
	InputType  string               `json:"input_type,omitempty"`
	Event
}

func (r *Request) String() string {
	return r.Page + " " + r.Event.String()
}

func (r *Request) AllFields() string {
	return "page=\"" + r.Page +
		"\"&protocol=\"" + r.Protocol +
		"\"&subscriber=\"" + r.Subscriber +
		"\"&user_id=\"" + r.UserId +
		"\"&service=\"" + r.Service +
		"\"&lang=\"" + r.Lang +
		"\"&abonent=\"" + r.Abonent +
		"\"&service_id=\"" + r.ServiceId +
		"\"&scenario=\"" + r.Scenario +
		"\"&input_type=\"" + r.InputType +
		"\"&bad_command=\"" + r.BadCommand +
		"\"&date=\"" + strconv.FormatInt(r.Date, 10) +
		"\"&request=\"" + r.RequestUrl + "\""
}

func (r *Request) User() User {
	return User{
		Protocol:   r.Protocol,
		Subscriber: r.Subscriber,
		Service:    r.Service,
	}
}

func (r *Request) GetField(key string) string {
	if r.Data == nil {
		return ""
	}
	return Unescaped(r.Data.Get(key))
}

func (r *Request) SetField(key string, value string) {
	if r.Data == nil {
		r.Data = url.Values{}
	}

	r.Data.Set(key, url.QueryEscape(value))
}

const (
	EventLink    = "link"
	EventMessage = "message"
	EventPush    = "push"

	EventTypeFile = "file" //message
	EventTypeText = "text" //message
	EventTypeHttp = "http" //push

	EventMediaTypePhoto = "photo" //file
)

type Event struct {
	Event     string `json:"event,omitempty"`
	Type      string `json:"event.type,omitempty"`
	Text      string `json:"event.text,omitempty"`
	Order     int64  `json:"event.order,omitempty"`
	Url       string `json:"event.url,omitempty"`
	Referer   string `json:"event.referer,omitempty"`
	MediaType string `json:"event.media_type,omitempty"`
	Source    string `json:"event.source,omitempty"`
	Size      int64  `json:"size,omitempty"`
}

func (e *Event) AllFields() string {
	return "event=\"" + e.Event +
		"\"&type=\"" + e.Type +
		"\"&text=\"" + e.Text +
		"\"&order=\"" + strconv.FormatInt(e.Order, 10) +
		"\"&url=\"" + e.Url +
		"\"&referer=\"" + e.Referer +
		"\"&media_type=\"" + e.MediaType +
		"\"&size=\"" + strconv.FormatInt(e.Size, 10) +
		"\"&source=\"" + e.Source + "\""
}

func (e *Event) String() string {
	switch e.Event {
	case EventPush, EventLink, EventMessage:
		return e.AllFields()
	case "":
		return "Empty event"
	default:
		log.Print("New diffirent event!")
		return e.AllFields()
	}

	return e.Event
}

func Decode(s string) (Request, error) {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return Request{}, err
	}

	var r Request
	r.RequestUrl = u.String()
	r.Page = u.Path[1:]

	r.Data = u.Query()
	r.Protocol = Unescaped(r.Data.Get("protocol"))
	r.Subscriber = Unescaped(r.Data.Get("subscriber"))
	r.UserId = Unescaped(r.Data.Get("user_id"))
	r.Service = Unescaped(r.Data.Get("service"))
	r.Lang = Unescaped(r.Data.Get("lang"))
	r.Abonent = Unescaped(r.Data.Get("abonent"))
	r.ServiceId = Unescaped(r.Data.Get("serviceId"))
	r.BadCommand = Unescaped(r.Data.Get("bad_command"))
	r.Date = time.Now().Unix()
	r.Scenario = Unescaped(r.Data.Get("scenario"))
	r.InputType = Unescaped(r.Data.Get("input_type"))
	r.Event.Event = Unescaped(r.Data.Get("event"))
	r.Event.Type = Unescaped(r.Data.Get("event.type"))
	r.Event.Text = Unescaped(r.Data.Get("event.text"))
	r.Event.Url = Unescaped(r.Data.Get("event.url"))
	r.Event.Referer = Unescaped(r.Data.Get("event.referer"))
	r.Event.MediaType = Unescaped(r.Data.Get("event.media_type"))
	r.Event.Source = Unescaped(r.Data.Get("event.source"))
	r.Event.Order, _ = strconv.ParseInt(Unescaped(r.Data.Get("event.order")), 10, 64)
	r.Event.Size, _ = strconv.ParseInt(Unescaped(r.Data.Get("event.size")), 10, 64)

	return r, nil
}

func Unescaped(s string) string {
	u, err := url.QueryUnescape(s)
	if err != nil {
		log.Print(err)
		return s
	}
	return u
}
