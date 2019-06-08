package gbch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"net/url"

	"github.com/Songmu/gitsemvers"
	"github.com/pkg/errors"
	"github.com/vvatanabe/errsgroup"
	backlog "github.com/vvatanabe/go-backlog/backlog/v2"
)

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// Run the ghch
func (gb *Gbch) Run() error {
	ctx := context.Background()
	if err := gb.initialize(ctx); err != nil {
		return err
	}
	if gb.All {
		return gb.runAll(ctx)
	}
	return gb.run(ctx)
}

func (gb *Gbch) runAll(ctx context.Context) error {
	chlog := Changelog{}
	vers := append(gb.versions(), "")
	prevRev := ""
	for _, rev := range vers {
		r, err := gb.getSection(ctx, rev, prevRev)
		if err != nil {
			return err
		}
		if prevRev == "" && gb.NextVersion != "" {
			r.ToRevision = gb.NextVersion
		}
		chlog.Sections = append(chlog.Sections, r)
		prevRev = rev
	}

	if gb.Format != "markdown" { // json
		encoder := json.NewEncoder(gb.OutStream)
		encoder.SetIndent("", "  ")
		return encoder.Encode(chlog)
	}
	results := make([]string, len(chlog.Sections))
	for i, v := range chlog.Sections {
		results[i], _ = v.toMkdn()
	}
	if gb.Write {
		content := "# Changelog\n\n" + strings.Join(results, "\n\n")
		if err := ioutil.WriteFile(gb.ChangelogMd, []byte(content), 0644); err != nil {
			return err
		}
	} else {
		fmt.Fprintln(gb.OutStream, strings.Join(results, "\n\n"))
	}
	return nil
}

func (gb *Gbch) run(ctx context.Context) error {
	if gb.Latest {
		vers := gb.versions()
		if len(vers) > 0 {
			gb.To = vers[0]
		}
		if gb.From == "" && len(vers) > 1 {
			gb.From = vers[1]
		}
	} else if gb.From == "" && gb.To == "" {
		gb.From = gb.getLatestSemverTag()
	}
	r, err := gb.getSection(ctx, gb.From, gb.To)
	if err != nil {
		return err
	}
	if r.ToRevision == "" && gb.NextVersion != "" {
		r.ToRevision = gb.NextVersion
	}

	if gb.Format != "markdown" { // json
		encoder := json.NewEncoder(gb.OutStream)
		encoder.SetIndent("", "  ")
		return encoder.Encode(r)
	}
	str, err := r.toMkdn()
	if err != nil {
		return err
	}
	if gb.Write {
		content := ""
		if exists(gb.ChangelogMd) {
			byt, err := ioutil.ReadFile(gb.ChangelogMd)
			if err != nil {
				return err
			}
			content = insertNewChangelog(byt, str)
		} else {
			content = "# Changelog\n\n" + str + "\n"
		}
		if err := ioutil.WriteFile(gb.ChangelogMd, []byte(content), 0644); err != nil {
			return err
		}
	} else {
		fmt.Fprintln(gb.OutStream, str)
	}
	return nil
}

func (gb *Gbch) initialize(ctx context.Context) error {
	if gb.Write {
		gb.Format = "markdown"
		if gb.ChangelogMd == "" {
			gb.ChangelogMd = "CHANGELOG.md"
		}
	}
	if gb.OutStream == nil {
		gb.OutStream = os.Stdout
	}

	remoteURL, err := gb.getRemoteURL()
	if err != nil {
		return err
	}

	spaceDomain := gb.spaceDomain(remoteURL)
	gb.BaseURL = fmt.Sprintf("https://%s", spaceDomain)
	gb.ProjectKey, gb.RepoName = gb.projectKeyAndRepo(remoteURL)

	gb.client = backlog.NewClient(spaceDomain, nil)
	gb.setAPIKey()
	if gb.APIKey == "" {
		return errors.New("backlog api key is empty")
	}
	gb.client.SetAPIKey(gb.APIKey)
	return nil
}

func (gb *Gbch) setAPIKey() {
	if gb.APIKey != "" {
		return
	}
	if gb.APIKey = os.Getenv("BACKLOG_API_KEY"); gb.APIKey != "" {
		return
	}
	return
}

func (gb *Gbch) gitProg() string {
	if gb.GitPath != "" {
		return gb.GitPath
	}
	return "git"
}

func (gb *Gbch) cmd(argv ...string) (string, error) {
	arg := []string{"-C", gb.RepoPath}
	arg = append(arg, argv...)
	cmd := exec.Command(gb.gitProg(), arg...)
	cmd.Env = append(os.Environ(), "LANG=C")

	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return b.String(), err
}

func (gb *Gbch) versions() []string {
	sv := gitsemvers.Semvers{
		RepoPath: gb.RepoPath,
		GitPath:  gb.GitPath,
	}
	return sv.VersionStrings()
}

func (gb *Gbch) getRemote() string {
	if gb.Remote != "" {
		return gb.Remote
	}
	return "origin"
}

type RemoteURL struct {
	Protocol string
	Host     string
	Port     string
	Path     string
}

