package swot

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_IsAcademic(t *testing.T) {
	tests := map[string]bool{
		"lreilly@stanford.edu":          true,
		"LREILLY@STANFORD.EDU":          true,
		"Lreilly@Stanford.Edu":          true,
		"lreilly@slac.stanford.edu":     true,
		"lreilly@strath.ac.uk":          true,
		"lreilly@soft-eng.strath.ac.uk": true,
		"lee@ugr.es":                    true,
		"lee@uottawa.ca":                true,
		"lee@mother.edu.ru":             true,
		"lee@ucy.ac.cy":                 true,
		"lee@leerilly.net":              false,
		"lee@gmail.com":                 false,
		"lee@stanford.edu.com":          false,
		"lee@strath.ac.uk.com":          false,
		"stanford.edu":                  true,
		"slac.stanford.edu":             true,
		"www.stanford.edu":              true,
		"http://www.stanford.edu":       true,
		"http://www.stanford.edu:9393":  true,
		"strath.ac.uk":                  true,
		"soft-eng.strath.ac.uk":         true,
		"ugr.es":                        true,
		"uottawa.ca":                    true,
		"mother.edu.ru":                 true,
		"ucy.ac.cy":                     true,
		"leerilly.net":                  false,
		"gmail.com":                     false,
		"stanford.edu.com":              false,
		"strath.ac.uk.com":              false,
		"":                              false,
		"the":                           false,
		" stanford.edu":                 true,
		"lee@strath.ac.uk ":             true,
		" gmail.com":                    false,
		"lee@stud.uni-corvinus.hu":      true,
		"lee@harvard.edu":               true,
		"lee@mail.harvard.edu":          true,
		"imposter@si.edu":               false,
		"lee@acmt.ac.ir":                true,
		"lee@australia.edu":             false,
		"si.edu":                        false,
		"foo.si.edu":                    false,
		"america.edu":                   false,
		"folger.edu":                    false,
		"foo@bar.invalid":               false,
		".com":                          false,
	}

	for key, want := range tests {
		got := IsAcademic(key)
		if got != want {
			t.Fatalf("%s, want:%t, got:%t\n", key, want, got)
		}
	}

}

func Test_GetSchoolName(t *testing.T) {
	tests := map[string]string{
		"lreilly@cs.strath.ac.uk":        "University of Strathclyde",
		"lreilly@fadi.at":                "BRG Fadingerstraße Linz, Austria",
		"abadojack@students.uonbi.ac.ke": "University of Nairobi",
		"foo@shop.com":                   "",
		"bar@gmail.com":                  "",
		"harvard.edu":                    "Harvard University",
		"stanford.edu":                   "Stanford University",
	}

	for key, want := range tests {
		got := GetSchoolName(key)
		if strings.Compare(got, want) != 0 {
			fmt.Println(strings.Compare(got, want))
			t.Fatalf("%s, want:%s, got:%s\n", key, want, got)
		}
	}
}

func Test_domainFiles(t *testing.T) {
	err := filepath.Walk("domains", walkFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func walkFunc(path string, info os.FileInfo, err error) error {
	if !info.IsDir() {
		//Should contain only text files
		if filepath.Ext(info.Name()) != ".txt" {
			return errors.New(info.Name() + " should have a .txt extension.")
		}

		//Each file should contain only a single line of text.
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(f)
		lines := 0
		for scanner.Scan() {
			lines++
		}

		if lines > 1 {
			return errors.New(info.Name() + " should only have a single line of text.")
		}
	}
	return nil
}
