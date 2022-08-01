package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"good/internal/config"
	"good/internal/core"
	"good/internal/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "good",
	Short: "good is an out-of-the-box go framework",
	Long:  "good is an out-of-the-box go framework, which provide many component for micro service.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preRun()
	},
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

var (
	configPath string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "./config.yaml", "specify config file path")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("root cmd execute error: %v\n", err)
	}
}

func preRun() {
	if err := config.Setup(configPath); err != nil {
		log.Fatalf("config setup error: %v\n", err)
	}
	if err := core.Setup(config.C); err != nil {
		log.Fatalf("core setup error: %v\n", err)
	}
}

func run() {
	r := router.Init()

	srv := http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.C.Application.Host, config.C.Application.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(config.C.Application.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.C.Application.WriteTimeout) * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe: %s", err)
		}
	}()
	log.Printf("ListenAndServe: %s\n", srv.Addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server down error: %s", err.Error())
	}
	log.Println("server closed")
}
