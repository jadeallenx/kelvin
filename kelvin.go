package main

import (
  "path/filepath"
  "os"
  "flag"
  "fmt"
  //"archive/zip"
  "log"
  "strings"
  "encoding/base64"
  //"github.com/bmizerany/aws4"
  //"github.com/glacjay/goini"
)

func visit(path string, f os.FileInfo, err error) error {
  fmt.Printf("Visited: %s\n", path)
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

func prompt(p string) string {
    fmt.Printf("%s\n", p)
    var prompt string
    fmt.Scanf("%s", &prompt)
    return strings.TrimSpace(prompt)
}

func ynprompt(p string) bool {
    response := strings.ToLower(prompt(p + " (y/n) "))

    if strings.HasPrefix(response, "y") {
        return true
    }

    return false
}

func make_config_string(ak, sk, rg, bk string) string {
    return fmt.Sprintf("[aws]\naccess_key = \"%s\"\nsecret_key = \"%s\"\nregion = \"%s\"\naccount_id = \"-\"\n\n[kelvin]\naes_key = \"%s\"\n",
            ak, sk, rg, bk)
}

func main() {

  config_file := filepath.Join(os.Getenv("HOME"), ".kelvin.ini")
  cfg, err := os.Open(config_file)
  if err != nil {
    if os.IsNotExist(err) {
        fmt.Printf("No configuration file found: %s\n", config_file)
        if ynprompt("Would you like to create one?") {
            access_key := prompt("What is your AWS Access Key?")
            secret_key := prompt("What is your AWS Secret Key?")
            region := prompt("What AWS region do you want to use?")

            var b64_key string
            if ynprompt("Do you want to generate an AES key?") {
                // generate key
                f, err := os.Open("/dev/random")
                if err != nil {
                    log.Fatal(err)
                }

                aes_key := make([]byte, 32)
                _, err = f.Read(aes_key)
                if err != nil {
                    log.Fatal(err)
                }

                enc := base64.StdEncoding
                b64_key = enc.EncodeToString(aes_key)
            }

            f, err := os.Create(config_file)
            if err != nil {
                log.Fatal(err)
            }
            defer f.Close()

            _, err = f.WriteString(make_config_string(access_key, secret_key, region, b64_key))
            if err != nil {
                log.Fatal(err)
            }

            os.Exit(0)

        } else {
            fmt.Printf("Exiting. No configuration file created.\n")
            os.Exit(1)
        }
      } else {
        log.Fatal(err)
      }
  }

  defer cfg.Close()

  flag.Parse()
  root := flag.Arg(0)

  err = filepath.Walk(root, visit)
  fmt.Printf("filepath.Walk() returned %v\n", err)
}

