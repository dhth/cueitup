package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	t "github.com/dhth/cueitup/internal/types"
)

var (
	errCouldntStartServer     = errors.New("couldn't start server")
	errForcefulShutdownFailed = errors.New("forceful shutdown failed")
)

func Serve(
	sqsClient *sqs.Client,
	config t.Config,
	open bool,
) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", getIndex)
	mux.HandleFunc("GET /priv/static/favicon.png", getFavicon)
	mux.HandleFunc("GET /priv/static/cueitup.css", getCSS)
	mux.HandleFunc("GET /priv/static/cueitup.mjs", getJS)
	mux.HandleFunc("GET /api/config", getConfig(config))
	mux.HandleFunc("GET /api/fetch", getMessages(sqsClient, config))
	muxWithCors := corsMiddleware(mux)

	port, ok := findOpenPort(startPort, endPort)
	if !ok {
		return fmt.Errorf("%w; checked between %d-%d", errNoPortOpen, startPort, endPort)
	}

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	server := &http.Server{
		Addr:    addr,
		Handler: muxWithCors,
	}

	addrWithProtocol := fmt.Sprintf("http://%s", addr)

	serverErrChan := make(chan error)

	go func(errChan chan<- error) {
		if open {
			fmt.Printf("Starting server at %s.\n", addrWithProtocol)
		} else {
			fmt.Printf("Starting server. Open %s in your browser.\n", addrWithProtocol)
		}
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}(serverErrChan)

	if open {
		go func() {
			time.Sleep(time.Millisecond * 1000)
			err := openURL(addrWithProtocol)
			if err != nil {
				fmt.Fprintf(os.Stderr, "couldn't open URL: %s", err.Error())
			}
		}()
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		shutDownCtx, shutDownRelease := context.WithTimeout(context.Background(), time.Second*3)
		defer shutDownRelease()

		if err := server.Shutdown(shutDownCtx); err != nil {
			fmt.Printf("Error shutting down: %s\nTrying forceful shutdown...\n", err.Error())
			if err := server.Close(); err != nil {
				return fmt.Errorf("%w:: %s", errForcefulShutdownFailed, err.Error())
			}
		}
		fmt.Printf("\nbye 👋\n")
	case err := <-serverErrChan:
		return fmt.Errorf("%w: %s", errCouldntStartServer, err.Error())
	}

	return nil
}
