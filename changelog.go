package gbch

import (
	"bufio"
	"bytes"
	"context"
	"log"
	"strings"
	"text/template"
	"time"

	"fmt"

	backlog "github.com/vvatanabe/go-backlog/backlog/v2"
)

// Changelog contains Sections
type Changelog struct {
	Sections []Section `json:"Sections"`
}

func insertNewChangelog(orig []byte, section, header string) string {
	var bf bytes.Buffer
	lineSnr := bufio.NewScanner(bytes.NewReader(orig))
	inserted := false
	headerL2 := header + header + " "
	for lineSnr.Scan() {
		line := lineSnr.Text()
		if !inserted && strings.HasPrefix(line, headerL2) {
			bf.WriteString(section)
			bf.WriteString("\n\n")
			inserted = true
		}
		bf.WriteString(line)
		bf.WriteString("\n")
	}
	if !inserted {
		bf.WriteString(section)
	}
	return bf.String()
}

// Section contains changes between two revisions
type Section struct {
	PullRequests []*backlog.PullRequest `json:"pull_requests"`
	FromRevision string                 `json:"from_revision"`
	ToRevision   string                 `json:"to_revision"`
	ChangedAt    time.Time              `json:"changed_at"`
	Project      string                 `json:"project"`
	Repo         string                 `json:"repo"`
	BaseURL      string                 `json:"-"`
	HTMLURL      string                 `json:"html_url"`
	ShowUniqueID bool                   `json:"-"`
}

var markdownTmplStr = `{{$ret := . -}}
## [{{.ToRevision}}]({{.HTMLURL}}/compare/{{.FromRevision}}...{{.ToRevision}}) ({{.ChangedAt.Format "2006-01-02"}})
{{range .PullRequests}}
* {{.Summary}} [#{{.Number}}]({{$ret.HTMLURL}}/pullRequests/{{.Number}}) ([{{.CreatedUser.Name}}]({{$ret.BaseURL}}/user/{{.CreatedUser.UserID}})){{if and ($ret.ShowUniqueID) (.CreatedUser.NulabAccount)}} @{{.CreatedUser.NulabAccount.UniqueID}}{{end}}
{{- end}}`

var backlogTmplStr = `{{$ret := . -}}
** [[{{.ToRevision}}:{{.HTMLURL}}/compare/{{.FromRevision}}...{{.ToRevision}}]] ({{.ChangedAt.Format "2006-01-02"}}){{range .PullRequests}}
- {{.Summary}} [[[#{{.Number}}:{{$ret.HTMLURL}}/pullRequests/{{.Number}}]]] ([[{{.CreatedUser.Name}}:{{$ret.BaseURL}}/user/{{.CreatedUser.UserID}}]]){{if and ($ret.ShowUniqueID) (.CreatedUser.NulabAccount)}} @{{.CreatedUser.NulabAccount.UniqueID}}{{end}}
{{- end}}`

var (
	mdTmpl *template.Template
	blTmpl *template.Template
)

func init() {
	var err error
	mdTmpl, err = template.New("md-changelog").Parse(markdownTmplStr)
	if err != nil {
		log.Fatal(err)
	}
	blTmpl, err = template.New("bl-changelog").Parse(backlogTmplStr)
	if err != nil {
		log.Fatal(err)
	}
}

func (rs Section) toMkdn() (string, error) {
	var b bytes.Buffer
	err := mdTmpl.Execute(&b, rs)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func (rs Section) toBacklog() (string, error) {
	var b bytes.Buffer
	err := blTmpl.Execute(&b, rs)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func (gb *Gbch) getSection(ctx context.Context, from, to string) (Section, error) {
	if from == "" {
		from, _ = gb.cmd("rev-list", "--max-parents=0", "HEAD")
		from = strings.TrimSpace(from)
		if len(from) > 12 {
			from = from[:12]
		}
	}
	r, err := gb.mergedPRs(ctx, from, to)
	if err != nil {
		return Section{}, err
	}
	t, err := gb.getChangedAt(to)
	if err != nil {
		return Section{}, err
	}

	return Section{
		PullRequests: r,
		FromRevision: from,
		ToRevision:   to,
		ChangedAt:    t,
		Project:      gb.ProjectKey,
		Repo:         gb.RepoName,

		BaseURL:      gb.BaseURL,
		HTMLURL:      fmt.Sprintf("%s/git/%s/%s", gb.BaseURL, gb.ProjectKey, gb.RepoName),
		ShowUniqueID: gb.ShowUniqueID,
	}, nil
}
