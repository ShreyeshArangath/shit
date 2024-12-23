package models

// Config represents the configuration for the repository.
type Config struct {
	Core CoreSection `ini:"core"`
}

// CoreSection represents the [core] section in the config file.
type CoreSection struct {
	RepositoryFormatVersion int  `ini:"repositoryformatversion"`
	FileMode                bool `ini:"filemode"`
	Bare                    bool `ini:"bare"`
}
