package src

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/zeebo/blake3"
	"golang.org/x/exp/rand"
	"golang.org/x/sync/errgroup"
)

var (
	hasher     = blake3.New()
	tabSpaceRe = regexp.MustCompile(`\t+|\s+`)
)

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

func ShortenTabSpaces(s []byte) []byte {
	return tabSpaceRe.ReplaceAll(s, []byte(" "))
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

func ParallelGroup(parallel int) *errgroup.Group {
	g := errgroup.Group{}
	if parallel > 1 {
		g.SetLimit(parallel)
	}
	return &g
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logrus.Debugln("local ip not found, get net interface failed")
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				logrus.Debugln("found local ip:", ip)
				return ip
			}
		}
	}
	logrus.Debugln("local ip not found")
	return ""
}

func Confirm(msg string) bool {
	prompt := promptui.Prompt{
		Label:     msg,
		IsConfirm: true,
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

func RandStr(length int) string {
	b := make([]byte, length+2)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}

func hashstr(h *blake3.Hasher, s string) [32]byte {
	_, _ = h.WriteString(s)
	result := h.Sum(nil)
	h.Reset()
	return [32]byte(result)
}

func hash(h *blake3.Hasher, b []byte) [32]byte {
	_, _ = h.Write(b)
	result := h.Sum(nil)
	h.Reset()
	return [32]byte(result)
}
