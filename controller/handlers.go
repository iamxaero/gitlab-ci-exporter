package controller

import (
	"fmt"
	"io"
	"net/http"

	"github.com/cloudflare/cfssl/log"
)

func (c *Controller) Health(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Gitlab CI Exporter is OK")
}

func (c *Controller) Webhook(w http.ResponseWriter, r *http.Request) {
	msg := "skipped"
	var err error

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Debug(err)
	}
	defer r.Body.Close()

	if r.Header.Get("X-Gitlab-Event") == "Job Hook" {
		msg, err = c.inJob(data)
		if err != nil {
			log.Debugf("Error! %v: %v\n", msg, err)
			log.Debug(string(data))
		}
	}

	if r.Header.Get("X-Gitlab-Event") == "Pipeline Hook" {
		msg, err = c.inPipe(data)
		if err != nil {
			log.Debugf("Error! %v: %v\n", msg, err)
			log.Debug(string(data))
		}
	}

	fmt.Fprint(w, msg)
}
