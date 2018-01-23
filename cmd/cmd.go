package cmd

import (
	"fmt"

	"github.com/Masterminds/semver"

	"github.com/ckeyer/logrus"
	"github.com/spf13/cobra"
)

const (
	KeyMajor = "major"
	KeyMinor = "minor"
	KeyPatch = "patch"

	defaultVerFile = "VERSION"
)

var (
	debug    bool
	verFile  string
	confFile string

	rootCmd = cobra.Command{
		Use:   "wang",
		Short: "汪",
		Long:  "汪汪汪",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logrus.SetFormatter(&logrus.TextFormatter{})
			if debug {
				logrus.SetLevel(logrus.DebugLevel)
			}
			logrus.Debugf("start.")
		},
	}
)

// Execute
func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "for debug")
	// rootCmd.PersistentFlags().StringVarP(&verFile, "version-file", "f", "VERSION", "version filename")
	// rootCmd.PersistentFlags().StringVarP(&confFile, "config-file", "c", "", "config filename")

	for k, fc := range map[string]func(semver.Version) semver.Version{
		KeyMajor: func(old semver.Version) semver.Version {
			logrus.WithField("old", old.String()).Debugf("IncMajor()")
			return old.IncMajor()
		},
		KeyMinor: func(old semver.Version) semver.Version {
			logrus.WithField("old", old.String()).Debugf("IncMinor()")
			return old.IncMinor()
		},
		KeyPatch: func(old semver.Version) semver.Version {
			logrus.WithField("old", old.String()).Debugf("IncPatch()")
			return old.IncPatch()
		},
	} {
		rootCmd.AddCommand(incCmd(k, fc))
	}
}

// incCmd
func incCmd(k string, fc func(semver.Version) semver.Version) *cobra.Command {
	cmd := &cobra.Command{
		Use:   k,
		Short: fmt.Sprintf("inc %s version.", k),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				verFile = args[0]
			} else {
				verFile = defaultVerFile
			}

			v, err := readVer(verFile)
			if err != nil {
				logrus.Errorf("read version file %s failed, %s", verFile, err)
				return
			}

			next := fc(*v)

			err = writeVer(verFile, next)
			if err != nil {
				logrus.Errorf("write version file %s failed, %s", verFile, err)
				return
			}
			logrus.WithField("old", v.String()).WithField("new", next.String()).Info("inc %s successful.", k)
		},
	}
	return cmd
}
