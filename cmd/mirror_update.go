package cmd

import (
	"fmt"
	"github.com/smira/commander"
	"github.com/smira/flag"
)

func aptlyMirrorUpdate(cmd *commander.Command, args []string) error {
	var err error
	if len(args) != 1 {
		cmd.Usage()
		return err
	}

	name := args[0]

	repo, err := context.collectionFactory.RemoteRepoCollection().ByName(name)
	if err != nil {
		return fmt.Errorf("unable to update: %s", err)
	}

	err = context.collectionFactory.RemoteRepoCollection().LoadComplete(repo)
	if err != nil {
		return fmt.Errorf("unable to update: %s", err)
	}

	ignoreMismatch := cmd.Flag.Lookup("ignore-checksums").Value.Get().(bool)

	verifier, err := getVerifier(cmd)
	if err != nil {
		return fmt.Errorf("unable to initialize GPG verifier: %s", err)
	}

	err = repo.Fetch(context.downloader, verifier)
	if err != nil {
		return fmt.Errorf("unable to update: %s", err)
	}

	err = repo.Download(context.progress, context.downloader, context.collectionFactory, context.packagePool, ignoreMismatch)
	if err != nil {
		return fmt.Errorf("unable to update: %s", err)
	}

	err = context.collectionFactory.RemoteRepoCollection().Update(repo)
	if err != nil {
		return fmt.Errorf("unable to update: %s", err)
	}

	context.progress.Printf("\nMirror `%s` has been successfully updated.\n", repo.Name)
	return err
}

func makeCmdMirrorUpdate() *commander.Command {
	cmd := &commander.Command{
		Run:       aptlyMirrorUpdate,
		UsageLine: "update <name>",
		Short:     "update mirror",
		Long: `
Updates remote mirror (downloads package files and meta information). When mirror is created,
this command should be run for the first time to fetch mirror contents. This command could be
run many times to get updated repository contents. If interrupted, command could be restarted safely.

Example:

  $ aptly mirror update wheezy-main
`,
		Flag: *flag.NewFlagSet("aptly-mirror-update", flag.ExitOnError),
	}

	cmd.Flag.Bool("ignore-checksums", false, "ignore checksum mismatches while downloading package files and metadata")
	cmd.Flag.Bool("ignore-signatures", false, "disable verification of Release file signatures")
	cmd.Flag.Var(&keyRings, "keyring", "gpg keyring to use when verifying Release file (could be specified multiple times)")

	return cmd
}
