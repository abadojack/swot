package swot

import (
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	valid "github.com/asaskevich/govalidator"
)

//IsAcademic returns true if the email address/URL belongs to an academic institution.
func IsAcademic(emailOrURL string) bool {
	domainName, err := getDomainName(emailOrURL)
	if err != nil {
		return false
	}

	if isBlacklisted(domainName) {
		return false
	} else if isAcademicTLD(domainName) {
		return true
	}

	_, err = getInstitutionName(domainName)
	return err == nil
}

func getDomainName(emailOrURL string) (string, error) {
	var domainName string

	emailOrURL = strings.ToLower(strings.TrimSpace(emailOrURL))

	if valid.IsEmail(emailOrURL) {
		return strings.Split(emailOrURL, "@")[1], nil
	} else if valid.IsURL(emailOrURL) {
		if valid.IsRequestURL(emailOrURL) {
			url, err := url.Parse(emailOrURL)
			if err != nil {
				return "", err
			}

			domainName = url.Host
			domainName = strings.Split(domainName, ":")[0]
		} else {
			domainName = emailOrURL
		}

		domainName = strings.TrimPrefix(domainName, "www")

		return domainName, nil
	}

	return "", errors.New("Domain name not found.")
}

func isBlacklisted(domainName string) bool {
	for _, dn := range blacklist {
		if strings.HasSuffix(domainName, dn) {
			return true
		}
	}
	return false
}

//isAcademicTLD returns true if the domainName is a top level academic domain
//or false otherwise.
func isAcademicTLD(domainName string) bool {
	for _, tld := range academicTLDs {
		if strings.HasSuffix(domainName, tld) {
			return true
		}
	}
	return false
}

//fileExits returns true if the file exists or false otherwise.
func fileExits(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

//GetSchoolName returns the name of the academic institution
//or an empty string if the name of the institution is not found.
func GetSchoolName(emailOrURL string) string {
	domainName, err := getDomainName(emailOrURL)
	if err != nil {
		return ""
	}

	s, err := getInstitutionName(domainName)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(s)
}

func getInstitutionName(domainName string) (string, error) {
	domainParts := splitDomainName(domainName)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("Path to 'domains' directory not found")
	}

	path := filepath.Join(filepath.Dir(filename), "domains", domainParts[len(domainParts)-1])
	for i := len(domainParts) - 2; i >= 0; i-- {
		path = filepath.Join(path, domainParts[i])
		if fileExits(path + ".txt") {
			b, err := ioutil.ReadFile(path + ".txt")
			if err != nil {
				return "", err
			}

			return string(b), nil
		}
	}
	return "", errors.New("Name of school not found")
}

//splitDomainName splits the domain name at the dots and returns a string array
//of the split parts.
//For example: uonbi.ac.ke ==> [uonbi ac ke]
func splitDomainName(domainName string) []string {
	return strings.Split(domainName, ".")
}
