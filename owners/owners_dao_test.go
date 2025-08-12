package owners

import (
	"fmt"
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var postgresURI string

func init() {
	viper.SetConfigName("app")
	viper.AddConfigPath("../config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	postgresURI = viper.GetString("postgres.prod") // Change to "postgres.dev" for development/local db
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006/01/02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("\t%s:%d", filename, f.Line)
		},
	})
	log.SetLevel(log.DebugLevel)
}
