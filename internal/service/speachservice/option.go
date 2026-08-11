package speachservice

import (
	"time"

	"github.com/Piktet/tg_bot/internal/repository/speach"
)

const defaultStatusTimeout = time.Second

type speachOption struct {
	connSpeach    *speach.SpeachConnection
	host          string
	statusTimeout time.Duration
}

type Option func(t *speachOption)

func WithSpeach(connSpeach *speach.SpeachConnection) Option {
	return func(x *speachOption) {
		x.connSpeach = connSpeach
	}
}

func WithHost(host string) Option {
	return func(x *speachOption) {
		x.host = host
	}
}

func WithStatusTimeout(t time.Duration) Option {
	return func(x *speachOption) {
		x.statusTimeout = t
	}
}

func newSpeachOption(opts ...Option) *speachOption {
	o := &speachOption{
		statusTimeout: defaultStatusTimeout,
	}
	for _, v := range opts {
		v(o)
	}
	return o
}
