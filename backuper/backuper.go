package backuper

import (
	"fmt"
	"os"
)

func New(bosh BoshDirector, artifactCreator ArtifactCreator) *Backuper {
	return &Backuper{
		BoshDirector:    bosh,
		ArtifactCreator: artifactCreator,
	}
}

//go:generate counterfeiter -o fakes/fake_artifact_creator.go . ArtifactCreator
type ArtifactCreator func(string) (Artifact, error)

//go:generate counterfeiter -o fakes/fake_artifact.go . Artifact
type Artifact interface {
	CreateFile()
}

type Backuper struct {
	BoshDirector
	ArtifactCreator
}

//Backup checks if a deployment has backupable instances and backs them up.
func (b Backuper) Backup(deploymentName string) error {
	instances, err := b.FindInstances(deploymentName)
	if err != nil {
		return err
	}
	defer instances.Cleanup()

	backupable, err := instances.AllBackupable()
	if err != nil {
		return err
	}
	if len(backupable) == 0 {
		return fmt.Errorf("Deployment '%s' has no backup scripts", deploymentName)
	}

	b.ArtifactCreator(deploymentName)

	err = os.MkdirAll(deploymentName, 0700)
	if err != nil {
		panic("oh my christ")
	}

	return backupable.Backup()
}

//go:generate counterfeiter -o fakes/fake_bosh_director.go . BoshDirector
type BoshDirector interface {
	FindInstances(deploymentName string) (Instances, error)
}
