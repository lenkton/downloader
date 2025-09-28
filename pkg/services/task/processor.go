package task

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/lenkton/downloader/pkg/fileutil"
)

type processor struct {
}

func newProcessor() *processor {
	return &processor{}
}

// TODO: wait for goroutines to finish
func (p *processor) Start(t *task) {
	t.status = started
	go downloadLinks(t)
}

func downloadLinks(t *task) {
	taskDir := filepath.Join(DownloadDir, strconv.Itoa(t.id))
	fileutil.EnsureDir(taskDir)

	for i, link := range t.links {
		fp := filepath.Join(taskDir, strconv.Itoa(i))
		err := downloadLink(link, fp)
		if err != nil {
			log.Printf("ERROR: downloading file #%v from task %v: %v\n", i, t.id, err)
		}
	}
	t.status = finished
}

func downloadLink(link string, filepath string) error {
	// NOTE: maybe we should create the file before downloading it
	// TODO: check the response code (what if it is 404?)
	resp, err := http.Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	log.Printf("INFO: file %v downloaded successfully\n", filepath)
	return nil
}
