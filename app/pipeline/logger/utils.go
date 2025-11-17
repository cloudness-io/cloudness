package logger

import (
	"bufio"
	"errors"
	"io"
	"sync"

	"github.com/rs/zerolog"
)

func Log(zlog *zerolog.Logger, logger *Logger, rc io.ReadCloser) error {
	defer rc.Close()
	defer logger.Close()

	var uploads sync.WaitGroup
	uploads.Go(func() {
		r := bufio.NewReader(rc)
		for {
			line, err := r.ReadBytes('\n')
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrClosedPipe) {
				logger.Write(line)
				return
			} else if err != nil {
				zlog.Error().Err(err).Msg("logger: error reading line")
				return
			}
			err = logger.Write(line)
			if err != nil {
				zlog.Error().Err(err).Msg("logger: error writing line")
				return
			}
		}
	})

	uploads.Wait()
	return nil
}
