package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

func calculateKey(uid, key, ts, qname, ak string) string {
	data := uid + key + ts + qname + ak
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func forwardDNSRequest(uid, ak, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ts := fmt.Sprintf("%d", time.Now().Unix())
		qname := r.URL.Query().Get("name")

		calculatedKey := calculateKey(uid, secret, ts, qname, ak)

		params := url.Values{}
		params.Add("uid", uid)
		params.Add("ak", ak)
		params.Add("ts", ts)
		params.Add("key", calculatedKey)
		params.Add("name", qname)

		dohURL := "https://223.5.5.5/resolve?" + params.Encode()

		resp, err := http.Get(dohURL)
		if err != nil {
			http.Error(w, "Failed to forward request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/dns-message")
		w.Write(body)
	}
}

func main() {
	uid := flag.String("uid", "", "User ID")
	ak := flag.String("ak", "", "AccessKey ID")
	secret := flag.String("secret", "", "AccessKey Secret")
	port := flag.String("port", "8080", "Port to listen")
	flag.Parse()

	if *uid == "" || *ak == "" || *secret == "" {
		fmt.Println("--uid, --ak and --secret is mandatory and cannot leave it empty.")
		os.Exit(1)
	}

	http.HandleFunc("/dns-query", forwardDNSRequest(*uid, *ak, *secret))
	fmt.Println("Server is running on port: " + *port)
	http.ListenAndServe(":" + *port, nil)
}
