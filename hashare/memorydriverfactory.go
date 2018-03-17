package hashconnect

import (
	"github.com/donomii/fbox"
	"github.com/donomii/hashare"
)

type HashareDriverFactory struct {
	Conf     hashare.Config
	Files    map[string]*HashareFile
	Username string
	Password string
}

func (f *HashareDriverFactory) NewDriver() (d fbox.FTPDriver, err error) {
	return &HashareDriver{
		Conf:     f.Conf,
		Files:    f.Files,
		Username: f.Username,
		Password: f.Password,
	}, nil
}
