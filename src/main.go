package main

import (
  "path/filepath"
  "os"
  "flag"
  "fmt"
)

func visit(path string, f os.FileInfo, err error) error {
  fmt.Printf("Visited: %s\n", path)
  return nil
}


func main() {
  flag.Parse()
  root := flag.Arg(0)
  err := filepath.Walk(root, visit)
  fmt.Printf("filepath.Walk() returned %v\n", err)
}

