package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"slices"

	"dev.acmcsuf.com/fullyhacks-qrms/server"
	"dev.acmcsuf.com/fullyhacks-qrms/sqldb"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"libdb.so/hserve"
	"libdb.so/xcsv"
)

var (
	addr    = ":5767"
	db      = "fullyhacks-qrms.db"
	verbose = false
)

func init() {
	flag.StringVar(&addr, "addr", addr, "address to listen on")
	flag.StringVar(&db, "db", db, "path to the database file")
	flag.BoolVar(&verbose, "verbose", verbose, "enable verbose logging")
	flag.Usage = func() {
		out := flag.CommandLine.Output()

		fmt.Fprintf(out, "Usage: %s [flags] [command]\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(out)	

		fmt.Fprintln(out, "Commands:")
		fmt.Fprintln(out, "  import-users <file>  Import users from a CSV file")
		fmt.Fprintln(out, "  serve                Start the server")
		fmt.Fprintln(out)

		fmt.Fprintln(out, "Flags:")
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := &tint.Options{
		Level:   slog.LevelInfo,
		NoColor: !isatty.IsTerminal(os.Stderr.Fd()),
	}
	if verbose {
		opts.Level = slog.LevelDebug
	}

	logger := slog.New(tint.NewHandler(os.Stderr, opts))
	slog.SetDefault(logger)

	switch flag.Arg(0) {
	case "import-users":
		if err := importUsers(ctx, logger); err != nil {
			logger.Error("Failed to import users", "error", err)
			os.Exit(1)
		}
	case "serve":
		fallthrough
	default:
		if err := startServer(ctx, logger); err != nil {
			logger.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}
}

func importUsers(ctx context.Context, logger *slog.Logger) error {
	file := flag.Arg(1)
	if file == "" {
		return fmt.Errorf("no file provided")
	}

	db, err := openDB()
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)

	header, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	expectedHeader := []string{"Name", "School", "Email", "Food", "Major"}
	type Record struct {
		Name  string
		_     string
		Email string
	}

	if !slices.Equal(header, expectedHeader) {
		return fmt.Errorf("unexpected header %v, expecting %v", header, expectedHeader)
	}

	records, err := xcsv.Unmarshal[Record](csvReader)
	if err != nil {
		return fmt.Errorf("failed to unmarshal records: %w", err)
	}

	var errored bool
	err = db.Tx(func(q *sqldb.Queries) error {
		for _, record := range records {
			if err := q.AddUser(ctx, sqldb.AddUserParams{
				UUID:  sqldb.GenerateUUID(),
				Name:  record.Name,
				Email: record.Email,
			}); err != nil {
				slog.Error(
					"Failed to insert user",
					"error", err,
					"user.name", record.Name,
					"user.email", record.Email)
				errored = true
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to add users: %w", err)
	}
	if errored {
		return fmt.Errorf("failed to add some users")
	}

	return nil
}

func startServer(ctx context.Context, logger *slog.Logger) error {
	db, err := openDB()
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	handler := server.NewHandler(db, logger.With("component", "server"))

	logger.Info("Starting server", "addr", addr)
	return hserve.ListenAndServe(ctx, addr, handler)
}

func openDB() (*sqldb.Database, error) {
	slog.Debug(
		"Opening database",
		"path", db)

	return sqldb.Open(db)
}
