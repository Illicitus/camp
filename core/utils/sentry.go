package utils

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"log"
)

const TagLocation = "location"

type LocalHub struct {
	hub          *sentry.Hub
	tag          string
	functionName string
	prodMode     bool
}

func NewLocalHub(functionName string, prodMode bool) *LocalHub {
	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag(TagLocation, functionName)
	})

	return &LocalHub{
		hub:          localHub,
		tag:          TagLocation,
		functionName: functionName,
		prodMode:     prodMode,
	}
}

func (lh *LocalHub) ErrorMessageHandler(err error, msg string) {
	if err != nil {
		if lh.prodMode {
			lh.hub.CaptureMessage(fmt.Sprintf("%s: %s", msg, err))
		}
		log.Fatalf(fmt.Sprintf("%s: %s", msg, err))
	}
}

func (lh *LocalHub) ErrorHandler(err error) {
	if err != nil {
		if lh.prodMode {
			lh.hub.CaptureException(err)
		}
		log.Fatal(err)
	}
}
