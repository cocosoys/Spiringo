package main

import (
	"flag"
	"log"
	"os"

	"github.com/spiringo/spiringo/internal/core/app"
	"github.com/spiringo/spiringo/internal/modules/builtin"
)

// 中文：main 是当前命令的程序入口。
// English: main is the entry point for this command.
func main() {
	env := flag.String("env", "", "runtime environment name")
	configDir := flag.String("config", "configs", "configuration directory")
	flag.Parse()

	if *env == "" {
		*env = firstNonEmpty(os.Getenv("APP_ENV"), "production")
	}
	if os.Getenv("SP_SERVER_ADDR") == "" {
		if port := os.Getenv("PORT"); port != "" {
			_ = os.Setenv("SP_SERVER_ADDR", ":"+port)
		}
	}

	application := app.New(app.WithConfigDir(*configDir), app.WithEnv(*env))
	builtin.RegisterAll(application)
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}

// 中文：firstNonEmpty 执行当前包中的对应流程。
// English: firstNonEmpty executes the corresponding workflow in this package.
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
