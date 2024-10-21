package cli_test

import (
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/MrNemo64/go-n-i18n/internal/cli"
	"github.com/stretchr/testify/require"
)

type stubWalker struct {
	files   []stubFileEntry
	current int
}
type stubFileEntry struct {
	path     []string
	language string
	contents string
}

func (fe *stubFileEntry) Path() []string                { return fe.path }
func (fe *stubFileEntry) Language() string              { return fe.language }
func (fe *stubFileEntry) FullPath() string              { return strings.Join(fe.path, "/") }
func (fe *stubFileEntry) ReadContents() ([]byte, error) { return []byte(fe.contents), nil }

func StubWalker(files []stubFileEntry) *stubWalker {
	return &stubWalker{
		current: -1,
		files:   files,
	}
}

func (walker *stubWalker) Next() (cli.FileEntry, error) {
	walker.current++
	if walker.current >= len(walker.files) || walker.files == nil {
		return nil, cli.ErrNoMoreFiles
	}
	return &walker.files[walker.current], nil
}

func TestJsonParser(t *testing.T) {
	t.Parallel()

	t.Run("Parser on empty walker", parserOnEmptyWalker)
	t.Run("Parser on one file and one language", parserOnOneFileOneLang)
	t.Run("Parser on several files and one language", parserOnSeveralFilesOneLang)
	t.Run("Parser on several files with colliding bags that merge and one language", parserOnSeveralFilesCollidingOneLang)
	t.Run("Parser on several files with colliding bags that merge and several languages", parserOnSeveralFilesCollidingSeveralLang)

	t.Run("Parsin fails on duplicated key", parseFailOnDuplicate)
	t.Run("Parsin fails on parent key is not bag", parseFailParentKeyNotBag)

	t.Run("Parsin parametrized merging works", parseParametrizedMergingWorks)
}

func parserOnEmptyWalker(t *testing.T) {
	t.Parallel()
	result, err := cli.ParseJson(StubWalker(make([]stubFileEntry, 0)))
	require.NoError(t, err)
	require.Equal(t, cli.MessageEntryMessageBag{}.With("", make([]cli.MessageEntry, 0)), result)
}

func parserOnOneFileOneLang(t *testing.T) {
	t.Parallel()

	result, err := cli.ParseJson(StubWalker([]stubFileEntry{
		{
			path:     []string{},
			language: "en-EN",
			contents: `{
				"key1": "value1",
				"key2": "value2",
				"key3": {
					"key4": "value4",
					"key5": {
						"key6": "value6"
					},
					"args":  {
						"any-param": "param {arg} has no type",
						"typed-param": "param {arg:int} has type int",
						"formatted-param": "param {arg:str:v} has type string and format %v"
					}
				}
			}
			`,
		},
	}))
	require.NoError(t, err)
	expected := bag("", []cli.MessageEntry{
		literal("key1", "en-EN", "value1"),
		literal("key2", "en-EN", "value2"),
		bag("key3", []cli.MessageEntry{
			literal("key4", "en-EN", "value4"),
			bag("key5", []cli.MessageEntry{
				literal("key6", "en-EN", "value6"),
			}),
			bag("args", []cli.MessageEntry{
				param("any-param", "en-EN", "param {arg} has no type"),
				param("typed-param", "en-EN", "param {arg:int} has type int"),
				param("formatted-param", "en-EN", "param {arg:str:v} has type string and format %v"),
			}),
		}),
	})
	require.Equal(t, expected, result)
}

func parserOnSeveralFilesOneLang(t *testing.T) {
	t.Parallel()

	result, err := cli.ParseJson(StubWalker([]stubFileEntry{
		{
			path:     []string{},
			language: "en-EN",
			contents: `{
				"key1": "value1",
				"key2": "value2"
			}
			`,
		},
		{
			path:     []string{"key3"},
			language: "en-EN",
			contents: `{
				"key4": "value4",
				"key5": {
					"key6": "value6"
				}
			}
			`,
		},
	}))
	require.NoError(t, err)
	expected := bag("", []cli.MessageEntry{
		literal("key1", "en-EN", "value1"),
		literal("key2", "en-EN", "value2"),
		bag("key3", []cli.MessageEntry{
			literal("key4", "en-EN", "value4"),
			bag("key5", []cli.MessageEntry{
				literal("key6", "en-EN", "value6"),
			}),
		}),
	})
	require.Equal(t, expected, result)
}

