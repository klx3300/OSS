package staticres

import (
	"errtypes"
	"io/ioutil"
	"logger"
	"os"
)

// Resource represents static resource the program is intended to access
type Resource struct {
	Name    string
	Content []byte
}

var log = logger.Logger{
	LogLevel: 0,
	Name:     "StaticResourceLoader:",
}

// ResSearchPath will give hints on where to search resources
var ResSearchPath = []string{
	"/Users/imzhwk/HUST/DatabaseSystem",
}

// GetResource reads resource from disk
func GetResource(name string, filename string) (Resource, error) {
	log.Debugln("Loading resource<"+name+">", filename)
	for _, respath := range ResSearchPath {
		tmpPath := respath + string(os.PathSeparator) + filename
		log.Debugln("Trying path", tmpPath)
		resfile, reserr := os.Open(tmpPath)
		if reserr != nil {
			log.Debugln("Failed for path", tmpPath)
			continue
		}
		log.Debugln("Found for path", tmpPath)
		cont, rdallerr := ioutil.ReadAll(resfile)
		if rdallerr != nil {
			log.Debugln("Read failed for path", tmpPath)
			continue
		}
		return Resource{
			Name:    name,
			Content: cont,
		}, nil
	}
	log.Debugln("Trying program execution path.")
	resfile, reserr := os.Open(filename)
	if reserr != nil {
		log.Debugln("Failed for execution path")
		return Resource{}, errtypes.ENEXIST
	}
	log.Debugln("Found for execution path.")
	cont, rdallerr := ioutil.ReadAll(resfile)
	if rdallerr != nil {
		log.Debugln("Read failed for exec path")
		return Resource{}, errtypes.ENEXIST
	}
	return Resource{
		Name:    name,
		Content: cont,
	}, nil
}
