package goat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var logPath, remoteURL, agent string
var isLogEnabled, isLogPushEnabledToRemote bool

func populateConfigurableVariables(log *viper.Viper) {
	isLogEnabled = log.GetBool("IsLogEnabled")
	isLogPushEnabledToRemote = log.GetBool("IsLogPushEnabledToRemote")
	logPath = log.GetString("LogPath")
	remoteURL = log.GetString("RemoteUrl")
	agent = log.GetString("Agent")
}

func loadConfig() {
	log := viper.New()
	log.SetConfigName("config")
	log.AddConfigPath(".")

	err := log.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err.Error()))
	}
	populateConfigurableVariables(log)
	log.WatchConfig()
	log.OnConfigChange(func(e fsnotify.Event) {
		populateConfigurableVariables(log)
	})
}

func pushError(input map[string]interface{}) {
	payloadBytes, err := json.Marshal(input)
	if err == nil {
		if isLogEnabled {
			f, err := os.OpenFile(logPath+time.Now().Format("2006-01-02")+".txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
			defer f.Close()
			if err == nil {
				stringMessage := string(payloadBytes)
				f.WriteString(stringMessage)
			} else {
				fmt.Println(err.Error())
			}
		}
		if isLogPushEnabledToRemote {
			client := &http.Client{}
			req, _ := http.NewRequest("POST", remoteURL, bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")
			client.Do(req)
		}
	}
}

//RecoverAndLogPanic - It catches the unexpected panic in the application and helps in recovering from it.
//If config file is placed in the working directory it logs data based on the attributes.
func RecoverAndLogPanic(next http.Handler) http.Handler {
	loadConfig()
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				input := collectErrorData(*req, err)
				pushError(input)

				res.WriteHeader(http.StatusInternalServerError)
				res.Header().Set("Content-Type", "application/json")
				errorText := fmt.Sprintf(" PANIC Defered : %v ", err)
				res.Write([]byte(errorText))

			}
		}()
		next.ServeHTTP(res, req)
	})
}

//RecoverAndLogGoRoutinePanic - It catches the unexpected panic in the goroutines and helps in recovering from it.
//If config file is placed in the working directory it logs data based on the attributes.
func RecoverAndLogGoRoutinePanic(req http.Request) {
	if err := recover(); err != nil {
		input := collectErrorData(req, err)
		pushError(input)
	}
}

func collectErrorData(req http.Request, err interface{}) map[string]interface{} {
	errorText := fmt.Sprintf("%v ", err)
	bodyBytes, err := ioutil.ReadAll(req.Body)
	bodyText := string(bodyBytes)
	trace := make([]byte, 1024)
	runtime.Stack(trace, false)
	stackTrace := fmt.Sprintf("%s", trace)
	input := map[string]interface{}{
		"Agent":   agent,
		"API":     req.URL.RequestURI(),
		"METHOD":  req.Method,
		"BODY":    bodyText,
		"PANIC":   errorText,
		"REFERER": req.Referer(),
		"IP":      req.RemoteAddr,
		"STACK":   stackTrace,
	}

	return input
}

//ReadData : It helps in reading the data from reader and populating it again it with the same data.
//This helps in using the request body again.
//e.g. In case of panic request body  is essential to the developers to simulate and debug the issue.
func ReadData(reader *io.ReadCloser, inputType interface{}) error {
	defer (*reader).Close()

	var err error
	if reader != nil {
		var buf bytes.Buffer
		tee := io.TeeReader(*reader, &buf)
		bodyBytes, err := ioutil.ReadAll(tee)
		*reader = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		if err == nil {
			err = json.Unmarshal(bodyBytes, &inputType)
		}
	} else {
		err = errors.New("reader is nil")
	}

	return err
}