func parserOnSeveralFilesCollidingOneLang(t *testing.T) {
	t.Parallel()

	result, err := cli.ParseJson(StubWalker([]stubFileEntry{
		{
			path:     []string{},
			language: "en-EN",
			contents: `{
				"key1": "value1",
				"key2": "value2",
				"key3": {
					"key4": "value4"
				}
			}
			`,
		},
		{
			path:     []string{"key3"},
			language: "en-EN",
			contents: `{
				"key5": {
					"key6": "value6"
				}
			}
			`,
		},
		{
			path:     []string{"key3", "key5"},
			language: "en-EN",
			contents: `{
				"key7": "value7"
			}
			`,
		},
		{
			path:     []string{"key3", "key5", "key8"},
			language: "en-EN",
			contents: `{
				"key9": "value9",
				"key10": "value10"
			}
			`,
		},
	}))
	require.NoError(t, err)
	expected := bag("", []cli.MessageEntry{
		literal("key1", "en-EN", "value1"),
		literal("key2", "en-EN", "value2"),
		bag("key3", []cli.MessageEntry{
			literal("key4", "en-EN", "value4"),
			bag("key5", []cli.MessageEntry{
				literal("key6", "en-EN", "value6"),
				literal("key7", "en-EN", "value7"),
				bag("key8", []cli.MessageEntry{
					literal("key9", "en-EN", "value9"),
					literal("key10", "en-EN", "value10"),
				}),
			}),
		}),
	})
	require.Equal(t, expected, result)
}

func parserOnSeveralFilesCollidingSeveralLang(t *testing.T) {
	t.Parallel()

	result, err := cli.ParseJson(StubWalker([]stubFileEntry{
		{
			path:     []string{},
			language: "en-EN",
			contents: `{
				"key1": "value1",
				"key2": "value2",
				"key3": {
					"key4": "value4"
				}
			}
			`,
		},
		{
			path:     []string{"key3"},
			language: "en-EN",
			contents: `{
				"key5": {
					"key6": "value6"
				}
			}
			`,
		},
		{
			path:     []string{"key3", "key5"},
			language: "en-EN",
			contents: `{
				"key7": "value7"
			}
			`,
		},
		{
			path:     []string{"key3", "key5", "key8"},
			language: "en-EN",
			contents: `{
				"key9": "value9",
				"key10": "value10"
			}
			`,
		},
		{
			path:     []string{},
			language: "es-ES",
			contents: `{
				"key1": "valor1",
				"key2": "valor2",
				"key3": {
					"key4": "valor4"
				}
			}
			`,
		},
		{
			path:     []string{"key3"},
			language: "es-ES",
			contents: `{
				"key5": {
					"key6": "valor6"
				}
			}
			`,
		},
		{
			path:     []string{"key3", "key5"},
			language: "es-ES",
			contents: `{
				"key7": "valor7"
			}
			`,
		},
		{
			path:     []string{"key3", "key5", "key8"},
			language: "es-ES",
			contents: `{
				"key9": "valor9",
				"key10": "valor10"
			}
			`,
		},
	}))
	require.NoError(t, err)
	expected := bag("", []cli.MessageEntry{
		literal("key1", "en-EN", "value1", "es-ES", "valor1"),
		literal("key2", "en-EN", "value2", "es-ES", "valor2"),
		bag("key3", []cli.MessageEntry{
			literal("key4", "en-EN", "value4", "es-ES", "valor4"),
			bag("key5", []cli.MessageEntry{
				literal("key6", "en-EN", "value6", "es-ES", "valor6"),
				literal("key7", "en-EN", "value7", "es-ES", "valor7"),
				bag("key8", []cli.MessageEntry{
					literal("key9", "en-EN", "value9", "es-ES", "valor9"),
					literal("key10", "en-EN", "value10", "es-ES", "valor10"),
				}),
			}),
		}),
	})
	require.Equal(t, expected, result)
}

func parseFailOnDuplicate(t *testing.T) {
	t.Parallel()
	_, err := cli.ParseJson(StubWalker([]stubFileEntry{
		{
			path:     []string{},
			language: "en-EN",
			contents: `{
				"key1": "value1",
				"key3": {
					"key4": "value4",
					"key5": {
						"key6": "value6"
					}
				}
			}
			`,
		},
		{
			path:     []string{},
			language: "en-EN",
			contents: `{
				"key1": "again"
			}
			`,
		},
	}))
	require.ErrorIs(t, err, cli.ErrLiteralMessageRedefinition)
}

