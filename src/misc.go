package src

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/goccy/go-json"
	"github.com/gogs/chardet"
	"github.com/manifoldco/promptui"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/xyproto/randomstring"
	"github.com/zeebo/blake3"
	"golang.org/x/exp/rand"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
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
			return nil, fmt.Errorf("invalid file path: %s, error: %v", s, err)
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

func (e *DummyEncoder) Encode(b []byte) ([]byte, error) {
	return b, nil
}

func MustJsonMarshal(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

func Cast2[R int8 | int16 | int | int32 | int64 | float32 | float64 | string | time.Time](v1, v2 any) (r1, r2 R, err error) {
	r1, err = Cast[R](v1)
	if err != nil {
		return
	}
	r2, err = Cast[R](v2)
	return
}

func Cast[R int8 | int16 | int | int32 | int64 | float32 | float64 | string | time.Time](v any) (r R, err error) {
	var r_ any

	switch any(r).(type) {
	case int8:
		r_, err = cast.ToInt8E(v)
	case int16:
		r_, err = cast.ToInt16E(v)
	case int:
		r_, err = cast.ToIntE(v)
	case int32:
		r_, err = cast.ToInt32E(v)
	case int64:
		r_, err = cast.ToInt64E(v)
	case float32:
		r_, err = cast.ToFloat32E(v)
	case float64:
		r_, err = cast.ToFloat64E(v)
	case string:
		r_, err = cast.ToInt16E(v)
	case time.Time:
		r_, err = cast.ToTimeE(v)
	default:
		return r, fmt.Errorf("unsupported cast type '%T' to '%T'", v, r)
	}

	return r_.(R), err
}

func RandomStr(lenMin, lenMax int) string {
	length := gofakeit.IntRange(lenMin, lenMax)
	if length < 20 {
		return randomstring.HumanFriendlyString(length)
	}
	return randomstring.CookieFriendlyString(length)
}

func IsStringType(colType string) bool {
	switch colType {
	case "VARCHAR", "CHAR", "TEXT", "STRING":
		return true
	}
	return false
}
