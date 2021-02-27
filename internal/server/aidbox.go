package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/zerolog"
)

func errorResp(log zerolog.Logger, w http.ResponseWriter, err error) {
	log.Error().Err(err).Msg("handling AiDBOX operation")
	b, _ := json.Marshal(AppOperationErrorResponse{Error: err.Error()})
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(b); err != nil {
		log.Error().Err(err).Msg("handling AiDBOX operation, writing error response")
	}
}

func AidboxOperationHandler(uri string, log zerolog.Logger, next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != uri {
			next.ServeHTTP(w, r)
			return
		}

		var op AppOperation
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&op); err != nil {
			errorResp(log, w, err)
			return
		}

		body := bytes.NewBuffer(op.Request.Resource)
		if len(op.Operation.Request) < 3 {
			errorResp(log, w, fmt.Errorf("invalid operation request: %v", op.Operation.Request))
			return
		}

		query := make(url.Values)
		for k, v := range op.Request.Params {
			query.Add(k, v)
		}
		url := &url.URL{Scheme: r.URL.Scheme, Host: r.Host, Path: strings.Join(op.Operation.Request[2:], "/"), RawQuery: query.Encode()}
		newReq, err := http.NewRequestWithContext(
			r.Context(),
			strings.ToUpper(op.Operation.Request[0]),
			url.String(),
			body,
		)
		if err != nil {
			errorResp(log, w, err)
			return
		}
		for k, v := range op.Request.Headers {
			newReq.Header.Add(k, v)
		}

		next.ServeHTTP(w, newReq)
	}
	return http.HandlerFunc(f)
}

type AppOperationErrorResponse struct {
	Error string `json:"error"`
}

type AppOperation struct {
	Type      string    `json:"type"`
	Request   Request   `json:"request"`
	Box       Box       `json:"box"`
	Operation Operation `json:"operation"`
}

type Box struct {
	BaseURL string `json:"base-url"`
}

type Operation struct {
	App          App         `json:"app"`
	Action       string      `json:"action"`
	Module       string      `json:"module"`
	Request      []string    `json:"request"`
	ID           string      `json:"id"`
	ResourceType string      `json:"resourceType"`
	Meta         Meta        `json:"meta"`
	W            interface{} `json:"w"`
}

type App struct {
	ID           string `json:"id"`
	ResourceType string `json:"resourceType"`
}

type Meta struct {
	LastUpdated string `json:"lastUpdated"`
	CreatedAt   string `json:"createdAt"`
	VersionID   string `json:"versionId"`
}

type Request struct {
	Resource    json.RawMessage `json:"resource"`
	Params      Values          `json:"params"`
	RouteParams Values          `json:"route-params"`
	Headers     Values          `json:"headers"`
	OauthClient OauthClient     `json:"oauth/client"`
}

type Headers struct {
	AcceptEncoding string `json:"accept-encoding"`
	Authorization  string `json:"authorization"`
	Connection     string `json:"connection"`
	ContentLength  string `json:"content-length"`
	ContentType    string `json:"content-type"`
	Host           string `json:"host"`
	UserAgent      string `json:"user-agent"`
	XRequestID     string `json:"x-request-id"`
}

type OauthClient struct {
	Secret       string   `json:"secret"`
	Source       string   `json:"source"`
	FirstParty   bool     `json:"first_party"`
	GrantTypes   []string `json:"grant_types"`
	ResourceType string   `json:"resourceType"`
	ID           string   `json:"id"`
	Meta         Meta     `json:"meta"`
}

type Values map[string]string
