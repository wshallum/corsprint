package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/wshallum/corsprint/printlib"
	"net/http"
	"os"
)

var (
	allowedOrigin string
)

type listPrintersResultJSON struct {
	Printers []listPrintersPrinterJSON `json:"printers"`
}

type listPrintersPrinterJSON struct {
	Name string `json:"name"`
}

func toJson(printers []printlib.Printer) ([]byte, error) {
	var result listPrintersResultJSON
	result.Printers = make([]listPrintersPrinterJSON, len(printers))
	for i, p := range printers {
		result.Printers[i] = listPrintersPrinterJSON{Name: p.Name()}
	}
	return json.Marshal(result)
}

func originMatches(claimedOrigin, allowedOrigin string) bool {
	if allowedOrigin == "*" {
		return true
	}
	return claimedOrigin == allowedOrigin
}

func listPrintersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		fmt.Fprintf(w, "This resource supports GET and OPTIONS\n")
	} else if r.Method == "GET" {
		if !originMatches(r.Header.Get("Origin"), allowedOrigin) {
			w.WriteHeader(403) // forbidden
			fmt.Fprintf(w, "Origin not allowed\n")
			return
		}
		printers, err := printlib.ListPrinters()
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error: %s\n", err.Error())
			return
		}
		bytes, err := toJson(printers)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Error: %s\n", err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	} else {
		w.WriteHeader(405) // method not allowed
	}

}

func printHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		fmt.Fprintf(w, "This resource supports POST and OPTIONS\n")
	} else if r.Method == "POST" {
		if !originMatches(r.Header.Get("Origin"), allowedOrigin) {
			w.WriteHeader(403) // forbidden
			fmt.Fprintf(w, "Origin not allowed\n")
			return
		}
		printer := r.FormValue("printer")
		data := r.FormValue("data")
		if data == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			fmt.Fprintf(w, "{\"error\": \"NO_DATA_SPECIFIED\"}")
			return
		}
		if printer == "" {
			defPrinter, err := printlib.GetDefaultPrinter()
			if err != nil {
				// just assume no default printer if we can't get one
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(400)
				fmt.Fprintf(w, "{\"error\": \"NO_DEFAULT_PRINTER\"}")
			}
			printer = defPrinter
		}
		_, err := printlib.Print(printer, data)
		if err != nil {

		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"ok\": true}")
	} else {
		w.WriteHeader(405) // method not allowed
	}
}

func catchAllHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	fmt.Fprintf(w, "404 Not Found\n")
}

func main() {
	listenAddress := flag.String("listen-address", "127.0.0.1:8080", "Listen address (ip:port)")
	origin := flag.String("origin", "", "Allowed origin")

	flag.Parse()
	if *origin == "" {
		fmt.Fprintf(os.Stderr, "Error: must specify allowed origin (-origin)\n")
		os.Exit(1)
	}
	fmt.Printf("Addr %s Origin %s\n", *listenAddress, *origin)
	allowedOrigin = *origin
	http.HandleFunc("/printers", listPrintersHandler)
	http.HandleFunc("/print", printHandler)
	http.HandleFunc("/", catchAllHandler)
	http.ListenAndServe(*listenAddress, nil)
}

// vim: ft=go:sw=8:ts=8:sts=8
