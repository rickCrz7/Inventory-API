package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rickCrz7/Inventory-API/utils"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	devFlag := flag.Bool("dev", false, "is it running in development mode")
	flag.Parse()

	// Load configuration from config file
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panicf("Could not load app.yml configuration file: %v", err)
	}

	// Setup logger
	lumberjackLogrotate := &lumberjack.Logger{
		Filename:   viper.GetString("log.file"),
		MaxSize:    viper.GetInt("log.max-size"),    // Max megabytes before log is rotated
		MaxBackups: viper.GetInt("log.max-backups"), // Max number of old log files to keep
		MaxAge:     viper.GetInt("log.max-age"),     // Max number of days to retain log files
		Compress:   false,
	}

	mode := "production"
	postgresURI := viper.GetString("postgres.prod")
	setLimits := true
	if *devFlag {
		log.SetReportCaller(true)
		log.SetFormatter(&log.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: "2006/01/02 15:04:05",
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return "", fmt.Sprintf("\t%s:%d", filename, f.Line)
			},
		})
		mode = "development"
		postgresURI = viper.GetString("postgres.dev")
		setLimits = false
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}

	logMultiWriter := io.MultiWriter(os.Stdout, lumberjackLogrotate)
	log.SetOutput(logMultiWriter)
	switch viper.GetString("log.level") {
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "TRACE":
		log.SetLevel(log.TraceLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}

	log.WithFields(log.Fields{
		"Runtime Version": runtime.Version(),
		"Number of CPUs":  runtime.NumCPU(),
		"Arch":            runtime.GOARCH,
	}).Infof("Starting %s", viper.GetString("app.name"))

	// Setup Postgres connection
	log.Infof("connecting to Postgres: %s", mode)
	pdb, err := utils.OpenDB(postgresURI, setLimits)
	if err != nil {
		log.Fatalf("Could not connect to Postgres: %v", err)
	}
	err = pdb.Ping(context.Background())
	if err != nil {
		log.Fatalf("Could not ping Postgres: %v", err)
	}
	log.Printf("Connected to Postgres: %s", mode)

	// Setup router
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	srv := &http.Server{
		Handler: r,
		Addr:    viper.GetString("app.addr"),
	}
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server has stopped: %v", err)
		}
	}()
	log.Printf("Server for %s started on %s", viper.GetString("app.name"), viper.GetString("app.addr"))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	<-done
	log.Printf("Shutting down server for %s", viper.GetString("app.name"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		log.Println("Closing")
		if pdb != nil {
			pdb.Close()
			log.Println("Database closed gracefully")
		}

		cancel()
	}()

	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Server Shutdown Failed: %v", err)
	}
	log.Print("Server shutdown gracefully")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		// skip logging for health check
		if request.URL.Path == "/purchase/api/v1/healthz" {
			next.ServeHTTP(response, request)
			return
		}
		start := time.Now()
		next.ServeHTTP(response, request)
		// log.Printf("%s %s %s %s", getIPAddress(request), request.Method, request.RequestURI, time.Since(start).String())
		log.WithFields(log.Fields{
			"IP":     getIPAddress(request),
			"Method": request.Method,
			"URI":    request.RequestURI,
			"Cost":   time.Since(start).String(),
		}).Info("Handler called")
	})
}

func getIPAddress(r *http.Request) string {
	// for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
	for _, h := range []string{"X-Forwarded-For"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		for i := 0; i < len(addresses); i++ {
			ip := strings.TrimSpace(addresses[i])
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() {
				// bad address, go to next
				continue
			}
			return ip
		}
	}
	return "localhost"
}
