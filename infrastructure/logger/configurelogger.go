package logger

import (
	"bytes"
	"fmt"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v3"
	"os"
	"strings"
	"time"
)

type outputSplitter struct{}

func (splitter *outputSplitter) Write(p []byte) (n int, err error) {
	if bytes.Contains(p, []byte("\"level\":\"info\"")) ||
		bytes.Contains(p, []byte("\"level\":\"debug\"")) ||
		bytes.Contains(p, []byte("\"level\":\"warning\"")) {
		return os.Stdout.Write(p)
	}
	return os.Stderr.Write(p)
}

var (
	logConfig     LoggerConfiguration
	elasticClient *elastic.Client
)

type LoggerConfiguration struct {
	ServiceName string
	Environment string
	ElasticUrl  string
}

func ConfigureLogger(config LoggerConfiguration) {
	logConfig = config
	if logConfig.Environment == "Production" || logConfig.Environment == "Staging" {
		log.SetLevel(log.InfoLevel)
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05-0700",
		})
		log.SetReportCaller(true)
		log.SetOutput(&outputSplitter{})

		client, err := elastic.NewClient(elastic.SetURL(logConfig.ElasticUrl), elastic.SetSniff(false))
		if err != nil {
			log.Panic(err)
		}
		elasticClient = client
		hook, err := elogrus.NewAsyncElasticHookWithFunc(elasticClient, logConfig.ServiceName,
			log.InfoLevel, getTodayElasticIndexName)

		if err != nil {
			log.Panic(err)
		}
		log.AddHook(hook)
	} else {
		log.SetLevel(log.DebugLevel)
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05-0700",
			PrettyPrint:     true,
		})
		log.SetOutput(os.Stdout)
	}
}

func getTodayElasticIndexName() string {
	return fmt.Sprintf("jabama_%s-log-%s", strings.ToLower(logConfig.ServiceName), time.Now().Format("2006-01-02"))
}
