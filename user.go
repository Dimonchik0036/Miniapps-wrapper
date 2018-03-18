package mapps

type User struct {
	Subscriber string `json:"subscriber,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
	Service    string `json:"service,omitempty"`
}

func (u *User) Key() string {
	return u.Protocol + "=" + u.Subscriber
}

func (u *User) SendMessage(text string) error {
	if u.Protocol == Telegram {
		return u.SendMessageTelegram(text)
	}
	return u.SendMessageDefault(text)
}

func (u *User) SendMessageDefault(text string) error {
	return Push(PushConfig{
		Service:    u.Service,
		Protocol:   u.Protocol,
		Subscriber: u.Subscriber,
		Document:   OnlyTextDefault(text),
	})
}

func (u *User) SendMessageTelegram(text string) error {
	return Push(PushConfig{
		Service:    u.Service,
		Protocol:   u.Protocol,
		Subscriber: u.Subscriber,
		Document:   OnlyTextTelegram(text),
	})
}

func (u *User) SendMessageBlock(div string) error {
	return Push(PushConfig{
		Service:    u.Service,
		Protocol:   u.Protocol,
		Subscriber: u.Subscriber,
		Document:   Page("", div),
	})
}
