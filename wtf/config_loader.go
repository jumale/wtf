package wtf

import (
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/radovskyb/watcher"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// ConfigLoader is responsible for unmarshaling a yaml configuration, provided
// in a file, into application config and widgets configs. Also it creates
// a sample config file, if there is no config yet.
type ConfigLoader struct {
	// Path to source config file
	configPath string

	// Unmarshaled appView config
	appConfig *AppConfig

	// Intermediate (partially unmarshaled) widgets configs, ready to be
	// finally unmarshaled into specific config objects.
	widgetsConfig []interimWidgetConfig
}

// Creates a new instance of ConfigLoader. The "appDir" specifies
func NewConfigLoader(configFile string) (*ConfigLoader, error) {
	fs := FileSystem{}
	configFile, err := fs.ExpandHomeDir(configFile)
	if err != nil {
		return nil, err
	}

	configPath, err := createConfigFile(configFile)
	if err != nil {
		return nil, err
	}

	return &ConfigLoader{configPath: configPath}, nil
}

func (cl *ConfigLoader) LoadConfig() error {
	cl.appConfig = &AppConfig{}

	data, err := ioutil.ReadFile(cl.configPath)
	if err != nil {
		return loadError(err)
	}

	loaded := Config{}
	err = yaml.Unmarshal(data, &loaded)
	if err != nil {
		return loadError(err)
	}

	cfg := Config{App: GetDefaultAppConfig()}
	err = mergo.Merge(&cfg, loaded, mergo.WithOverride)
	if err != nil {
		return loadError(err)
	}

	cl.appConfig = &cfg.App
	cl.widgetsConfig = cfg.WidgetsConfig

	return nil
}

func (cl *ConfigLoader) unmarshalWidgetConfig(
	cfg interimWidgetConfig,
	target interface{},
	logger Logger,
) error {

	if globalsMerger, ok := target.(defaultConfigSetter); ok {
		logger.Debugf("APP: set config globals for '%s' widget", cfg.Type)
		globalsMerger.setDefaults(*cl.appConfig)
	}

	err := yaml.Unmarshal(cfg.ConfigYaml, target)
	if err != nil {
		return loadError(err)
	}

	return nil
}

type OnConfigChange func() error

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

				err = onChange()
				if err != nil {
					log.Fatalln(err)
				}

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

type interimWidgetConfig struct {
	BaseWidgetConfig `yaml:",inline"`
	ConfigYaml       []byte
}

func (iwc *interimWidgetConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// unmarshal base widget config
	err := unmarshal(&iwc.BaseWidgetConfig)
	if err != nil {
		return err
	}

	// store the raw yaml config for the widget
	mapConfig := make(map[interface{}]interface{})
	err = unmarshal(&mapConfig)
	if err != nil {
		return err
	}
	iwc.ConfigYaml, err = yaml.Marshal(mapConfig)
	return err
}

func loadError(err error) error {
	return errors.Errorf("could not properly load WTF config: %s", err)
}

// CreateConfigFile creates a simple config file in the config directory if
// one does not already exist
func createConfigFile(filePath string) (string, error) {
	fs := FileSystem{}
	filePath, err := fs.CreateFile(filePath)
	if err != nil {
		return "", err
	}

	// If the file is empty, write to it
	file, _ := os.Stat(filePath)

	if file.Size() == 0 {
		err = ioutil.WriteFile(filePath, []byte(sampleConfig), 0644)
	}
	return filePath, errors.WithStack(err)
}

func expandConfigDir(configDir string) (string, error) {
	fs := FileSystem{}
	return fs.ExpandHomeDir(configDir)
}
