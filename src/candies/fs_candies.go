package candies

import (
	"encoding/base64"
	"errtypes"
	"io/ioutil"
	"os"

	"fyne.io/fyne"
)

// FuckBase64 fucks the base64 out of the file
func FuckBase64(fn string) (string, error) {
	file, foerr := os.Open(fn)
	if foerr != nil {
		return "", foerr
	}
	cont, rderr := ioutil.ReadAll(file)
	if rderr != nil {
		return "", rderr
	}
	return base64.URLEncoding.EncodeToString(cont), nil
}

// Base64Fuck generate resource out of base64
func Base64Fuck(name string, encoded string) (*fyne.StaticResource, error) {
	res := new(fyne.StaticResource)
	res.StaticName = name
	cont, decerr := base64.URLEncoding.DecodeString(encoded)
	if decerr != nil {
		return nil, errtypes.ECORRUPT
	}
	res.StaticContent = cont
	return res, nil
}
