package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/34South/envr"
	dropbox "github.com/tj/go-dropbox"
	dropy "github.com/tj/go-dropy"
)

const fetchLatestSnapshotURL = "https://www.autobus.io/api/snapshots/latest?token=%s"

func init() {
	envr.New("backupdb", []string{
		"AUTOBUS_API_TOKEN",
		"DROPBOX_ACCESS_TOKEN",
	}).Auto()
}

func main() {
	fmt.Println("Fetch latest database snapshot...")
	resp, err := http.Get(fmt.Sprintf(fetchLatestSnapshotURL, os.Getenv("AUTOBUS_API_TOKEN")))
	if err != nil {
		log.Fatalf("http.Get() err = %s\n", err)
	}

	xb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll() err = %s\n", err)
	}
	snapshotURL := string(xb)

	filename := fmt.Sprintf("%v.sql.gz", time.Now().Unix())
	err = DownloadFile(filename, snapshotURL)
	if err != nil {
		log.Fatalf("DownloadFile() err = %s\n", err)
	}
	fmt.Printf("Downloaded snapshot to %s\n", filename)

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("os.Open() err = %s\n", err)
	}
	defer f.Close()

	fmt.Printf("Copying %s to Dropbox...\n", filename)
	client := dropy.New(dropbox.New(dropbox.NewConfig(os.Getenv("DROPBOX_ACCESS_TOKEN"))))
	err = client.Upload("/"+filename, f)
	if err != nil {
		log.Fatalf("client.Upload() err = %s\n", err)
	}

	fmt.Println("Cleaning up...")
	err = os.Remove(filename)
	if err != nil {
		fmt.Printf("os.Remove() err = %s\n", err)
	}

	fmt.Println("Done!")
}

// Code below pinched from here: https://golangcode.com/download-a-file-with-progress/

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
