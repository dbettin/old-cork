package cork

import (
	"errors"
	"net/http"
)

type Cork struct {
	Routes   func(r Routes)
	Services *Services
	Error    Handler

	configured bool
}

func Pop() *Cork {
	cork := new(Cork)
	cork.Services = NewServices()

	return cork
}

type Services struct {
	Dispatcher     Dispatcher
	Router         Router
	Action         Action
	MessageCreator MessageCreator
}

func NewServices() *Services {
	return &Services{
		Dispatcher:     new(defaultDispatcher),
		Router:         new(defaultRouter),
		Action:         new(DefaultAction),
		MessageCreator: new(defaultMessageCreator),
	}
}

type SettingsHandler interface {
	Configure(cork *Cork)
}

func (c *Cork) configure() error {

	if c.Routes == nil {
		return errors.New("Routes setting must be configured to properly dispatch requests.")
	}

	if c.Services.Router == nil {
		return errors.New("Router setting must be configured to properly dispatch requests.")
	}
	c.Services.Router.Configure(c)

	if c.Services.Dispatcher == nil {
		return errors.New("Dispatcher setting must be configured to properly dispatch requests.")
	}
	c.Services.Dispatcher.Configure(c)

	c.configured = true
	return nil
}

func (c *Cork) Handler() http.Handler {

	if !c.configured {
		err := c.configure()
		if err != nil {
			panic("Failed to properly configure Cork settings: " + err.Error())
		}
	}

	return c.Services.Dispatcher
}
