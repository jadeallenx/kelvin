package kelvin

import (
  "path/filepath"
  "os"
  "flag"
  "fmt"
//  "log"
)


func visit(path string, f os.FileInfo, err error) error {
//  fmt.Printf("Visited: %s\n", path)
  return nil
}

/* 
   check for config
   if no config, create config

   else

   walk path
   zip files
   check vault
   if no vault, make one
   upload file to vault
*/


func main() {
  cfg := GetConfig(GetConfigLocation())
  fmt.Printf("config: %v", cfg)

  r := GlacierRequest("GET", "vaults", "", cfg)
  fmt.Println("foo")
  fmt.Println(r.Write(os.Stderr))

  flag.Parse()
  root := flag.Arg(0)

  fmt.Printf("%v", root)

  // make the current directory a default
  if root == "" {
      root = "."
  }

  err := filepath.Walk(root, visit)
  fmt.Printf("filepath.Walk() returned %v\n", err)
}
