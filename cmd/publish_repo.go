package cmd

import (
	"github.com/smira/commander"
	"github.com/smira/flag"
)

func makeCmdPublishRepo() *commander.Command {
	cmd := &commander.Command{
		Run:       aptlyPublishSnapshotOrRepo,
		UsageLine: "repo <name> [<prefix>]",
		Short:     "publish local repository",
		Long: `
Command publishes current state of local repository ready to be consumed
by apt tools. Published repostiories appear under rootDir/public directory.
Valid GPG key is required for publishing.

It is not recommended to publish local repositories directly unless the
repository is for testing purposes and changes happen frequently. For
production usage please take snapshot of repository and publish it
using publish snapshot command.

Example:

    $ aptly publish repo testing
`,
		Flag: *flag.NewFlagSet("aptly-publish-repo", flag.ExitOnError),
	}
	cmd.Flag.String("distribution", "", "distribution name to publish")
	cmd.Flag.String("component", "", "component name to publish")
	cmd.Flag.String("gpg-key", "", "GPG key ID to use when signing the release")
	cmd.Flag.String("keyring", "", "GPG keyring to use (instead of default)")
	cmd.Flag.String("secret-keyring", "", "GPG secret keyring to use (instead of default)")
	cmd.Flag.Bool("skip-signing", false, "don't sign Release files with GPG")

	return cmd
}
