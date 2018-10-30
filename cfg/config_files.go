package cfg

import (
	"github.com/senorprogrammer/wtf/wtf"
	"log"
	"os"
)

// ConfigDirV1 defines the path to the first version of configuration. Do not use this
const ConfigDirV1 = "~/.wtf/"

// ConfigDirV2 defines the path to the second version of the configuration. Use this.
const ConfigDirV2 = "~/.config/wtf/"

/* -------------------- AppConfig Migration -------------------- */

// MigrateOldConfig copies any existing configuration from the old location
// to the new, XDG-compatible location
func MigrateOldConfig() {
	fs := wtf.FileSystem{}
	srcDir, _ := fs.ExpandHomeDir(ConfigDirV1)
	destDir, _ := fs.ExpandHomeDir(ConfigDirV2)

	// If the old config directory doesn't exist, do not move
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return
	}

	// If the new config directory already exists, do not move
	if _, err := os.Stat(destDir); err == nil {
		return
	}

	// Time to move
	err := Copy(srcDir, destDir)
	if err != nil {
		panic(err)
	} else {
		log.Printf("Copied old config from %s to %s", srcDir, destDir)
	}

	// Delete the old directory if the new one exists
	if _, err := os.Stat(destDir); err == nil {
		err := os.RemoveAll(srcDir)
		if err != nil {
			log.Print(err.Error())
		}
	}
}
