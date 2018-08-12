package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const VersionMajor = 0
const VersionMinor = 0
const VersionRevision = 1

type VersionInfo struct {
	Major    int
	Minor    int
	Revision int
}

func MakeVersionInfo() VersionInfo {
	return VersionInfo{
		VersionMajor,
		VersionMinor,
		VersionRevision,
	}
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

func (writer *LoggingResponseWriter) WriteHeader(code int) {
	writer.status = code
	writer.ResponseWriter.WriteHeader(code)
}

func LogHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Printf("--> %s %s %s\n", r.RemoteAddr, r.Method, r.URL)

		writer := NewLoggingResponseWriter(w)

		handler.ServeHTTP(writer, r)

		log.Printf("<-- %d - %s %s %s %s", writer.status, http.StatusText(writer.status), r.RemoteAddr, r.Method, r.URL)
	})
}

func main() {

	var key string
	flag.StringVar(&key, "key", "key.pem", "ssl key")
	var certificate string
	flag.StringVar(&certificate, "certificate", "certificate.pem", "ssl certificate file")
	var logfile string
	flag.StringVar(&logfile, "log", "logfile.log", "name of log file")
	var port int
	flag.IntVar(&port, "port", 443, "port to run server on")

	flag.Parse()

	file, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log-file: %v", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, file)

	log.SetOutput(multiWriter)

	serveMux := http.NewServeMux()
	apiMux := http.NewServeMux()

	serveMux.Handle("/", http.FileServer(http.Dir("./com")))

	serveMux.Handle("/api/", http.StripPrefix("/api", apiMux))

	apiMux.HandleFunc("/version/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		js, err := json.Marshal(MakeVersionInfo())
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Write(js)
		}
	})

	server := http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: LogHandler(serveMux),
	}

	log.Println("starting notion server on port:", port)
	err = server.ListenAndServeTLS(certificate, key)
	if err != nil {
		log.Fatal("Failed to launch http server")
	}
}