func parseFailParentKeyNotBag(t *testing.T) {
	t.Parallel()
	_, err := cli.ParseJson(StubWalker([]stubFileEntry{
		{
			path:     []string{},
			language: "en-EN",
			contents: `{
				"key1": "value1",
				"key2": "value2",
				"key3": {
					"key4": "value4",
					"key5": {
						"key6": "value6"
					}
				}
			}
			`,
		},
		{
			path:     []string{},
			language: "en-EN",
			contents: `{
				"key2": {
					"key7": "since key 2 has been defined in the file before as 'value2' this will fail because it requires key2 to be a bag"
				}
			}
			`,
		},
	}))
	require.Error(t, err)
}

func parseParametrizedMergingWorks(t *testing.T) {
	t.Parallel()

	result, err := cli.ParseJson(StubWalker([]stubFileEntry{
		{
			path:     []string{},
			language: "en-EN",
			contents: `{
				"key1": "{arg}",
				"key2": "{arg2:int}",
				"key3": "{arg3:int:v}"
			}
			`,
		},
		{
			path:     []string{},
			language: "es-ES",
			contents: `{
				"key1": "{arg:string}",
				"key2": "{newArg}",
				"key3": "{arg3}"
			}
			`,
		},
	}))
	require.NoError(t, err)
	expected := bag("", []cli.MessageEntry{
		param("key1", "en-EN", "{arg}", "es-ES", "{arg:string}"),
		param("key2", "en-EN", "{arg2:int}", "es-ES", "{newArg}"),
		param("key3", "en-EN", "{arg3:int:v}", "es-ES", "{arg3}"),
	})
	require.Equal(t, expected, result)
	strType := cli.FindArgumentType("string")
	intType := cli.FindArgumentType("int")
	anyType := cli.AnyKind()
	entry, err := result.AsBag().GetEntry("key1")
	require.NoError(t, err)
	require.Equal(t, []*cli.MessageArgument{{Name: "arg", Type: strType, Format: strType.DefaultFormat}}, entry.AsParametrized().Args())
	entry, err = result.AsBag().GetEntry("key2")
	require.NoError(t, err)
	require.Equal(t, []*cli.MessageArgument{
		{Name: "arg2", Type: intType, Format: intType.DefaultFormat},
		{Name: "newArg", Type: anyType, Format: anyType.DefaultFormat},
	}, entry.AsParametrized().Args())
	entry, err = result.AsBag().GetEntry("key3")
	require.NoError(t, err)
	require.Equal(t, []*cli.MessageArgument{{Name: "arg3", Type: intType, Format: "v"}}, entry.AsParametrized().Args())
}

func bag(key string, entries []cli.MessageEntry) *cli.MessageEntryMessageBag {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key() < entries[j].Key()
	})
	e := cli.MessageEntryMessageBag{}.With(key, entries)
	for _, entry := range entries {
		entry.AssignParent(e)
	}
	return e
}

func literal(key string, message ...string) *cli.MessageEntryLiteralString {
	if len(message)%2 != 0 {
		panic("")
	}
	m := map[string]string{}
	for i := 1; i < len(message); i += 2 {
		m[message[i-1]] = message[i]
	}
	return cli.MessageEntryLiteralString{}.With(key, m)
}

func param(key string, message ...string) *cli.MessageEntryParametrizedString {
	re := regexp.MustCompile(`\{([a-zA-Z_][a-zA-Z0-9_]*)(?::([a-zA-Z0-9_]+))?(?::([a-zA-Z0-9_%.]+))?\}`)
	if len(message)%2 != 0 {
		panic("")
	}
	m := map[string]string{}
	for i := 1; i < len(message); i += 2 {
		m[message[i-1]] = message[i]
	}
	p := cli.MessageEntryParametrizedString{}.With(key, m)
	for _, msg := range m {
		args := re.FindAllStringSubmatch(msg, -1)
		for _, arg := range args {
			if err := p.AddArgument(arg[1], arg[2], arg[3]); err != nil {
				panic(err)
			}
		}
	}
	return p
}
