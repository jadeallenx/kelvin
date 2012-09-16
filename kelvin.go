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
  "github.com/bmizerany/aws4"
  "github.com/glacjay/goini"
  //"crypto/tls"
  "net/http"
  //"io/ioutil"
  "time"
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

type KelvinCfg struct {
    aws_keys *aws4.Keys
    aws_service *aws4.Service
    aws_account_id string
    kelvin_aes_key []byte
    kelvin_default_vault string
}

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

func get_config(cfg ini.Dict, loc string) (c KelvinCfg) {
    access_key, ok := cfg.GetString("aws", "access_key")
    if !ok {
        log.Fatal("Couldn't find aws/access_key in ", loc)
    }

    secret_key, ok := cfg.GetString("aws", "secret_key")
    if !ok {
        log.Fatal("Couldn't find aws/secret_key in ", loc)
    }

    k := &aws4.Keys{
        AccessKey: access_key,
        SecretKey: secret_key,
    }

    region, ok := cfg.GetString("aws", "region")
    if !ok {
        log.Fatal("Couldn't find aws/region in ", loc)
    }

    s := &aws4.Service{
        Name: "glacier",
        Region: region,
    }

    account_id, ok := cfg.GetString("aws", "account_id")
    if !ok {
        log.Fatal("Couldn't find aws/account_id in ", loc)
    }

    var aes_key_bytes []byte
    aes_key_b64, ok := cfg.GetString("kelvin", "aes_key")
    if ok {
        var err error
        enc := base64.StdEncoding
        aes_key_bytes, err = enc.DecodeString(aes_key_b64)
        if err != nil {
            log.Fatal(err)
        }
    } else {
        log.Fatal("Couldn't find kelvin/aes_key in ", loc)
    }

    default_vault, ok := cfg.GetString("kelvin", "default_vault")
    if !ok {
        log.Fatal("Couldn't find kelvin/default_vault in ", loc)
    }

    return KelvinCfg{
        aws_keys: k,
        aws_service: s,
        aws_account_id: account_id,
        kelvin_aes_key: aes_key_bytes,
        kelvin_default_vault: default_vault,
    }

}

func build_host(region string) string {
    return "glacier." + region + ".amazonaws.com"
}

func build_url(region, url, account_id string) string {

    return "https://" + build_host(region) + "/" + account_id + "/" + url

}

func build_glacier_request(operation, url, data string, cfg KelvinCfg) *http.Request {

    url2 := build_url(cfg.aws_service.Region, url, cfg.aws_account_id)
    r, _ := http.NewRequest(operation, url2, nil)
    r.Header.Set("Host", build_host(cfg.aws_service.Region))
    r.Header.Set("x-amz-glacier-version", "2012-06-01")
    r.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))

    fmt.Printf("%v", r)

    if err := cfg.aws_service.Sign(cfg.aws_keys, r); err != nil {
		log.Fatal(err)
	}

    return r

}

func main() {
 /* tr := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
  }

  c := &http.Client{Transport: tr} */

  config_file := filepath.Join(os.Getenv("HOME"), ".kelvin.ini")
  cff, err := ini.Load(config_file)
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
            log.Fatal("Exiting. No configuration file created.")
        }
      } else {
        log.Fatal(err)
      }
  }

  cfg := get_config(cff, config_file)
  fmt.Println(cfg)

  r := build_glacier_request("GET", "vaults", "", cfg)
  fmt.Println("foo")
  fmt.Println(r.Write(os.Stderr))

  flag.Parse()
  root := flag.Arg(0)

  err = filepath.Walk(root, visit)
  fmt.Printf("filepath.Walk() returned %v\n", err)
}

