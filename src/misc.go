package src

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/gogs/chardet"
	"github.com/manifoldco/promptui"
	"github.com/samber/lo"
	"github.com/zeebo/blake3"
	"golang.org/x/exp/rand"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"gopkg.in/yaml.v3"
)

var (
	NumberRe = regexp.MustCompile(`\d+`)
)

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

func ExpandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		dirname, _ := os.UserHomeDir()
		path = filepath.Join(dirname, path[2:])
	}
	return path
}

func WriteFile(path string, content string) error {
	// append newline if not exists
	b := []byte(content)
	if b[len(b)-1] != '\n' {
		b = append(b, '\n')
	}
	return os.WriteFile(path, b, 0600)
}

func ReadFileOrStdin(path string) (string, error) {
	var (
		input []byte
		err   error
	)
	switch path {
	case "-":
		// read from stdin
		input, err = io.ReadAll(os.Stdin)
	default:
		input, err = os.ReadFile(path)
	}
	return string(input), err
}

func ParallelGroup(parallel int) *errgroup.Group {
	g := errgroup.Group{}
	if parallel >= 1 {
		g.SetLimit(parallel)
	}
	return &g
}

func Confirm(msg string) bool {
	prompt := promptui.Prompt{
		Label:     msg,
		IsConfirm: true,
	}
	defaultYes := os.Getenv("DORIS_YES")
	if defaultYes == "0" {
		prompt.Stdin = io.NopCloser(bytes.NewReader([]byte("N")))
	} else if defaultYes != "" {
		prompt.Stdin = io.NopCloser(bytes.NewReader([]byte("y")))
	}
	result, _ := prompt.Run()
	return result == "y"
}

func Choose(msg string, items []string) (string, error) {
	prompt := promptui.Select{
		Label:             msg,
		Items:             items,
		Size:              20,
		StartInSearchMode: true,
		Searcher: func(input string, index int) bool {
			item := items[index]
			return strings.Contains(item, input)
		},
	}
	_, result, err := prompt.Run()
	return result, err
}

func hashstr(h *blake3.Hasher, s string) [32]byte {
	_, _ = h.WriteString(s)
	result := h.Sum(nil)
	h.Reset()
	return [32]byte(result)
}

func DetectCharset(r *bufio.Reader) (string, error) {
	hdr, err := r.Peek(4096)
	if len(hdr) == 0 {
		return "", fmt.Errorf("cannot read file: %v", err)
	}
	ress, err := chardet.NewTextDetector().DetectAll(hdr)
	if err != nil {
		return "", fmt.Errorf("cannot detect encoding: %v", err)
	}
	if _, utf8 := lo.Find(ress, func(r chardet.Result) bool { return r.Charset == "UTF-8" }); utf8 {
		return "UTF-8", nil
	}

	return ress[0].Charset, nil
}

func FileGlob(paths []string) ([]string, error) {
	files := []string{}
	for _, s := range paths {
		// '-' represents stdin
		if s == "-" {
			files = append(files, "-")
			continue
		}
		localPaths, err := filepath.Glob(s)
		if err != nil {
			return nil, fmt.Errorf("invalid file path: %s, err: %v", s, err)
		}

		files = append(files, localPaths...)
	}

	return lo.Uniq(files), nil
}

func GetEncoding(name string) (encoding.Encoding, error) {
	enc, err := htmlindex.Get(name)
	if err != nil {
		return nil, fmt.Errorf("invalid encoding: %s", name)
	}
	//nolint:revive
	switch enc {
	case simplifiedchinese.GBK:
		enc = simplifiedchinese.GB18030
	}

	return enc, nil
}

type BytesEncoder interface {
	Encode(b []byte) ([]byte, error)
}

func NewBytesEncoder(srcEncoding encoding.Encoding) BytesEncoder {
	if srcEncoding == unicode.UTF8 {
		return &DummyEncoder{}
	}
	return &Utf8Encoder{
		decoder: srcEncoding.NewDecoder(),
		encoder: unicode.UTF8.NewEncoder(),
	}
}

type Utf8Encoder struct {
	decoder *encoding.Decoder
	encoder *encoding.Encoder
}

func (e *Utf8Encoder) Encode(b []byte) ([]byte, error) {
	dec, err := e.decoder.Bytes(b)
	if err != nil {
		return nil, fmt.Errorf("cannot decode: %s , err: %v", string(b), err)
	}

	enc, err := e.encoder.Bytes(dec)
	if err != nil {
		return nil, fmt.Errorf("cannot encode: %s , err: %v", string(b), err)
	}

	return enc, nil
}

type DummyEncoder struct {
}

func (*DummyEncoder) Encode(b []byte) ([]byte, error) {
	return b, nil
}

func MustJsonMarshal(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

func MustYamlMarshal(v any) []byte {
	data, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

func IsStringType(colType string) bool {
	switch colType {
	case "VARCHAR", "CHAR", "TEXT", "STRING":
		return true
	}
	return false
}
