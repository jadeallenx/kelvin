package kelvin

import (
    "fmt"
    "os"
    "io/ioutil"
    "archive/zip"
)

type ZipFile struct {
    name string
    body []byte
}

var ( zip_files = []ZipFile{} )

func Visit(path string, f os.FileInfo, err error) error {
  if f.IsDir() {
      return nil
  }

  body, err := ioutil.ReadFile(path)
  if err != nil {
      return err
  }

  zf := ZipFile{
        name: path,
        body: body,
  }

  zip_files = append(zip_files, zf)

  return nil
}

// FIXME: Should write the file during visit instead of stuffing it all into memory
func WriteZipFile(archive *os.File) (string, error) {
    w := zip.NewWriter(archive)

    for _, file := range zip_files {
        f, err := w.Create(file.name)
        if err != nil {
            return "", err
        }
        _, err = f.Write(file.body)
        if err != nil {
            return "", err
        }
    }

    // Make sure to check the error on Close.
    err := w.Close()
    if err != nil {
        return "", err
    }
    return archive.Name(), nil
}

func DumpZipFiles() error {
  fmt.Printf("%v", zip_files)
  return nil
}

