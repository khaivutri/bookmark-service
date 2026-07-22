package sqldb

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Host     		string 		`default:"localhost" envconfig:"DB_HOST"`
	User     		string 		`default:"admin" envconfig:"DB_USER"`
	Password 		string 		`default:"admin" envconfig:"DB_PASS"`
	DBName   		string 		`default:"bookmark" envconfig:"DB_NAME"`
	Port     		string 		`default:"5432" envconfig:"DB_PORT"`
}

func newConfig(envPrefix string) (*config, error) {
	cfg := &config{}
	if err := envconfig.Process(envPrefix, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", 
						cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port)
}