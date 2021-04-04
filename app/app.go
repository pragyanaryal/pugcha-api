package app

import (
	"crypto/tls"
	"gitlab.com/ProtectIdentity/pugcha-backend/config"
	_ "gitlab.com/ProtectIdentity/pugcha-backend/db"
	"gitlab.com/ProtectIdentity/pugcha-backend/log"
	"gitlab.com/ProtectIdentity/pugcha-backend/migrations"
	"gitlab.com/ProtectIdentity/pugcha-backend/router"
	"net/http"
	"time"
)

func BootstrapApp() {
	migrations.Migrate()

	routes := router.New()

	s := &http.Server{
		Addr:         config.Configuration.ServerPort,
		Handler:      routes,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig: &tls.Config{
			NextProtos: []string{"h2", "http/1.1"},
		},
	}

	if config.Configuration.TLS == true {
		
		go http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
		}))

		log.Info().Msgf("Starting server in %s in https mode", config.Configuration.ServerPort)
		
		if config.Configuration.Local == true {
			log.Fatal().Err(s.ListenAndServeTLS("certs/dev.pugcha.com.crt", "certs/dev.pugcha.com.key")).Msg("Fatal Server Error")
		} else {
			log.Fatal().Err(s.ListenAndServeTLS("/etc/letsencrypt/live/pugcha.com/fullchain.pem", "/etc/letsencrypt/live/pugcha.com/privkey.pem")).Msg("Fatal Server Error")
		}
		
	} else {
		log.Info().Msgf("Starting server in %s in http mode", config.Configuration.ServerPort)
		log.Fatal().Err(s.ListenAndServe()).Msg("Fatal Server Error")
	}
}
