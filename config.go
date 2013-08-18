package kelvin

import (
    "path/filepath"
    "github.com/glacjay/goini"
    "github.com/mrallen1/aws4"
    "fmt"
    "os"
    "strings"
    "encoding/base64"
    "log"
)

type KelvinCfg struct {
    aws_keys *aws4.Keys
    aws_service *aws4.Service
    aws_account_id string
    kelvin_aes_key []byte
    kelvin_default_vault string
}

func make_config_string(ak, sk, rg, bk string) string {
    return fmt.Sprintf("[aws]\naccess_key = \"%s\"\nsecret_key = \"%s\"\nregion = \"%s\"\naccount_id = \"-\"\n\n[kelvin]\naes_key = \"%s\"\n",
            ak, sk, rg, bk)
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

func GetConfig(config_file string) KelvinCfg {
    cfg, err := ini.Load(config_file)

    if err != nil {
        if os.IsNotExist(err) {
            if err := setupNewConfig(config_file); err != nil {
                log.Fatal("Error setting up new configuration file: ", err)
            }
        }
        cfg, err = ini.Load(config_file)
        if err != nil {
            log.Fatal("Couldn't load config file: ", config_file)
        }
    }

    access_key, ok := cfg.GetString("aws", "access_key")
    if !ok {
        log.Fatal("Couldn't find aws/access_key in ", config_file)
    }

    secret_key, ok := cfg.GetString("aws", "secret_key")
    if !ok {
        log.Fatal("Couldn't find aws/secret_key in ", config_file)
    }

    k := &aws4.Keys{
        AccessKey: access_key,
        SecretKey: secret_key,
    }

    region, ok := cfg.GetString("aws", "region")
    if !ok {
        log.Fatal("Couldn't find aws/region in ", config_file)
    }

    s := &aws4.Service{
        Name: "glacier",
        Region: region,
    }

    account_id, ok := cfg.GetString("aws", "account_id")
    if !ok {
        log.Fatal("Couldn't find aws/account_id in ", config_file)
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
        log.Fatal("Couldn't find kelvin/aes_key in ", config_file)
    }

    default_vault, ok := cfg.GetString("kelvin", "default_vault")
    if !ok {
        log.Fatal("Couldn't find kelvin/default_vault in ", config_file)
    }

    return KelvinCfg{
        aws_keys: k,
        aws_service: s,
        aws_account_id: account_id,
        kelvin_aes_key: aes_key_bytes,
        kelvin_default_vault: default_vault,
    }

}

func GetConfigLocation() string {
    return filepath.Join(os.Getenv("HOME"), ".kelvin.ini")
}

func setupNewConfig(config_file string) error {
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
                return err
            }

            aes_key := make([]byte, 32)
            _, err = f.Read(aes_key)
            if err != nil {
                return err
            }

            enc := base64.StdEncoding
            b64_key = enc.EncodeToString(aes_key)
        }

        f, err := os.Create(config_file)
        if err != nil {
            return err
        }
        defer f.Close()

        _, err = f.WriteString(make_config_string(access_key, secret_key, region, b64_key))
        if err != nil {
            return err
        }


    } else {
        return fmt.Errorf("Exiting. No configuration file created.")
    }

    return nil
}
