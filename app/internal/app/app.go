package app

import (
	"emails_extractor/app/internal/config"
	"emails_extractor/app/internal/extractors/csv"
	"emails_extractor/app/internal/extractors/sql"
	"emails_extractor/app/internal/extractors/txt"
	logger "emails_extractor/app/pkg/logger/zap"
	"emails_extractor/app/pkg/utils"
	"fmt"
	"go.uber.org/zap"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type App struct {
	cfg   *config.Config
	files []string
}

func NewApp(cfg *config.Config) *App {
	if !utils.IsDirectoryExists(cfg.Databases.Input) {
		logger.GetLogger().Warn("directory does not exists. Trying to create it.", zap.String("directory_path", cfg.Databases.Input))

		if err := utils.CreateDir(cfg.Databases.Input); err != nil {
			logger.GetLogger().Fatal("failed to create directory", zap.String("directory_path", cfg.Databases.Input))
		}
	}

	if !utils.IsDirectoryExists(cfg.Databases.Output) {
		logger.GetLogger().Warn("directory does not exists. Trying to create it.", zap.String("directory_path", cfg.Databases.Output))

		if err := utils.CreateDir(cfg.Databases.Output); err != nil {
			logger.GetLogger().Fatal("failed to create directory", zap.String("directory_path", cfg.Databases.Output))
		}
	}

	if !utils.IsDirectoryExists(cfg.Databases.MergedOutput) {
		logger.GetLogger().Warn("directory does not exists. Trying to create it.", zap.String("directory_path", cfg.Databases.MergedOutput))

		if err := utils.CreateDir(cfg.Databases.MergedOutput); err != nil {
			logger.GetLogger().Fatal("failed to create directory", zap.String("directory_path", cfg.Databases.MergedOutput))
		}
	}

	files, err := utils.GetFilesInDirectory(cfg.Databases.Input)
	if err != nil {
		return nil
	}

	return &App{
		cfg:   cfg,
		files: files,
	}
}

func (a App) ExtractEmails() error {
	emailsChan := make(chan []string)
	quitChan := make(chan bool)
	doneChan := make(chan bool)

	mergeAndWrite := func(emailsChan chan []string, quitChan chan bool) {
		var wg sync.WaitGroup
		mergedEmails := make([]string, 0)

		wg.Add(1)
		go func() {
			for {
				select {
				case emails := <-emailsChan:
					mergedEmails = append(mergedEmails, emails...)
					logger.GetLogger().Info("got new emails")
				case <-quitChan:
					logger.GetLogger().Info("sort emails and delete duplicates")
					utils.Sort(mergedEmails)

					logger.GetLogger().Info("write sorted and cleared emails to file")

					file := fmt.Sprintf("%d_all.txt", time.Now().Unix())

					p := path.Join(a.cfg.Databases.MergedOutput, file)

					output, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE, 0644)
					if err != nil {
						close(quitChan)
						close(emailsChan)
						panic("failed to open file")
					}

					_, err = output.Write([]byte(strings.Join(mergedEmails, "\n")))
					if err != nil {
						close(quitChan)
						close(emailsChan)
						panic("failed to write to fail")
					}

					logger.GetLogger().Info("done! Close all the channels and finish my job")
					close(quitChan)
					close(emailsChan)

					wg.Done()
				}
			}
		}()

		wg.Wait()
		doneChan <- true
	}
	go mergeAndWrite(emailsChan, quitChan)

	var wg sync.WaitGroup
	for _, file := range a.files {
		wg.Add(1)
		fileExtension := utils.GetFileExtension(file)

		switch fileExtension {
		case config.CSV:
			csvExtractor := csv.NewExtractor(a.cfg, file)
			go csvExtractor.Start(&wg, emailsChan)

			break
		case config.TXT:
			txtExtractor := txt.NewExtractor(a.cfg, file)
			go txtExtractor.Start(&wg, emailsChan)

			break
		case config.SQL:
			sqlExtractor := sql.NewExtractor(a.cfg, file)
			go sqlExtractor.Start(&wg, emailsChan)

			break
		default:
			logger.GetLogger().Warn("unknown or unsupported file extension", zap.String("extension", fileExtension))
			break
		}
	}

	wg.Wait()
	quitChan <- true

	<-doneChan

	return nil
}
