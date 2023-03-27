package command

import (
	"flag"
	"fmt"
	flags "github.com/da-moon/dare-cli/cmd/dare/flags"
	model "github.com/da-moon/dare-cli/model"
	hashsink "github.com/da-moon/dare-cli/pkg/hashsink"
	jsonutil "github.com/da-moon/dare-cli/pkg/jsonutil"
	cli "github.com/mitchellh/cli"
	stacktrace "github.com/palantir/stacktrace"
	"math"
	mathrand "math/rand"
	"os"
	"strconv"
	"strings"
)

// DDCommand is a Command implementation that generates an encryption
// key.
type DDCommand struct {
	args []string
	Ui   cli.Ui
}

var _ cli.Command = &DDCommand{}

// Run ...
func (c *DDCommand) Run(args []string) int {
	c.Ui = &cli.PrefixedUi{
		OutputPrefix: "==> ",
		Ui:           c.Ui,
	}

	c.args = args
	const entrypoint = "dd"
	cmdFlags := flag.NewFlagSet(entrypoint, flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Info(c.Help()) }
	sizeString := flags.DDSizeFlag(cmdFlags)
	pathString := flags.DDPathFlag(cmdFlags)
	err := cmdFlags.Parse(c.args)
	if err != nil {
		return 1
	}
	if len(*sizeString) == 0 {
		c.Ui.Error("[ERROR] size value is needed")
		c.Ui.Info(c.Help())
		return 1
	}
	if len(*pathString) == 0 {
		c.Ui.Error("[ERROR] path value is needed")
		c.Ui.Info(c.Help())
		return 1
	}
	parsedSize, err := parseSize(*sizeString)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("[ERROR] could not parse given size: %s", err))
		return 1
	}
	os.Remove(*pathString)
	result, err := createRandomFile(*pathString, int(parsedSize))
	if err != nil {
		c.Ui.Error(fmt.Sprintf("[ERROR] could not create random file: %s", err))
		return 1
	}
	if result == nil {
		c.Ui.Error("[ERROR] could not calculate generated file hashes")
		return 1
	}
	c.Ui.Output(fmt.Sprintf("output path : %s", *pathString))
	c.Ui.Output(fmt.Sprintf("MD5 Hash : %s", result.GetMd5()))
	c.Ui.Output(fmt.Sprintf("SHA256 Hash : %s", result.GetSha256()))
	return 0
}

// Synopsis ...
func (c *DDCommand) Synopsis() string {
	return "Generates a new file used for testing"
}

// Help ...
func (c *DDCommand) Help() string {
	helpText := `
Usage: dare dd [options]

  generates a new human readable JSON lorem ipsum file. 

Options:

  --size=1MB file size to generate.
  --path=/tmp/plain target path to store the file.
`
	return strings.TrimSpace(helpText)
}

// parseSize ...
func parseSize(s string) (int64, error) {
	ss := []byte(strings.ToUpper(s))
	if !(strings.Contains(string(ss), "K") || strings.Contains(string(ss), "KB") ||
		strings.Contains(string(ss), "M") || strings.Contains(string(ss), "MB") ||
		strings.Contains(string(ss), "G") || strings.Contains(string(ss), "GB") ||
		strings.Contains(string(ss), "T") || strings.Contains(string(ss), "TB")) {
		return -1, stacktrace.NewError("wrong format for input string")
	}

	var unit int64 = 1
	p, _ := strconv.Atoi(string(ss[:len(ss)-1]))
	unitstr := string(ss[len(ss)-1])

	if ss[len(ss)-1] == 'B' {
		p, _ = strconv.Atoi(string(ss[:len(ss)-2]))
		unitstr = string(ss[len(ss)-2:])
	}

	switch unitstr {
	default:
		fallthrough
	case "T", "TB":
		unit *= 1024
		fallthrough
	case "G", "GB":
		unit *= 1024
		fallthrough
	case "M", "MB":
		unit *= 1024
		fallthrough
	case "K", "KB":
		unit *= 1024
	}
	return int64(p) * unit, nil
}

func createRandomFile(path string, maxSize int) (*model.Hash, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		err = stacktrace.Propagate(err, "Can't open %s for writing", path)
		return nil, err
	}
	defer file.Close()
	hashWriter := hashsink.NewWriter(file)
	size := maxSize/2 + mathrand.Int()%(maxSize/2)
	loremString := path + `---Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin facilisis mi sapien, vitae accumsan libero malesuada in. Suspendisse sodales finibus sagittis. Proin et augue vitae dui scelerisque imperdiet. Suspendisse et pulvinar libero. Vestibulum id porttitor augue. Vivamus lobortis lacus et libero ultricies accumsan. Donec non feugiat enim, nec tempus nunc. Mauris rutrum, diam euismod elementum ultricies, purus tellus faucibus augue, sit amet tristique diam purus eu arcu. Integer elementum urna non justo fringilla fermentum. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Quisque sollicitudin elit in metus imperdiet, et gravida tortor hendrerit. In volutpat tellus quis sapien rutrum, sit amet cursus augue ultricies. Morbi tincidunt arcu id commodo mollis. Aliquam laoreet purus sed justo pulvinar, quis porta risus lobortis. In commodo leo id porta mattis.`
	byteSizeOfDefaultLorem := len([]byte(loremString))
	repetitions := int(math.Round(float64(size / byteSizeOfDefaultLorem)))
	for i := 0; i < repetitions; i++ {
		enc, _ := jsonutil.EncodeJSONWithIndentation(map[int]string{
			i: (loremString),
		})
		hashWriter.Write([]byte(enc))
	}
	result := &model.Hash{
		Md5:    hashWriter.MD5HexString(),
		Sha256: hashWriter.SHA256HexString(),
	}
	return result, nil
}
