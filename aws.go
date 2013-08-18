package kelvin

import (
//    "fmt"
    "net/http"
    "time"
    "log"
)

func build_host(region string) string {
    return "glacier." + region + ".amazonaws.com"
}

func build_url(region, url, account_id string) string {

    return "https://" + build_host(region) + "/" + account_id + "/" + url

}

func GlacierRequest(operation, url, data string, cfg KelvinCfg) *http.Request {

    aws_url := build_url(cfg.aws_service.Region, url, cfg.aws_account_id)
    r, _ := http.NewRequest(operation, aws_url, nil)
    r.Header.Set("Host", build_host(cfg.aws_service.Region))
    r.Header.Set("x-amz-glacier-version", "2012-06-01")
    r.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))

    if err := cfg.aws_service.Sign(cfg.aws_keys, r); err != nil {
		log.Fatal(err)
	}

    return r

}


