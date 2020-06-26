package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/justinas/alice"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

type Config struct {
	Addr string `default:":8080"`
}

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Str("svc", "duck").Logger()

	cfg, err := parseConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("parsing config")
	}

	mw := setupMiddleware(logger)

	mux := http.NewServeMux()

	mux.Handle("/healthz", mw.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "healthy")
	})))

	mux.Handle("/", mw.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "quack")
	})))

	srv := &http.Server{Addr: cfg.Addr, Handler: mux}

	logger.Info().Msgf("server starting, listening on addr %s", cfg.Addr)

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal().Err(err).Msg("server error")
	}
}

func setupMiddleware(log zerolog.Logger) alice.Chain {
	c := alice.New()
	c = c.Append(hlog.NewHandler(log))

	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("served request")
	}))
	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))
	c = c.Append(hlog.RefererHandler("referer"))
	c = c.Append(hlog.UserAgentHandler("useragent"))

	return c
}

func parseConfig() (Config, error) {
	var cfg Config
	err := envconfig.Process("duck", &cfg)

	return cfg, err
}
