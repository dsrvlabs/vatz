package utils

import (
	"bytes"
	"fmt"
	"net/url"
	"os/exec"
)

func DownloadFileWithWgetOrCurl(url string) ([]byte, error) {
	// First, try to download using wget
	cmd := exec.Command("wget", "-qO-", url)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err == nil {
		return out.Bytes(), nil
	}

	// If wget fails, try using curl
	cmd = exec.Command("curl", "-s", url)
	out.Reset()
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to download file using wget or curl: %v", err)
	}

	return out.Bytes(), nil
}

func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
