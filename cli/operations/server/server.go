package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudness-io/cloudness/profiler"
	"github.com/cloudness-io/cloudness/types"

	"github.com/alecthomas/kingpin/v2"
	"github.com/joho/godotenv"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type command struct {
	envfile     string
	initializer func(context.Context, *types.Config) (*System, error)
}

func (c *command) run(*kingpin.ParseContext) error {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// load environment variables from file.
	// no error handling needed when file is not present
	_ = godotenv.Load(c.envfile)

	// create the system configuration store by loading
	// data from the environment.
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("encountered an error while loading configuration: %w", err)
	}

	// configure the log level
	SetupLogger(config)

	// configure profiler
	SetupProfiler(config)

	// add logger to context
	log := log.Logger.With().Logger()
	ctx = log.WithContext(ctx)

	// initialize system
	system, err := c.initializer(ctx, config)
	if err != nil {
		return fmt.Errorf("encountered an error while wiring the system: %w", err)
	}

	// bootstrap the system
	err = system.bootstrap(ctx)
	if err != nil {
		return fmt.Errorf("encountered an error while bootstrapping the system: %w", err)
	}

	// gCtx is canceled if any of the following occurs:
	// - any go routine launched with g encounters an error
	// - ctx is canceled
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := system.services.Cleanup.Register(gCtx); err != nil {
			log.Error().Err(err).Msg("failed to register cleanup service")
			return err
		}

		return system.services.JobScheduler.Run(gCtx)
	})

	// start server
	gHTTP, shutdownHTTP := system.server.ListenAndServe()
	g.Go(gHTTP.Wait)
	//start runner agent for CI deployments
	g.Go(func() error {
		return system.agent.Start(gCtx)
	})
	// if c.enableCI {
	// 	// start poller for CI build executions.
	// 	// g.Go(func() error {
	// 	// 	system.poller.Poll(
	// 	// 		logger.WithWrappedZerolog(ctx),
	// 	// 		config.CI.ParallelWorkers,
	// 	// 	)
	// 	// 	return nil
	// 	// })
	// }

	log.Info().
		Int("port", config.Server.HTTP.Port).
		Msg("server started")

	// wait until the error group context is done
	<-gCtx.Done()

	// restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Info().Msg("shutting down gracefully (press Ctrl+C again to force)")

	// shutdown servers gracefully
	shutdownCtx, cancel := context.WithTimeout(context.Background(), config.GracefulShutdownTime)
	defer cancel()

	if sErr := shutdownHTTP(shutdownCtx); sErr != nil {
		log.Err(sErr).Msg("failed to shutdown http server gracefully")
	}

	log.Info().Msg("Waiting for subroutines to complete")
	err = g.Wait()

	return err
}

// SetupLogger configures the global logger from the loaded configuration.
func SetupLogger(config *types.Config) {
	// configure the log level
	switch {
	case config.Trace:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case config.Debug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// configure time format (ignored if running in terminal)
	zerolog.TimeFieldFormat = time.RFC3339Nano

	// if the terminal is a tty we should output the
	// logs in pretty format
	if isatty.IsTerminal(os.Stdout.Fd()) {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:        os.Stderr,
				NoColor:    false,
				TimeFormat: "15:04:05.999",
			},
		)
	}
}

func SetupProfiler(config *types.Config) {
	profilerType, parsed := profiler.ParseType(config.Profiler.Type)
	if !parsed {
		log.Info().Msgf("No valid profiler so skipping profiling ['%s']", config.Profiler.Type)
		return
	}

	platformProfiler, _ := profiler.New(profilerType)
	platformProfiler.StartProfiling(config.Profiler.ServiceName, "1.0")
}

func Register(app *kingpin.Application, initializer func(context.Context, *types.Config) (*System, error)) {
	c := new(command)
	c.initializer = initializer

	cmd := app.Command("server", "starts the server").
		Action(c.run)

	cmd.Arg("envfile", "load the environment variable file").
		Default("").
		StringVar(&c.envfile)

}
