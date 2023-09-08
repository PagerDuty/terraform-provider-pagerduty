package persistentconfig

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"gopkg.in/ini.v1"
)

const (
	DefaultConfigFolder        = ".pagerduty"
	DefaultConfigProfileTag    = "profile"
	DefaultConfigProfile       = "default"
	DefaultConfigFileName      = "config"
	DefaultCredentialsFileName = "credentials"
)

// AppOauthScopedTokenParams parameters for setting up API calls authentication
// using App Scoped Oauth Token
type AppOauthScopedTokenParams struct {
	ClientID     string
	ClientSecret string
	PDSubDomain  string
	Region       string
	Token        string // App Oauth Scoped Token
}

type ClientPersistentConfig struct {
	AppOauthScopedTokenParams
	Profile         string
	Fs              afero.Fs
	configFile      string
	credentialsFile string
}

func (c *ClientPersistentConfig) Load() error {
	err := c.ensureConfigFiles()
	if err != nil {
		return err
	}

	err = c.SetActiveProfile(DefaultConfigProfile)
	if err != nil {
		return err
	}
	c.Profile = DefaultConfigProfile

	token, err := c.GetCredential("token")
	if err != nil {
		return err
	}

	c.Token = token

	return nil
}

func (c *ClientPersistentConfig) ensureConfigFiles() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	pagerDutyDir := filepath.Join(homeDir, DefaultConfigFolder)
	c.Fs.MkdirAll(pagerDutyDir, 0755)

	configFile := filepath.Join(pagerDutyDir, DefaultConfigFileName)
	credentialsFile := filepath.Join(pagerDutyDir, DefaultCredentialsFileName)

	for _, file := range []string{configFile, credentialsFile} {
		exists, err := afero.Exists(c.Fs, file)
		if err != nil {
			return fmt.Errorf("%w; error: persistent configuration %q file could not be created", err, file)
		}
		if !exists {
			_, err = c.Fs.Create(file)
			if err != nil {
				return fmt.Errorf("%w; error: persistent configuration %q file could not be created", err, file)
			}
			err = c.setCredentialsPermissions(file)
			if err != nil {
				return err
			}
		}
	}

	c.configFile = configFile
	c.credentialsFile = credentialsFile

	return nil
}

func (c ClientPersistentConfig) ReadConfigFile() (*ini.File, error) {
	return readConfig(c.Fs, c.configFile)
}

func (c ClientPersistentConfig) ReadCredentialsFile() (*ini.File, error) {
	return readConfig(c.Fs, c.credentialsFile)
}

func (c ClientPersistentConfig) WriteConfigFile(cfg *ini.File) error {
	return writeConfig(c.Fs, c.configFile, cfg)
}

func (c ClientPersistentConfig) WriteCredentialsFile(cfg *ini.File) error {
	return writeConfig(c.Fs, c.credentialsFile, cfg)
}

func (c ClientPersistentConfig) SetActiveProfile(profile string) error {
	cfg, err := c.ReadConfigFile()
	if err != nil {
		return err
	}

	cfg.Section(fmt.Sprintf("%s %s", DefaultConfigProfileTag, profile))
	return c.WriteConfigFile(cfg)
}

func (c ClientPersistentConfig) GetCredential(k string) (string, error) {
	cfg, err := c.ReadCredentialsFile()
	if err != nil {
		return "", err
	}

	section := cfg.Section(c.Profile)
	return section.Key(k).String(), nil
}

func (c *ClientPersistentConfig) SetCredential(k, v string) error {
	cfg, err := c.ReadCredentialsFile()
	if err != nil {
		return err
	}

	section := cfg.Section(c.Profile)
	section.Key(k).SetValue(v)
	return c.WriteCredentialsFile(cfg)
}

func (c *ClientPersistentConfig) setCredentialsPermissions(credentialsFile string) error {
	if err := c.Fs.Chmod(credentialsFile, 0600); err != nil {
		return fmt.Errorf("%w; error: file permissions could not be set to %q file", err, credentialsFile)
	}
	return nil
}

func readConfig(fs afero.Fs, file string) (*ini.File, error) {
	configFile, err := afero.ReadFile(fs, file)
	if err != nil {
		return nil, fmt.Errorf("%w; error: persistent configuration could not be read from %q file", err, file)
	}

	cfg, err := ini.Load(configFile)
	if err != nil {
		return nil, fmt.Errorf("%w; error: persistent configuration could not be read from %q file", err, file)
	}
	return cfg, nil
}

func writeConfig(fs afero.Fs, file string, cfg *ini.File) error {
	// Use a buffer to write the *ini.File content to a byte slice
	var buffer bytes.Buffer
	_, err := cfg.WriteTo(&buffer)
	if err != nil {
		return fmt.Errorf("%w; error: persistent configuration could not write ini.File to buffer", err)
	}

	// Write the buffer's contents to the file
	err = afero.WriteFile(fs, file, buffer.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("%w; error: persistent configuration could not be written to %q file", err, file)
	}
	return nil
}
