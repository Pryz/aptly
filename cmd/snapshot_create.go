package cmd

import (
	"fmt"
	"github.com/smira/aptly/debian"
	"github.com/smira/commander"
	"github.com/smira/flag"
)

func aptlySnapshotCreate(cmd *commander.Command, args []string) error {
	var (
		err      error
		snapshot *debian.Snapshot
	)

	if len(args) == 4 && args[1] == "from" && args[2] == "mirror" {
		// aptly snapshot create snap from mirror mirror
		var repo *debian.RemoteRepo

		repoName, snapshotName := args[3], args[0]

		repo, err = context.collectionFactory.RemoteRepoCollection().ByName(repoName)
		if err != nil {
			return fmt.Errorf("unable to create snapshot: %s", err)
		}

		err = context.collectionFactory.RemoteRepoCollection().LoadComplete(repo)
		if err != nil {
			return fmt.Errorf("unable to create snapshot: %s", err)
		}

		snapshot, err = debian.NewSnapshotFromRepository(snapshotName, repo)
		if err != nil {
			return fmt.Errorf("unable to create snapshot: %s", err)
		}
	} else if len(args) == 4 && args[1] == "from" && args[2] == "repo" {
		// aptly snapshot create snap from repo repo
		var repo *debian.LocalRepo

		localRepoName, snapshotName := args[3], args[0]

		repo, err = context.collectionFactory.LocalRepoCollection().ByName(localRepoName)
		if err != nil {
			return fmt.Errorf("unable to create snapshot: %s", err)
		}

		err = context.collectionFactory.LocalRepoCollection().LoadComplete(repo)
		if err != nil {
			return fmt.Errorf("unable to create snapshot: %s", err)
		}

		snapshot, err = debian.NewSnapshotFromLocalRepo(snapshotName, repo)
		if err != nil {
			return fmt.Errorf("unable to create snapshot: %s", err)
		}
	} else if len(args) == 2 && args[1] == "empty" {
		// aptly snapshot create snap empty
		snapshotName := args[0]

		packageList := debian.NewPackageList()

		snapshot = debian.NewSnapshotFromPackageList(snapshotName, nil, packageList, "Created as empty")
	} else {
		cmd.Usage()
		return err
	}

	err = context.collectionFactory.SnapshotCollection().Add(snapshot)
	if err != nil {
		return fmt.Errorf("unable to add snapshot: %s", err)
	}

	fmt.Printf("\nSnapshot %s successfully created.\nYou can run 'aptly publish snapshot %s' to publish snapshot as Debian repository.\n", snapshot.Name, snapshot.Name)

	return err
}

func makeCmdSnapshotCreate() *commander.Command {
	cmd := &commander.Command{
		Run:       aptlySnapshotCreate,
		UsageLine: "create <name> from mirror <mirror-name> | from repo <repo-name> | empty",
		Short:     "creates snapshot of mirror (local repository) contents",
		Long: `
Command create <name> from mirror makes persistent immutable snapshot of remote
repository mirror. Snapshot could be published or further modified using
merge, pull and other aptly features.

Command create <name> from repo makes persistent immutable snapshot of local
repository. Snapshot could be processed as mirror snapshots, and mixed with
snapshots of remote mirrors.

Command create <name> empty creates empty snapshot that could be used as a
basis for snapshot pull operations, for example. As snapshots are immutable,
creating one empty snapshot should be enough.

Example:

  $ aptly snapshot create wheezy-main-today from mirror wheezy-main
`,
		Flag: *flag.NewFlagSet("aptly-snapshot-create", flag.ExitOnError),
	}

	return cmd

}
