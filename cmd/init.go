package cmd

import (
	"os"
	"text/template"

	"github.com/ckeyer/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(InitCmd())
}

func InitCmd() *cobra.Command {
	var (
		force       bool
		versionFile string
		remote      string
		branch      string
	)
	cmd := &cobra.Command{
		Use: "init",
		Run: func(cmd *cobra.Command, args []string) {
			_, err := os.Open(versionFile)
			if err != nil {
				logrus.Errorf("open version file failed, %s", err)
				return
			}

			hookfile := ".git/hooks/pre-push"
			if force {
				os.RemoveAll(hookfile)
			}
			f, err := os.OpenFile(hookfile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0755)
			if err != nil {
				logrus.Errorf("create file %s failed, %s", hookfile, err)
				return
			}
			defer f.Close()

			data := map[string]interface{}{
				"VersionFile": versionFile,
				"Remote":      remote,
				"Branch":      branch,
			}

			tpl, err := template.New("default").Parse(prePushTpl)
			if err != nil {
				panic(err)
				return
			}

			err = tpl.Execute(f, data)
			if err != nil {
				logrus.Errorf("write hook file failed, %s", err)
				return
			}
			logrus.Info("init successful.")
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "remove old hooks.")
	cmd.Flags().StringVarP(&versionFile, "version-file", "v", "VERSION", "version file.")
	cmd.Flags().StringVarP(&remote, "remote", "r", "origin", "git remote name")
	cmd.Flags().StringVarP(&branch, "branch", "b", "master", "git branch")
	return cmd
}

const (
	prePushTpl = `#!/bin/sh

remote="$1"

while read local_ref local_sha remote_ref remote_sha
do
	if [[ $remote_ref = 'refs/heads/{{.Branch}}' ]] && [[ $remote = '{{.Remote}}' ]]; then
		wang patch {{.VersionFile}};
		git add {{.VersionFile}}
		git commit -m "upgrade version"
	fi
done

exit 0
`
)