func (gb *Gbch) getRemoteURL() (*RemoteURL, error) {
	out, _ := gb.cmd("remote", "-v")
	remotes := strings.Split(out, "\n")

	remote := gb.getRemote()
	var remoteURL string
	for _, r := range remotes {
		fields := strings.Fields(r)
		if len(fields) > 1 && fields[0] == remote {
			remoteURL = fields[1]
			break
		}
	}

	var ep *RemoteURL
	if isHTTP(remoteURL) {
		ep = toRemoteURLFromHTTP(remoteURL)
	} else if isSSH(remoteURL) {
		ep = toRemoteURLFromSSH(remoteURL)
	} else {
		return ep, errors.New("could not be used protocol except http and ssh")
	}

	return ep, nil
}

var repoURLReg = regexp.MustCompile(`([^/:]+)/([^/]+?)(?:\.git)?$`)

func (gb *Gbch) projectKeyAndRepo(remoteURL *RemoteURL) (projectKey, repo string) {
	if matches := repoURLReg.FindStringSubmatch(remoteURL.Path); len(matches) > 2 {
		return matches[1], matches[2]
	}
	return
}

var serviceDomains = []string{"backlog.jp", "backlog.com", "backlogtool.com"}

func (gb *Gbch) spaceDomain(remoteURL *RemoteURL) string {

	var isBacklogDomain bool
	for _, d := range serviceDomains {
		if strings.HasSuffix(remoteURL.Host, "."+d) {
			isBacklogDomain = true
			break
		}

	}

	if !isBacklogDomain {
		return remoteURL.Host
	}

	if strings.HasPrefix(remoteURL.Protocol, "http") {
		return remoteURL.Host
	}

	// ignore `git` from ssh host (foo.git.backlog.jp)
	delimitedHost := strings.Split(remoteURL.Host, ".")
	spaceKey := delimitedHost[0]
	domain := strings.Join(delimitedHost[len(delimitedHost)-2:], ".")
	return fmt.Sprintf("%s.%s", spaceKey, domain)
}

var errsGroupLimitSize = errsgroup.LimitSize(4)

func (gb *Gbch) mergedPRs(ctx context.Context, from, to string) (prs []*backlog.PullRequest, err error) {
	prlogs, err := gb.mergedPRLogs(from, to)
	if err != nil {
		return
	}
	prs = make([]*backlog.PullRequest, 0, len(prlogs))
	g := errsgroup.NewGroup(errsGroupLimitSize)
	for _, v := range prlogs {
		prlog := v
		g.Go(func() (err error) {
			pr, resp, err := gb.client.PullRequests.GetPullRequest(ctx, gb.ProjectKey, gb.RepoName, prlog.num)
			if err != nil {
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					return
				}
				log.Println(err)
				return
			}
			if pr.Branch != prlog.branch {
				return
			}
			prs = append(prs, pr)
			return
		})
	}
	for _, e := range g.Wait() {
		err = e
	}
	return
}

func (gb *Gbch) getLatestSemverTag() string {
	vers := gb.versions()
	if len(vers) < 1 {
		return ""
	}
	return vers[0]
}

type mergedPRLog struct {
	num    int
	branch string
}

func (gb *Gbch) mergedPRLogs(from, to string) (nums []*mergedPRLog, err error) {
	revisionRange := fmt.Sprintf("%s..%s", from, to)
	out, err := gb.cmd("log", revisionRange, "--merges", "--oneline")
	if err != nil {
		return []*mergedPRLog{}, err
	}
	return parseMergedPRLogs(out), nil
}

var prMergeReg = regexp.MustCompile(`^[a-f0-9]+ (?:\(.+\) )?Merge pull request #([0-9]+) (\S+) into \S+`)

func parseMergedPRLogs(out string) (prs []*mergedPRLog) {
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if matches := prMergeReg.FindStringSubmatch(line); len(matches) > 2 {
			i, _ := strconv.Atoi(matches[1])
			prs = append(prs, &mergedPRLog{
				num:    i,
				branch: matches[2],
			})
		}
	}
	return
}

func (gb *Gbch) getChangedAt(rev string) (time.Time, error) {
	if rev == "" {
		rev = "HEAD"
	}
	out, err := gb.cmd("show", "-s", rev+"^{commit}", `--format=%ct`)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "failed to get changed at from git revision. `git show` failed")
	}
	out = strings.TrimSpace(out)
	i, err := strconv.ParseInt(out, 10, 64)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "failed to get changed at from git revision. ParseInt failed")
	}
	return time.Unix(i, 0), nil
}

func isHTTP(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme == "https" && u.Host != ""
}

func toRemoteURLFromHTTP(str string) *RemoteURL {
	u, _ := url.Parse(str)
	return &RemoteURL{
		Protocol: u.Scheme,
		Host:     u.Host,
		Port:     u.Port(),
		Path:     u.Path,
	}
}

var sshUrlReg = regexp.MustCompile(`^(?:(?P<user>[^@]+)@)?(?P<host>[^:\s]+):(?:(?P<port>[0-9]{1,5})/)?(?P<path>[^\\].*)$`)

func isSSH(str string) bool {
	return sshUrlReg.MatchString(str)
}

func toRemoteURLFromSSH(str string) *RemoteURL {
	m := sshUrlReg.FindStringSubmatch(str)
	return &RemoteURL{
		Protocol: "ssh",
		Host:     m[2],
		Port:     m[3],
		Path:     m[4],
	}
}
