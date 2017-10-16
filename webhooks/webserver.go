package webhooks

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/nlopes/slack"

	"golang.org/x/crypto/acme/autocert"
)

const (
	domain           = "unidiv.space"
	certificateCache = "certs"
	port             = "443"
)

var messageBus *slack.RTM

func StartWebhooksServer(rtm *slack.RTM) {
	messageBus = rtm
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
		Cache:      autocert.DirCache(certificateCache),
	}

	http.HandleFunc("/webhooks/jenkins", jenkinsHandler)

	server := &http.Server{
		Addr: ":" + port,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}
	if os.Getenv("DEVELOP") != "" {
		log.Print("Starting webhooks listener on port 8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
			os.Exit(3)
		}
	} else {
		log.Print("Starting webhooks listener on port ", port)
		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Fatal(err)
			os.Exit(3)
		}
	}
}
