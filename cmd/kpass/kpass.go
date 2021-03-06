package main

import (
	"flag"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/seccom/kpass/src"
	"github.com/seccom/kpass/src/logger"
)

var (
	address  = flag.String("addr", "127.0.0.1:8088", `Auth service address to listen on.`)
	dbPath   = flag.String("dbpath", "./kpass.kdb", `KPass database pass.`)
	devMode  = flag.Bool("dev", false, "Development mode, will use memory database as default.")
	certFile = flag.String("certFile", "", `certFile path, used to create TLS service, support HTTP/2.`)
	keyFile  = flag.String("keyFile", "", `keyFile path, used to create TLS service, support HTTP/2.`)
)

func main() {
	flag.Parse()
	if !strings.HasSuffix(*dbPath, ".kdb") {
		panic(`Invalid dbpath, must have ".kdb" as extension.`)
	}

	if os.Getenv("APP_ENV") == "" {
		if *devMode {
			os.Setenv("APP_ENV", "development")
		} else {
			os.Setenv("APP_ENV", "production")
		}
	}

	env := os.Getenv("APP_ENV")
	app := src.New(*dbPath, env)

	go func() {
		host := "http://" + app.Server.Addr
		logger.Info("Start KPass: " + host)
		time.Sleep(time.Second)
		startBrowser(host)
	}()
	if *certFile != "" && *keyFile != "" {
		logger.Fatal(app.ListenTLS(*address, *certFile, *keyFile))
	} else {
		logger.Fatal(app.Listen(*address))
	}
}

// startBrowser tries to open the URL in a browser
// and reports whether it succeeds.
func startBrowser(url string) bool {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}
