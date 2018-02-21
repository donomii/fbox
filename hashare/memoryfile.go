package hashconnect

import (
	"os"
)

type HashareFile struct {
	File    os.FileInfo
	Content []byte
}
