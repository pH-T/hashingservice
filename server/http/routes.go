package http

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func initRoutes(hs hashservice) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/health", LoggingMiddleware(GETMiddleware(health())))
	mux.Handle("/hash", LoggingMiddleware(TimeMiddleware(POSTMiddleware(hash(hs)))))
	mux.Handle("/verify", LoggingMiddleware(TimeMiddleware(POSTMiddleware(verify(hs)))))
	return mux
}

func health() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("up"))
	})
}

func hash(hs hashservice) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type hashreq struct {
			Pw string `json:"pw"`
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			reqid := r.Context().Value("reqid")
			log.Printf("[%s] Error @ reading hash request: %s", reqid, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var req hashreq
		err = json.Unmarshal(b, &req)
		if err != nil {
			reqid := r.Context().Value("reqid")
			log.Printf("[%s] Error @ decoding hash request: %s", reqid, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		hash, err := hs.Hash(r.Context(), req.Pw)
		if err != nil {
			reqid := r.Context().Value("reqid")
			log.Printf("[%s] Error @ hashing request: %s", reqid, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		type hashresp struct {
			Hash string `json:"hash"`
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(hashresp{Hash: hash})
		if err != nil {
			reqid := r.Context().Value("reqid")
			log.Printf("[%s] Error @ responding to hashing request: %s", reqid, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})
}

func verify(hs hashservice) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type verifyreq struct {
			Pw   string `json:"pw"`
			Hash string `json:"hash"`
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			reqid := r.Context().Value("reqid")
			log.Printf("[%s] Error @ reading verify request: %s", reqid, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var req verifyreq
		err = json.Unmarshal(b, &req)
		if err != nil {
			reqid := r.Context().Value("reqid")
			log.Printf("[%s] Error @ decoding verify request: %s", reqid, err)

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		type hashresp struct {
			Verified bool `json:"verified"`
			Selfmade bool `json:"selfmade"`
		}
		resp := hashresp{}
		verified, selfmade, err := hs.Verify(r.Context(), req.Pw, req.Hash)
		if err != nil {
			reqid := r.Context().Value("reqid")
			log.Printf("[%s] Error @ verify request: %s", reqid, err)

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if verified {
			resp.Verified = true
			resp.Selfmade = false
			if selfmade {
				resp.Selfmade = true
			}

		} else {
			resp.Verified = false
			resp.Selfmade = false
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			reqid := r.Context().Value("reqid")
			log.Printf("[%s] Error @ responding to verify request: %s", reqid, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})
}
