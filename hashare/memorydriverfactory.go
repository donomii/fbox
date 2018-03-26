package hashconnect

import (
	"github.com/donomii/hashare"
	"github.com/donomii/vort"
)

type HashareDriverFactory struct {
	Conf     hashare.Config
	Files    map[string]*HashareFile
	Username string
	Password string
}

func (f *HashareDriverFactory) NewDriver() (d vort.FTPDriver, err error) {
	return &HashareDriver{
		Conf:     f.Conf,
		Files:    f.Files,
		Username: f.Username,
		Password: f.Password,
	}, nil
}
