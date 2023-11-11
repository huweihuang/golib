package logrus

import (
	"testing"

	"github.com/sirupsen/logrus"

	log "github.com/huweihuang/golib/logger/logrus"
)

func TestLogrus(t *testing.T) {
	// init log
	log.InitLogger("./log/logrus.log", "debug", "text", false)

	// Printf
	log.Logger().Debugf("test debugf, %s", "debugf")
	log.Logger().Infof("test infof, %s", "infof")
	log.Logger().Warnf("test warnf, %s", "warnf")
	log.Logger().Errorf("test errorf, %s", "errorf")

	// WithField
	log.Logger().WithField("field1", "debug").Debug("test field, debug")
	log.Logger().WithField("field1", "info").Info("test field, info")
	log.Logger().WithField("field1", "warn").Warn("test field, warn")
	log.Logger().WithField("field1", "error").Error("test field, error")

	// WithFields
	log.Logger().WithFields(logrus.Fields{
		"fields1": "fields1_value",
		"fields2": "fields2_value",
	}).Debug("test fields, debug")

	log.Logger().WithFields(logrus.Fields{
		"fields1": "fields1_value",
		"fields2": "fields2_value",
	}).Info("test fields, info")

	log.Logger().WithFields(logrus.Fields{
		"fields1": "fields1_value",
		"fields2": "fields2_value",
	}).Warn("test fields, warn")

	log.Logger().WithFields(logrus.Fields{
		"fields1": "fields1_value",
		"fields2": "fields2_value",
	}).Error("test fields, error")
}

func TestLog(t *testing.T) {
	log.Logger().WithField("field1", "field1").Info("test field, info")
}
