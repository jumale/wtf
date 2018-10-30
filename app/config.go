package app

import (
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/radovskyb/watcher"
	"github.com/senorprogrammer/wtf/wtf"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

var simpleConfig = `wtf:
  # disable auto-refresh, unless it's explicitly enabled in mods configs  
  refreshInterval: -1
  colors:
    border:
      focusable: darkslateblue
      focused: orange
      normal: gray
  grid:
    # uncomment to map columns manually by their size
    # columns: [40, 40]
    
    # or specify number of columns, each column with will be calculated automatically
    numCols: 2
    rows: [13, 13, 4]

  mods:
    clocks:
      enabled: true
      refreshInterval: 15
      sort: "alphabetical"
      position:
        top: 0
        left: 0
        height: 1
        width: 1
      colors:
        rows:
          even: "lightblue"
          odd: "white"
      locations:
        Avignon: "Europe/Paris"
        Barcelona: "Europe/Madrid"
        Dubai: "Asia/Dubai"
        Vancouver: "America/Vancouver"
        Toronto: "America/Toronto"

    system:
      enabled: true
      refreshInterval: 3600
      position:
        top: 0
        left: 1
        height: 1
        width: 1

    security:
      enabled: true
      refreshInterval: 3600
      position:
        top: 1
        left: 0
        height: 1
        width: 1

    textfile:
      enabled: true
      refreshInterval: 30
      filePath: "~/.config/wtf/config.yml"
      position:
        top: 1
        left: 1
        height: 1
        width: 1

    status:
      enabled: true
      refreshInterval: 1
      position:
        top: 2
        left: 0
        height: 1
        width: 2
`

type RootConfig struct {
	wtf.AppConfig `yaml:",inline"`
	WidgetsConfig map[string]interface{} `yaml:"mods"`
}

type WtfConfig struct {
	Root RootConfig `yaml:"wtf"`
}

type OnConfigChange func()

func NewConfigLoader(configDir, configFile string) (*ConfigLoader, error) {
	configPath, err := createConfigFile(path.Join(configDir, configFile))
	if err != nil {
		return nil, err
	}
	return &ConfigLoader{configPath: configPath}, nil
}

type ConfigLoader struct {
	configPath    string
	appConfig     *wtf.AppConfig
	widgetsConfig map[string][]byte
}

func (cl *ConfigLoader) LoadConfig() error {
	cl.appConfig = &wtf.AppConfig{}
	cl.widgetsConfig = make(map[string][]byte)

	data, err := ioutil.ReadFile(cl.configPath)
	if err != nil {
		return loadError(err)
	}

	cfg := WtfConfig{}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return loadError(err)
	}

	defaults := WtfConfig{Root: RootConfig{AppConfig: wtf.GetDefaultAppConfig()}}
	err = mergo.Merge(&cfg, defaults)
	if err != nil {
		return loadError(err)
	}

	cl.appConfig = &cfg.Root.AppConfig
	for widgetName, cfgMap := range cfg.Root.WidgetsConfig {
		cl.widgetsConfig[widgetName], err = yaml.Marshal(cfgMap)
		if err != nil {
			return loadError(err)
		}
	}

	return nil
}

func (cl *ConfigLoader) unmarshalWidgetConfig(name string, target interface{}, logger wtf.Logger) error {
	configYaml, ok := cl.widgetsConfig[name]
	if !ok {
		//return errors.Errorf("could not find config for '%s' widget in the config file", name)
		return nil // just skip a non-configured widget
	}

	if defaultsSetter, ok := target.(wtf.DefaultConfigSetter); ok {
		logger.Debugf("APP: set config defaults for '%s' widget", name)
		defaultsSetter.SetDefaults(*cl.appConfig)
	}

	err := yaml.Unmarshal(configYaml, target)
	if err != nil {
		return loadError(err)
	}

	return nil
}

func (cl *ConfigLoader) WatchChanges(onChange OnConfigChange) {
	watch := watcher.New()
	// notify write events.
	watch.FilterOps(watcher.Write)

	go func() {
		var err error
		for {
			select {
			case <-watch.Event:
				err = cl.LoadConfig()
				if err != nil {
					log.Fatalln(err)
				}
				onChange()

			case err := <-watch.Error:
				log.Fatalln(err)

			case <-watch.Closed:
				return
			}
		}
	}()

	// Watch config file for changes.
	if err := watch.Add(cl.configPath); err != nil {
		log.Fatalln(err)
	}

	// Start the watching process - it'll check for changes every 100ms.
	if err := watch.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func loadError(err error) error {
	return errors.Errorf("could not properly load WTF config: %s", err)
}

// CreateConfigFile creates a simple config file in the config directory if
// one does not already exist
func createConfigFile(filePath string) (string, error) {
	fs := wtf.FileSystem{}
	filePath, err := fs.CreateFile(filePath)
	if err != nil {
		return "", err
	}

	// If the file is empty, write to it
	file, _ := os.Stat(filePath)

	if file.Size() == 0 {
		err = ioutil.WriteFile(filePath, []byte(simpleConfig), 0644)
	}
	return filePath, errors.WithStack(err)
}

func expandConfigDir(configDir string) (string, error) {
	fs := wtf.FileSystem{}
	return fs.ExpandHomeDir(configDir)
}
