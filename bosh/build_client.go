package bosh

import (
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/instance"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/ssh"
	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/pkg/errors"

	boshuaa "github.com/cloudfoundry/bosh-cli/uaa"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

func BuildClient(targetUrl, username, password, caCert string, logger boshlog.Logger) (BoshClient, error) {
	config, err := director.NewConfigFromURL(targetUrl)
	if err != nil {
		return nil, errors.Errorf("invalid bosh URL - %s", err.Error())
	}

	config.CACert = caCert

	directorFactory := director.NewFactory(logger)

	info, err := getDirectorInfo(directorFactory, config)
	if err != nil {
		return nil, err
	}

	if info.Auth.Type == "uaa" {
		uaa, err := buildUaa(info, username, password, caCert, logger)
		if err != nil {
			return nil, err
		}

		config.TokenFunc = boshuaa.NewClientTokenSession(uaa).TokenFunc
	} else {
		config.Client = username
		config.ClientSecret = password
	}

	boshDirector, err := directorFactory.New(config, director.NewNoopTaskReporter(), director.NewNoopFileReporter())
	if err != nil {
		return nil, errors.Wrap(err, "error building bosh director client")
	}

	return NewClient(boshDirector, director.NewSSHOpts, ssh.NewSshRemoteRunner, logger, instance.NewJobFinder(logger), NewBoshManifestReleaseMapping), nil
}

func getDirectorInfo(directorFactory director.Factory, config director.FactoryConfig) (director.Info, error) {
	infoDirector, err := directorFactory.New(config, director.NewNoopTaskReporter(), director.NewNoopFileReporter())
	if err != nil {
		return director.Info{}, errors.Wrap(err, "error building bosh director client")
	}

	info, err := infoDirector.Info()
	if err != nil {
		return director.Info{}, errors.Wrap(err, "bosh director unreachable or unhealthy")
	}

	return info, nil
}

func buildUaa(info director.Info, username, password, cert string, logger boshlog.Logger) (boshuaa.UAA, error) {
	urlAsInterface := info.Auth.Options["url"]
	url, ok := urlAsInterface.(string)
	if !ok {
		return nil, errors.Errorf("Expected URL '%s' to be a string", urlAsInterface)
	}

	uaaConfig, err := boshuaa.NewConfigFromURL(url)
	if err != nil {
		return nil, errors.Wrap(err, "invalid UAA URL")
	}

	uaaConfig.CACert = cert
	uaaConfig.Client = username
	uaaConfig.ClientSecret = password

	return boshuaa.NewFactory(logger).New(uaaConfig)
}
