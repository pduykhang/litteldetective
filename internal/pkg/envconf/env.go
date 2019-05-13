package envconf

import (
	"github.com/PhamDuyKhang/littledetective/internal/pkg/flog"
	"github.com/kelseyhightower/envconfig"
)

func Load(envstruct interface{}) error {
	l := flog.New()
	l.SetLocal("envconf")
	if err := envconfig.Process("", envstruct); err != nil {
		l.Errorf("can't get config from environment variable,with err = %v ", err)
		return err
	}
	l.Logger.Infof("get env variable successfully")
	return nil
}
