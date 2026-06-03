package shrub

import (
	"errors"
	"fmt"
	"testing"
)

func TestWellformedOperations(t *testing.T) {
	cases := map[string]Command{
		"subprocess.exec":           CmdExec{},
		"shell.exec":                CmdExecShell{},
		"downstream_expansions.set": CmdSetExpansions{},
		"s3Copy.copy":               CmdS3Copy{},
		"s3.get":                    CmdS3Get{},
		"s3.put":                    CmdS3Put{AWSKey: "foo", AWSSecret: "bar", LocalFile: "baz"},
		"s3.push":                   CmdS3Push{},
		"s3.pull":                   CmdS3Pull{},
		"git.get_project":           CmdGetProject{},
		"attach.artifacts":          CmdAttachArtifacts{},
		"attach.results":            CmdResultsJSON{},
		"attach.xunit_results":      CmdResultsXunit{},
		"gotest.parse_files":        CmdResultsGoTest{},
		"archive.zip_pack":          CmdArchiveCreate{Format: ZIP},
		"archive.targz_pack":        CmdArchiveCreate{Format: TARBALL},
		"archive.zip_extract":       CmdArchiveExtract{Format: ZIP},
		"archive.targz_extract":     CmdArchiveExtract{Format: TARBALL},
		"archive.auto_extract":      CmdArchiveExtract{Format: ArchiveFormat("auto")},
		"host.create":               CmdHostCreate{},
		"host.list":                 CmdHostList{},
		"expansions.update":         CmdExpansionsUpdate{},
		"expansions.write":          CmdExpansionsWrite{},
		"json.send":                 CmdJSONSend{},
		"papertrail.trace":          CmdPapertrailTrace{},
		"perf.send":                 CmdPerfSend{},
		"timeout.update":            CmdTimeoutUpdate{},
		"github.generate_token":     CmdGitHubGenerateToken{},
		"test_selection.get":        CmdTestSelectionGet{},
		"cache.save": CmdCacheSave{
			CacheName:     "myCache",
			Bucket:        "my-bucket",
			RemotePath:    "prefix/",
			KeyExpansions: []string{"${project_id}"},
			AWSKey:        "key",
			AWSSecret:     "secret",
			Paths:         []string{"vendor/"},
		},
		"cache.restore": CmdCacheRestore{
			CacheName:     "myCache",
			Bucket:        "my-bucket",
			RemotePath:    "prefix/",
			KeyExpansions: []string{"${project_id}"},
			AWSKey:        "key",
			AWSSecret:     "secret",
		},
	}

	for name, cmd := range cases {
		t.Run("Validate_"+name, func(t *testing.T) {
			assert(t, cmd.Validate() == nil, fmt.Sprintf("validation for %T (%s)", cmd, name))
		})
		t.Run("Resolve_"+name, func(t *testing.T) {
			defer catch(t, name, "resolve")
			rcmd := cmd.Resolve()
			require(t, rcmd != nil, fmt.Sprintf("resolution for %T (%s)", cmd, name))
			assert(t, rcmd.CommandName == name, name, "not equal to", rcmd.FunctionName)
		})

	}
}

type unmarshableCmd struct {
	name string
}

func (u unmarshableCmd) Name() string                { return "foo" }
func (u unmarshableCmd) Validate() error             { return nil }
func (u unmarshableCmd) Resolve() *CommandDefinition { panic("always") }
func (u unmarshableCmd) MarshalJSON() ([]byte, error) {
	return nil, errors.New("always")
}

func TestPoorlyFormedOperations(t *testing.T) {
	cases := map[string]Command{
		"s3put.empty":         CmdS3Put{},
		"s3put.nocreds":       CmdS3Put{LocalFile: "baz"},
		"s3put.nofile":        CmdS3Put{AWSKey: "foo", AWSSecret: "bar"},
		"s3put.nosecret":      CmdS3Put{AWSKey: "foo", LocalFile: "baz"},
		"s3put.nokey":         CmdS3Put{AWSSecret: "bar", LocalFile: "baz"},
		"archive.create_auto":              CmdArchiveCreate{Format: ArchiveFormat("auto")},
		"archive.invalid":                  CmdArchiveExtract{Format: ArchiveFormat("bleh")},
		"cache.save.empty":                 CmdCacheSave{},
		"cache.save.noname":                CmdCacheSave{Bucket: "b", RemotePath: "r", KeyExpansions: []string{"k"}, AWSKey: "k", AWSSecret: "s", Paths: []string{"p"}},
		"cache.save.nobucket":              CmdCacheSave{CacheName: "c", RemotePath: "r", KeyExpansions: []string{"k"}, AWSKey: "k", AWSSecret: "s", Paths: []string{"p"}},
		"cache.save.noremotepath":          CmdCacheSave{CacheName: "c", Bucket: "b", KeyExpansions: []string{"k"}, AWSKey: "k", AWSSecret: "s", Paths: []string{"p"}},
		"cache.save.nokeyexpansions":       CmdCacheSave{CacheName: "c", Bucket: "b", RemotePath: "r", AWSKey: "k", AWSSecret: "s", Paths: []string{"p"}},
		"cache.save.nocreds":               CmdCacheSave{CacheName: "c", Bucket: "b", RemotePath: "r", KeyExpansions: []string{"k"}, Paths: []string{"p"}},
		"cache.save.nopaths":               CmdCacheSave{CacheName: "c", Bucket: "b", RemotePath: "r", KeyExpansions: []string{"k"}, AWSKey: "k", AWSSecret: "s"},
		"cache.save.rolearn_with_key":      CmdCacheSave{CacheName: "c", Bucket: "b", RemotePath: "r", KeyExpansions: []string{"k"}, RoleARN: "arn:aws:iam::123:role/r", AWSKey: "k", Paths: []string{"p"}},
		"cache.restore.empty":              CmdCacheRestore{},
		"cache.restore.noname":             CmdCacheRestore{Bucket: "b", RemotePath: "r", KeyExpansions: []string{"k"}, AWSKey: "k", AWSSecret: "s"},
		"cache.restore.nobucket":           CmdCacheRestore{CacheName: "c", RemotePath: "r", KeyExpansions: []string{"k"}, AWSKey: "k", AWSSecret: "s"},
		"cache.restore.noremotepath":       CmdCacheRestore{CacheName: "c", Bucket: "b", KeyExpansions: []string{"k"}, AWSKey: "k", AWSSecret: "s"},
		"cache.restore.nokeyexpansions":    CmdCacheRestore{CacheName: "c", Bucket: "b", RemotePath: "r", AWSKey: "k", AWSSecret: "s"},
		"cache.restore.nocreds":            CmdCacheRestore{CacheName: "c", Bucket: "b", RemotePath: "r", KeyExpansions: []string{"k"}},
		"cache.restore.rolearn_with_token": CmdCacheRestore{CacheName: "c", Bucket: "b", RemotePath: "r", KeyExpansions: []string{"k"}, RoleARN: "arn:aws:iam::123:role/r", AWSSessionToken: "t"},
	}

	for name, cmd := range cases {
		t.Run("ValidateFailsFor_"+name, func(t *testing.T) {
			assert(t, cmd.Validate() != nil, name)
		})
		t.Run("ResolvePanicsFor_"+name, func(t *testing.T) {
			defer expect(t, name)

			rcmd := cmd.Resolve()
			assert(t, rcmd == nil)
		})
	}

	t.Run("AlwaysPanicsWhenCannotMarshal", func(t *testing.T) {
		defer expect(t, "marshaling")

		res := exportCmd(unmarshableCmd{name: "sad"})
		assert(t, res == nil)
	})
}
