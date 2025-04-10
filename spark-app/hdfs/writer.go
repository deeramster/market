package hdfs

import (
	"fmt"
	"os"
	"time"

	"spark-app/config"

	"github.com/colinmarc/hdfs"
)

type Writer struct {
	client *hdfs.Client
}

func NewWriter() (*Writer, error) {
	client, err := hdfs.New(config.HDFSNamenode())
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к HDFS: %w", err)
	}

	return &Writer{client: client}, nil
}

func (w *Writer) WriteJSONLine(data []byte) error {
	today := time.Now().Format("2006-01-02")
	path := fmt.Sprintf("%s/%s.jsonl", config.HDFSDataPath(), today)

	exists := true
	_, err := w.client.Stat(path)
	if err != nil {
		if pathError, ok := err.(*os.PathError); ok && pathError.Err == os.ErrNotExist {
			exists = false
		} else {
			return fmt.Errorf("ошибка доступа к файлу: %w", err)
		}
	}

	var file *hdfs.FileWriter
	if exists {
		file, err = w.client.Append(path)
	} else {
		_, err := w.client.Stat(config.HDFSDataPath())
		if err != nil {
			if pathError, ok := err.(*os.PathError); ok && pathError.Err == os.ErrNotExist {
				err = w.client.MkdirAll(config.HDFSDataPath(), 0755)
				if err != nil {
					return fmt.Errorf("ошибка создания директории: %w", err)
				}
			} else {
				return fmt.Errorf("ошибка проверки директории: %w", err)
			}
		}

		file, err = w.client.Create(path)
	}
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()

	_, err = file.Write(append(data, '\n'))
	if err != nil {
		return fmt.Errorf("ошибка записи в файл: %w", err)
	}

	return nil
}

func (w *Writer) Close() error {
	return w.client.Close()
}
