// Package webroot implements a HTTP provider for solving the HTTP-01 challenge using web server's root path.
package webroot

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/xenolf/lego/acme"
)

// HTTPProviderWebroot implements ChallengeProvider for `http-01` challenge
type HTTPProviderWebroot struct {
	path string
}

// NewHTTPProviderWebroot returns a HTTPProviderWebroot instance with a configured webroot path
func NewHTTPProviderWebroot(path string) (*HTTPProviderWebroot, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("Webroot path does not exist")
	}

	c := &HTTPProviderWebroot{
		path: path,
	}

	return c, nil
}

// Present makes the token available at `HTTP01ChallengePath(token)` by creating a file in the given webroot path
func (w *HTTPProviderWebroot) Present(domain, token, keyAuth string) error {
	var err error

	challengeFilePath := path.Join(w.path, acme.HTTP01ChallengePath(token))
	err = os.MkdirAll(path.Dir(challengeFilePath), 0777)
	if err != nil {
		return fmt.Errorf("Could not create required directories in webroot for HTTP challenge -> %v", err)
	}

	err = ioutil.WriteFile(challengeFilePath, []byte(keyAuth), 0777)
	if err != nil {
		return fmt.Errorf("Could not write file in webroot for HTTP challenge -> %v", err)
	}

	return nil
}

// CleanUp removes the file created for the challenge
func (w *HTTPProviderWebroot) CleanUp(domain, token, keyAuth string) error {
	var err error
	err = os.Remove(path.Join(w.path, acme.HTTP01ChallengePath(token)))
	if err != nil {
		return fmt.Errorf("Could not remove file in webroot after HTTP challenge -> %v", err)
	}

	return nil
}