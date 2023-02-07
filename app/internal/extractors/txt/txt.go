package txt

import (
	"emails_extractor/app/internal/config"
	logger "emails_extractor/app/pkg/logger/zap"
	"go.uber.org/zap"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
)

type Extractor struct {
	cfg  *config.Config
	file string
}

func NewExtractor(cfg *config.Config, file string) *Extractor {
	return &Extractor{
		cfg:  cfg,
		file: file,
	}
}

func (e Extractor) Start(group *sync.WaitGroup, emailsChan chan []string) {
	defer group.Done()
	emails, err := e.startWrapper()
	if err != nil {
		logger.GetLogger().Warn("can`t extract data from file", zap.String("file", e.file))
	}

	emailsChan <- emails
}

func (e Extractor) startWrapper() ([]string, error) {

	p := path.Join(e.cfg.Databases.Input, e.file)

	input, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	defer input.Close()

	var data []byte
	buf := make([]byte, 1024)
	for {
		n, err := input.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		data = append(data, buf[:n]...)
	}

	r, err := regexp.Compile(config.EmailRegex)
	if err != nil {
		return nil, err
	}

	emails := r.FindAllString(string(data), -1)
	p = path.Join(e.cfg.Databases.Output, e.file)

	output, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	_, err = output.Write([]byte(strings.Join(emails, "\n")))
	if err != nil {
		return nil, err
	}

	return emails, nil
}
