package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/sir-hassan/luguz/captcha"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const apiPort = 8080

// ErrorReply is used to encode api errors.
type ErrorReply struct {
	Error string `json:"error"`
}

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	level.Info(logger).Log("msg", "starting luguz...")

	handler := createHandler(logger, captcha.NewCachedGenerator())

	http.Handle("/", handler)
	level.Info(logger).Log("msg", "listening started", "port", apiPort)
	ln, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(apiPort))
	if err != nil {
		level.Error(logger).Log("msg", "listening on port failed", "port", apiPort, "error", err)
		return
	}
	server := &http.Server{}

	go func() {
		err := server.Serve(ln)
		if err != http.ErrServerClosed {
			level.Error(logger).Log("msg", "closing server failed", "error", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	level.Info(logger).Log("msg", "signal received", "sig", sig)
	level.Info(logger).Log("msg", "terminating server...")

	if err := server.Shutdown(context.Background()); err != nil {
		level.Error(logger).Log("msg", "terminating server failed", "error", err)
	} else {
		level.Info(logger).Log("msg", "terminating server succeeded")
	}
}

func createHandler(logger log.Logger, generator captcha.Generator) http.Handler {
	handler := mux.NewRouter()
	handler.HandleFunc("/captcha", createCaptchaHandleFunc(logger, generator))
	handler.NotFoundHandler = http.HandlerFunc(createNotFoundHandleFunc(logger))
	return handler
}

func createNotFoundHandleFunc(logger log.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		writeErrorReply(logger, w, 404, "the requested path is not found")
	}
}

func createCaptchaHandleFunc(logger log.Logger, generator captcha.Generator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		cpt, err := generator.Generate()
		if err != nil {
			level.Error(logger).Log("msg", "generating captcha", "error", err)
		}
		if _, err = w.Write(cpt.Bytes()); err != nil {
			level.Error(logger).Log("msg", "writing connection", "error", err)
		}
		if err = generator.Release(cpt); err != nil {
			level.Error(logger).Log("msg", "releasing captcha", "error", err)
		}
	}
}

func writeReply(logger log.Logger, w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write([]byte(message + "\n")); err != nil {
		level.Error(logger).Log("msg", "writing connection", "error", err)
	}
}

func writeErrorReply(logger log.Logger, w http.ResponseWriter, statusCode int, message string) {
	jsonString, err := json.Marshal(ErrorReply{Error: message})
	if err != nil {
		level.Error(logger).Log("msg", "decoding reply", "error", err)
		return
	}
	writeReply(logger, w, statusCode, string(jsonString))
}
