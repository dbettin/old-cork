package testing

import (
	"errors"
	"github.com/dbettin/cork"
	"net/url"
)

var Settings *cork.Cork

func Get(route string) (*cork.Route, error) {
	if err := checkSettings(); err != nil {
		return nil, err
	}

	_ = Settings.Handler()

	url, err := url.Parse(route)

	if err != nil {
		return nil, err
	}

	foundRoute := Settings.Services.Router.Route(cork.GET, url)

	return foundRoute, nil
}

func checkSettings() error {
	if Settings == nil {
		return errors.New("Please add Cork settings before running route tests")
	}

	return nil
}
