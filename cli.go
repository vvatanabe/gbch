package gbch

import (
	"fmt"
	"io"
	"log"

	"github.com/jessevdk/go-flags"
	backlog "github.com/vvatanabe/go-backlog/backlog/v2"
)

// Gbch is main application struct
type Gbch struct {
	RepoPath     string `short:"r" long:"repo" default:"." description:"git repository path"`
	GitPath      string `short:"g" long:"git" default:"git" description:"git path"`
	From         string `short:"f" long:"from" description:"git commit revision range start from"`
	To           string `short:"t" long:"to" description:"git commit revision range end to"`
	Latest       bool   `          long:"latest" description:"output changes between latest two semantic versioned tags"`
	APIKey       string `          long:"apikey" description:"backlog api key"`
	Remote       string `          long:"remote" default:"origin" description:"default remote name"`
	Format       string `short:"F" long:"format" default:"json" description:"json or markdown or backlog"`
	All          bool   `short:"A" long:"all" description:"output all changes"`
	NextVersion  string `short:"N" long:"next-version"`
	Write        bool   `short:"w" description:"write result to file"`
	ShowUniqueID bool   `          long:"show-uid" description:"show the unique id on nulab account"`
	VerPrefix    string `          long:"ver-prefix" description:"custom version prefix (ex. foo-1.2.3)"`
	ChangelogMd  string
	// Tmpl string
	OutStream io.Writer

	client     *backlog.Client
	BaseURL    string
	ProjectKey string
	RepoName   string
}

const (
	exitCodeOK = iota
	exitCodeParseFlagError
	exitCodeErr
)

// CLI is struct for command line tool
type CLI struct {
	OutStream, ErrStream io.Writer
}

// Run the gbch
func (cli *CLI) Run(argv []string) int {
	log.SetOutput(cli.ErrStream)
	p, gh, err := cli.parseArgs(argv)
	if err != nil {
		if ferr, ok := err.(*flags.Error); !ok || ferr.Type != flags.ErrHelp {
			p.WriteHelp(cli.ErrStream)
		}
		return exitCodeParseFlagError
	}
	if err := gh.Run(); err != nil {
		log.Println(err)
		return exitCodeErr
	}
	return exitCodeOK
}

func (cli *CLI) parseArgs(args []string) (*flags.Parser, *Gbch, error) {
	gb := &Gbch{
		OutStream: cli.OutStream,
	}
	p := flags.NewParser(gb, flags.Default)
	p.Usage = fmt.Sprintf("[OPTIONS]\n\nVersion: %s (rev: %s)", version, revision)
	rest, err := p.ParseArgs(args)
	if gb.Write {
		if !gb.isMDFormat() {
			gb.Format = "markdown"
		}
		gb.ChangelogMd = "CHANGELOG.md"
		if len(rest) > 0 {
			gb.ChangelogMd = rest[0]
		}
	}
	return p, gb, err
}
