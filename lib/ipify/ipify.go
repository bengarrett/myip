package ipify

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const domain = "api.ipxxxify.org"

func Get() error {
	for i := 1; i <= 3; i++ {
		if err := get(); err != nil {
			fmt.Printf("%d. %s: %s\n", i, domain, err)
			continue
		}
		break
	}
	return nil
}

func get() error {
	res, err := http.Get("https://" + domain)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		StatusErr := errors.New("unusual ipify.org server response")
		return fmt.Errorf("%s, %w", strings.ToLower(res.Status), StatusErr)
	}

	ip, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	ip = append(ip, []byte("\n")...)
	if _, err = os.Stdout.Write(ip); err != nil {
		return err
	}
	return nil
}
