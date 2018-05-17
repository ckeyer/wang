package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/ckeyer/logrus"
)

// readVer
func readVer(name string) (*semver.Version, error) {
	bs, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	ver, err := semver.NewVersion(strings.TrimSpace(string(bs)))
	if err != nil {
		logrus.Debugf("parse version %s failed, %s", bs, err)
		return nil, err
	}

	return ver, nil
}

// writeVer
func writeVer(name string, ver semver.Version) error {
	f, err := os.OpenFile(name, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(ver.String() + "\n")
	if err != nil {
		return err
	}
	return nil
}

// getVerFile
func getVerFile(root string) string {
	var verf string
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		logrus.WithField("path", path).WithField("info.name", info.Name()).Debug(err)
		if strings.ToLower(info.Name()) == "version" {
			verf = info.Name()
		}
		return nil
	})
	return verf
}
