package cork

import (
	"testing"
)

func Setup(r Routes) {
}

func TestConfigure(t *testing.T) {
	c := Pop()
	c.Routes = Setup
	err := c.configure()

	if err != nil {
		t.Fail()
	}

	refute(t, c, nil)
}

func TestDefaultServices(t *testing.T) {
	c := Pop()

	refute(t, c.Services, nil)

	s := c.Services
	refute(t, s.Router, nil)
	refute(t, s.Dispatcher, nil)
	refute(t, s.Action, nil)
}
