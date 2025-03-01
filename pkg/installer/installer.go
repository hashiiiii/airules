package installer

// Installer interface defines methods that must be implemented by installers for each editor
type Installer interface {
	// InstallLocal installs the local configuration file
	InstallLocal() error
	
	// InstallGlobal installs the global configuration file
	InstallGlobal() error
	
	// InstallAll installs both configuration files
	InstallAll() error
}
