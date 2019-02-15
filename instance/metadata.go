package instance

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type MetadataParserFunc func(string) (*Metadata, error)

type LockBefore struct {
	JobName string `yaml:"job_name"`
	Release string `yaml:"release"`
}

type Metadata struct {
	BackupName                  string       `yaml:"backup_name"`
	RestoreName                 string       `yaml:"restore_name"`
	BackupShouldBeLockedBefore  []LockBefore `yaml:"backup_should_be_locked_before"`
	RestoreShouldBeLockedBefore []LockBefore `yaml:"restore_should_be_locked_before"`
}

func ParseJobMetadata(data string) (*Metadata, error) {
	metadata := &Metadata{}
	err := yaml.Unmarshal([]byte(data), metadata)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal job metadata")
	}

	for _, lockBefore := range append(metadata.BackupShouldBeLockedBefore, metadata.RestoreShouldBeLockedBefore...) {
		err = lockBefore.Validate()
		if err != nil {
			return nil, err
		}
	}

	return metadata, nil
}

func ParseJobMetadataOmitReleases(data string) (*Metadata, error) {
	metadata, err := ParseJobMetadata(data)
	if err != nil {
		return nil, err
	}

	metadata.BackupShouldBeLockedBefore = omitReleases(metadata.BackupShouldBeLockedBefore)
	metadata.RestoreShouldBeLockedBefore = omitReleases(metadata.RestoreShouldBeLockedBefore)

	return metadata, nil
}

func omitReleases(lockBefores []LockBefore) []LockBefore {
	var lockBeforesWithoutReleases []LockBefore

	for _, lockBefore := range lockBefores {
		lockBeforesWithoutReleases = append(
			lockBeforesWithoutReleases,
			LockBefore{JobName: lockBefore.JobName, Release: ""},
		)
	}

	return lockBeforesWithoutReleases
}

func (l LockBefore) Validate() error {
	if l.JobName == "" || l.Release == "" {
		return errors.New(
			"both job name and release should be specified for should be locked before")
	}
	return nil
}
