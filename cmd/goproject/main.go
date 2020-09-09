package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"

	"example.com/example/goproject/internal/app/goproject"
	"example.com/example/goproject/internal/pkg/config"
	"example.com/example/goproject/pkg/db"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := NewRootCommand().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "goproject",
	}
	cmd.AddCommand(NewApp())
	cmd.AddCommand(NewDBMigrate())
	cmd.AddCommand(NewDBRollback())
	return cmd
}

func NewApp() *cobra.Command {
	app := goproject.NewAPP()

	cmd := &cobra.Command{
		Use: "app",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return app.Init()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Start()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return app.Close()
		},
	}

	return cmd
}

func NewDBMigrate() *cobra.Command {
	cmd := &cobra.Command{
		Use: "dbmigrate",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			return RunMigration(func(m db.Migrator) error { return m.Up() })
		},
	}
	return cmd
}

func NewDBRollback() *cobra.Command {
	cmd := &cobra.Command{
		Use: "dbmigrate",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			return RunMigration(func(m db.Migrator) error { return m.Rollback() })
		},
	}
	return cmd
}

func RunMigration(f func(m db.Migrator) error) error {
	pgDB, err := db.New(config.NewDBConfig())
	if err != nil {
		return err
	}
	defer pgDB.Close()

	m, err := db.NewPostgresMigrator(config.MigrationPath(), pgDB.DB)
	if err != nil {
		return err
	}

	return f(m)
}
