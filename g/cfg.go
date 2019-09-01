package g

import (
	"encoding/json"
	"fmt"
	"github.com/toolkits/file"
	"goexamples/kafka/msgMiddleware/model"
	"runtime"
	"sync"
)

const (
	VERSION = "0.0.1"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

var (
	config *model.Config
	lock   = new(sync.RWMutex)
)

func Config() *model.Config {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func ParseConfig(cfg string) error {
	var c model.Config

	if cfg == "" {
		return fmt.Errorf("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		return fmt.Errorf("config file %s is nonexistent", cfg)
	}

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		return fmt.Errorf("read config file %s fail %s", cfg, err)
	}

	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return fmt.Errorf("parse config file %s fail %s", cfg, err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	return nil
}
