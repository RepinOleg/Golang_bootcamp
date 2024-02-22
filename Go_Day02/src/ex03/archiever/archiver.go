package archiever

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func CreateArchive(path, filename string, wg *sync.WaitGroup) {
	//убавляем счетчик для wg
	defer wg.Done()

	/*
		Вызываем функицю которая создает архив
	*/
	tarFile, err := createFile(path, filename)
	if err != nil {
		log.Fatal(err)
	}
	defer tarFile.Close()

	//инициализируем переменные для архивации
	gzipWriter := gzip.NewWriter(tarFile)
	defer gzipWriter.Close()
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Функция записи в tar файл
	err = writeLogToTar(filename, tarWriter)
	if err != nil {
		log.Fatal(err)
	}
}

func createFile(path, filename string) (*os.File, error) {

	//вычисляем разницу во времени между последней модификацией файла и настоящим моментом

	diffTime, err := calculateMTime(filename)
	if err != nil {
		return nil, err
	}

	// Проводим манипуляции с путем к файлу чтобы при создании он был корректный

	filenameWithoutSuffix := strings.Split(filename, filepath.Ext(filename))
	filenameWithoutPreffix := strings.TrimPrefix(filenameWithoutSuffix[0], "/")

	path += fmt.Sprintf("/%s_%s.tar.gz", filenameWithoutPreffix, diffTime)

	//создаем файл
	tarFile, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("could not create tarball file '%s', got error '%s'", path, err.Error())
	}

	//возвращаем указатель на файл
	return tarFile, nil
}

func writeLogToTar(filename string, tw *tar.Writer) error {
	// открываем файл
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	// берем описание файла
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	//Вызываем функцию из библиотеки tar которая создает хэдер файл из описания нашего файла
	header, err := tar.FileInfoHeader(stat, stat.Name())
	if err != nil {
		return err
	}
	// записываем этот хедер в архив
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	// копируем содержимое файла
	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}
	return nil
}

func calculateMTime(file string) (string, error) {
	info, err := os.Stat(file)

	if err != nil {
		return "", err
	}

	//время сейчас
	now := time.Now()
	//время модификации файла
	mtime := info.ModTime()
	//вычитаем, переводим в секунды и округляем
	diff := math.Round(now.Sub(mtime).Seconds())

	return fmt.Sprintf("%d", int(diff)), nil
}
