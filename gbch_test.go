package gbch

import (
	"reflect"
	"testing"
)

func TestParsePRLogs(t *testing.T) {

	input := `6191693 Merge pull request #225 mackerelio/fix-test-for-invaild-toml into master
dbb1d50 Merge pull request #224 mackerelio/retry-retire into master
922eb55 Merge pull request #223 mackerelio/remove_vet into master
aa18e36 Merge pull request #222 mackerelio/fix-comments into master
71da053 Merge pull request #221 yukiyan/fix-typo into master
0925081 Merge pull request #217 mackerelio/remove-usr-local-bin-again into master
b4bc51b Merge pull request #216 mackerelio/bump-version-0.30.2 into master
a8ea16b Merge pull request #215 mackerelio/revert-9e0c8ab1 into master
98b28f1 Merge pull request #214 mackerelio/bump-version-0.30.1 into master
19a0010 Merge pull request #213 mackerelio/workaround-amd64 into master
9e0c8ab Merge pull request #211 mackerelio/usr-bin into master
7d278aa Merge pull request #210 mackerelio/bump-version-0.30.0 into master
ce37096 Merge pull request #208 mackerelio/refactor-net-interface into master
8a07070 Merge pull request #207 mackerelio/subcommand-init into master
a39ca5e Merge pull request #209 mackerelio/remove-cpu-flags into master
c0ed1f1 Merge pull request #205 mackerelio/interface-ips into master
8cc281c Merge pull request #202 mackerelio/remove-deprecated-sensu into master
afeb5e5 Merge pull request #161 mackerelio/remove-uptime into master
a75f8b2 Merge pull request #206 mackerelio/bump-version-0.29.2 into master
fd40654 Merge pull request #174 mackerelio/travis-docker into master
a8665f5 Merge branch 'master' into travis-docker into master
bdb2271 Merge pull request #203 mackerelio/alternative-build into master
2ac5301 Merge pull request #199 mackerelio/fix-deb into master
7c79f92 Merge pull request #201 mackerelio/bump-version-0.29.1 into master
32a3e1f Merge pull request #200 mackerelio/bump-version-0.29.0 into master
b4b8c2c Merge pull request #197 hanazuki/check-timeouts into master
a30e851 Merge pull request #198 mackerelio/dont-ignore-logging-level_string into master
2ec717e Merge pull request #196 mackerelio/refactor-around-start into master
ca345ea Merge pull request #195 mackerelio/introduce-motemen-go-cli into master
843b32e Merge pull request #194 mackerelio/remove-deprecated into master
87375ec Merge pull request #193 mackerelio/bump-version-0.28.1 into master
82ccaa3 Merge branch 'master' of github.com:mackerelio/mackerel-agent
4a6d83c Merge pull request #192 mackerelio/deb_init_d_stop_retval into master
5b0a536 Merge pull request #191 mackerelio/gofmt-on-travis into master
`
	expect := []*mergedPRLog{
		{num: 225, branch: "mackerelio/fix-test-for-invaild-toml"},
		{num: 224, branch: "mackerelio/retry-retire"},
		{num: 223, branch: "mackerelio/remove_vet"},
		{num: 222, branch: "mackerelio/fix-comments"},
		{num: 221, branch: "yukiyan/fix-typo"},
		{num: 217, branch: "mackerelio/remove-usr-local-bin-again"},
		{num: 216, branch: "mackerelio/bump-version-0.30.2"},
		{num: 215, branch: "mackerelio/revert-9e0c8ab1"},
		{num: 214, branch: "mackerelio/bump-version-0.30.1"},
		{num: 213, branch: "mackerelio/workaround-amd64"},
		{num: 211, branch: "mackerelio/usr-bin"},
		{num: 210, branch: "mackerelio/bump-version-0.30.0"},
		{num: 208, branch: "mackerelio/refactor-net-interface"},
		{num: 207, branch: "mackerelio/subcommand-init"},
		{num: 209, branch: "mackerelio/remove-cpu-flags"},
		{num: 205, branch: "mackerelio/interface-ips"},
		{num: 202, branch: "mackerelio/remove-deprecated-sensu"},
		{num: 161, branch: "mackerelio/remove-uptime"},
		{num: 206, branch: "mackerelio/bump-version-0.29.2"},
		{num: 174, branch: "mackerelio/travis-docker"},
		{num: 203, branch: "mackerelio/alternative-build"},
		{num: 199, branch: "mackerelio/fix-deb"},
		{num: 201, branch: "mackerelio/bump-version-0.29.1"},
		{num: 200, branch: "mackerelio/bump-version-0.29.0"},
		{num: 197, branch: "hanazuki/check-timeouts"},
		{num: 198, branch: "mackerelio/dont-ignore-logging-level_string"},
		{num: 196, branch: "mackerelio/refactor-around-start"},
		{num: 195, branch: "mackerelio/introduce-motemen-go-cli"},
		{num: 194, branch: "mackerelio/remove-deprecated"},
		{num: 193, branch: "mackerelio/bump-version-0.28.1"},
		{num: 192, branch: "mackerelio/deb_init_d_stop_retval"},
		{num: 191, branch: "mackerelio/gofmt-on-travis"},
	}
	if !reflect.DeepEqual(parseMergedPRLogs(input), expect) {
		t.Errorf("somthing went wrong")
	}
}
