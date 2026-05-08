package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/spiringo/spiringo/internal/core/app"
	"github.com/spiringo/spiringo/internal/modules/builtin"
)

// 中文：main 是当前命令的程序入口。
// English: main is the entry point for this command.
func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

// 中文：run 执行当前包中的对应流程。
// English: run executes the corresponding workflow in this package.
func run(args []string) error {
	if len(args) > 0 && args[0] == "migrate" {
		return runMigrate(args[1:])
	}
	application, err := newApplication(args)
	if err != nil {
		return err
	}
	return application.Run()
}

// 中文：runMigrate 执行当前包中的对应流程。
// English: runMigrate executes the corresponding workflow in this package.
func runMigrate(args []string) error {
	if len(args) == 0 || args[0] != "up" {
		return fmt.Errorf("migrate requires: up")
	}
	application, err := newApplication(args[1:])
	if err != nil {
		return err
	}
	return application.Migrate(context.Background())
}

// 中文：newApplication 执行当前包中的对应流程。
// English: newApplication executes the corresponding workflow in this package.
func newApplication(args []string) (*app.App, error) {
	fs := flag.NewFlagSet("spiringo", flag.ContinueOnError)
	env := fs.String("env", "", "runtime environment name")
	configDir := fs.String("config", "configs", "configuration directory")
	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if *env == "" {
		*env = os.Getenv("APP_ENV")
	}

	opts := []app.Option{app.WithConfigDir(*configDir)}
	if *env != "" {
		opts = append(opts, app.WithEnv(*env))
	}

	application := app.New(opts...)
	builtin.RegisterAll(application)
	return application, nil
}
