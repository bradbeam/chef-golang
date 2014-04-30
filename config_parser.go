package chef

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func ParseConfig(configfile ...string) (*KnifeConfig, error) {
	var knifeFile string

	if len(configfile) == 1 && configfile[0] != "" {
		knifeFile = configfile[0]
	}

	// If knifeFile is not provided during initialization
	// iterate through the same was as the knife command
	if knifeFile == "" {
		knifeFiles := []string{}

		// Check current directory for .chef
		knifeFiles = append(knifeFiles, filepath.Join(".chef", "knife.rb"))
		// Chef ~/.chef
		knifeFiles = append(knifeFiles, filepath.Join(os.Getenv("HOME"), ".chef", "knife.rb"))

		for _, each := range knifeFiles {
			if _, err := os.Stat(each); err == nil {
				knifeFile = each
				break
			}
		}

		if knifeFile == "" {
			return nil, errors.New("knife.rb configuration file not found")
		}
	}

	file, err := os.Open(knifeFile)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	config := new(KnifeConfig)
	for scanner.Scan() {
		split := splitWhitespace(scanner.Text())
		if len(split) == 2 {
			switch split[0] {
			case "chef_server_url":
				config.ChefServerUrl = filterQuotes(split[1])
				chefUrl, err := url.Parse(config.ChefServerUrl)
				if err != nil {
					return nil, err
				}
				hostPort := strings.Split(chefUrl.Host, ":")
				if len(hostPort) == 2 {
					config.Host = hostPort[0]
					config.Port = hostPort[1]
				} else if len(hostPort) == 1 {
					config.Host = hostPort[0]
					switch chefUrl.Scheme {
					case "http":
						config.Port = "80"
					case "https":
						config.Port = "443"
					default:
						return nil, errors.New("Invalid http scheme")
					}

				} else {
					return nil, errors.New("Invalid host format")
				}
			case "client_key":
				key, err := KeyFromFile(filterQuotes(split[1]))
				if err != nil {
					return nil, err
				}
				config.ClientKey = key
			case "cookbook_copyright":
				config.CookbookCopyright = filterQuotes(split[1])
			case "cookbook_email":
				config.CookbookEmail = filterQuotes(split[1])
			case "cookbook_license":
				config.CookbookLicense = filterQuotes(split[1])
			case "data_bag_encrypt_version":
				config.DataBagEncryptVersion, _ = strconv.ParseInt(filterQuotes(split[1]), 0, 0)
			case "local_mode":
				config.LocalMode, _ = strconv.ParseBool(filterQuotes(split[1]))
			case "node_name":
				config.NodeName = filterQuotes(split[1])
			case "syntax_check_cache_path":
				config.SyntaxCheckCachePath = filterQuotes(split[1])
			case "validation_client_name":
				config.ValidationClientName = filterQuotes(split[1])
			case "validation_key":
				config.ValidationKey = filterQuotes(split[1])
			case "versioned_cookbooks":
				config.VersionedCookbooks, _ = strconv.ParseBool(filterQuotes(split[1]))
			}
		}
	}

	return config, nil
}

// Given a string with multiple consecutive spaces, splitWhitespace returns a
// slice of strings which represent the given string split by \s characters with
// all duplicates removed
func splitWhitespace(s string) []string {
	re := regexp.MustCompile(`\s+`)
	return strings.Split(re.ReplaceAllString(s, `\s`), `\s`)
}

// filterQuotes returns a string with surrounding quotes filtered
func filterQuotes(s string) string {
	re1 := regexp.MustCompile(`^(\'|\")`)
	re2 := regexp.MustCompile(`(\'|\")$`)
	return re2.ReplaceAllString(re1.ReplaceAllString(s, ``), ``)
}

// KeyFromFile reads an RSA private key given a filepath
func KeyFromFile(filename string) (*rsa.PrivateKey, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return KeyFromString(content)
}

// KeyFromString parses an RSA private key from a string
func KeyFromString(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, fmt.Errorf("block size invalid for '%s'", string(key))
	}
	rsaKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsaKey, nil
}
