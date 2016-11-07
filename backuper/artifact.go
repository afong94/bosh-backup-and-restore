package backuper

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

func DirectoryArtifactCreator(name string) (Artifact, error) {
	return &DirectoryArtifact{baseDirName: name}, os.MkdirAll(name, 0700)
}

type DirectoryArtifact struct {
	baseDirName string
}

type InstanceMetadata struct {
	InstanceName string `yaml:"instance_name"`
	InstanceID   string `yaml:"instance_id"`
	Checksum     string `yaml:"checksum"`
}

type metadata struct {
	MetadataForEachInstance []InstanceMetadata `yaml:"instances"`
}

func (d *DirectoryArtifact) CreateFile(inst Instance) (io.WriteCloser, error) {
	filename := inst.Name() + "-" + inst.ID() + ".tgz"
	return os.Create(path.Join(d.baseDirName, filename))
}

func (d *DirectoryArtifact) CalculateChecksum(inst Instance) (string, error) {
	sha := sha1.New()
	filename := d.instanceFilename(inst)
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(sha, file); err != nil {
		return "", err
	}
	checksum := sha.Sum(nil)
	return fmt.Sprintf("%x", checksum), nil
}

func (d *DirectoryArtifact) AddChecksum(inst Instance, shasum string) error {
	metadata, err := d.readMetadata()
	if err != nil {
		return err
	}

	metadata.MetadataForEachInstance = append(metadata.MetadataForEachInstance, InstanceMetadata{
		InstanceName: inst.Name(),
		InstanceID:   inst.ID(),
		Checksum:     shasum,
	})

	return d.saveMetadata(metadata)
}

func (d *DirectoryArtifact) saveMetadata(data metadata) error {
	contents, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(d.metadataFilename(), contents, 0666)
}

func (d *DirectoryArtifact) instanceFilename(inst Instance) string {
	return path.Join(d.baseDirName, inst.Name()+"-"+inst.ID()+".tgz")
}

func (d *DirectoryArtifact) metadataFilename() string {
	return path.Join(d.baseDirName, "metadata")
}

func (d *DirectoryArtifact) readMetadata() (metadata, error) {
	metadata := metadata{}

	fileInfo, _ := os.Stat(d.metadataFilename())
	if fileInfo != nil {
		contents, err := ioutil.ReadFile(d.metadataFilename())

		if err != nil {
			return metadata, err
		}

		if err := yaml.Unmarshal(contents, &metadata); err != nil {
			return metadata, err
		}
	}
	return metadata, nil
}
