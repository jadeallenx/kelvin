package kelvin

import (
    "fmt"
    "os"
)

func Visit(path string, f os.FileInfo, err error) error {
  fmt.Printf("Visited: %s\n", path)
  return nil
}


