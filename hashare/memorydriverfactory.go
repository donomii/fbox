package hashconnect

import (
	"github.com/donomii/fbox"
)

type HashareDriverFactory struct {
	Files    map[string]*HashareFile
	Username string
	Password string
}

func (f *HashareDriverFactory) NewDriver() (d fbox.FTPDriver, err error) {
	return &HashareDriver{
		Files:    f.Files,
		Username: f.Username,
		Password: f.Password,
	}, nil
}
