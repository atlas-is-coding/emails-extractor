package csv

import (
	"emails_extractor/app/internal/config"
	logger "emails_extractor/app/pkg/logger/zap"
	"encoding/csv"
	"go.uber.org/zap"
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
		logger.GetLogger().Warn("can`e extract data from file", zap.String("file", e.file))
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

	reader := csv.NewReader(input)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	emails := make([]string, 0)
	for _, record := range records {
		r, err := regexp.Compile(config.EmailRegex)
		if err != nil {
			return nil, err
		}

		foundEmails := r.FindAllString(strings.Join(record, " "), -1)

		if len(foundEmails) > 0 {
			emails = append(emails, strings.Join(foundEmails, "\n"))
		}
	}

	e.file = strings.Replace(e.file, ".csv", ".txt", -1)
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
