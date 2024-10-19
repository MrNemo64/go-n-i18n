package cli_test

import (
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
					}
				}
			}
			`,
		},
	}))
	require.NoError(t, err)
	require.Equal(t, bag("", []cli.MessageEntry{
		literal("key1", "en-EN", "value1"),
		literal("key2", "en-EN", "value2"),
		bag("key3", []cli.MessageEntry{
			literal("key4", "en-EN", "value4"),
			bag("key5", []cli.MessageEntry{
				literal("key6", "en-EN", "value6"),
			}),
		}),
	}), result)
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
	require.Equal(t, bag("", []cli.MessageEntry{
		literal("key1", "en-EN", "value1"),
		literal("key2", "en-EN", "value2"),
		bag("key3", []cli.MessageEntry{
			literal("key4", "en-EN", "value4"),
			bag("key5", []cli.MessageEntry{
				literal("key6", "en-EN", "value6"),
			}),
		}),
	}), result)
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

func bag(key string, entries []cli.MessageEntry) *cli.MessageEntryMessageBag {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key() < entries[j].Key()
	})
	return cli.MessageEntryMessageBag{}.With(key, entries)
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
