package editor_test

import (
	"context"
	"strings"
	"testing"

	"github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/app/editor"
	"github.com/fakovacic/editor/internal/app/editor/mocks"
	"github.com/stretchr/testify/assert"
)

func TestChangesComplex(t *testing.T) {
	const cssExample = `pre code.hljs {
	display: block;
	overflow-x: auto;
	padding: 1em;
}

code.hljs {
	padding: 3px 5px;
}
`

	cases := []struct {
		it string

		req []*app.ChangeMsg

		expectedResponse string
		expectedError    string
	}{
		{
			it: "change padding property",
			req: []*app.ChangeMsg{
				{
					Action: "remove",
					Start: app.ChangeRow{
						Row:    3,
						Column: 13,
					},
					End: app.ChangeRow{
						Row:    3,
						Column: 16,
					},
					Lines: []string{"1em"},
				},
				{
					Action: "insert",
					Start: app.ChangeRow{
						Row:    3,
						Column: 13,
					},
					End: app.ChangeRow{
						Row:    3,
						Column: 14,
					},
					Lines: []string{"2"},
				},
				{
					Action: "insert",
					Start: app.ChangeRow{
						Row:    3,
						Column: 14,
					},
					End: app.ChangeRow{
						Row:    3,
						Column: 15,
					},
					Lines: []string{"e"},
				},
				{
					Action: "insert",
					Start: app.ChangeRow{
						Row:    3,
						Column: 15,
					},
					End: app.ChangeRow{
						Row:    3,
						Column: 16,
					},
					Lines: []string{"m"},
				},
			},

			expectedResponse: `pre code.hljs {
	display: block;
	overflow-x: auto;
	padding: 2em;
}

code.hljs {
	padding: 3px 5px;
}
`,
		},
		{
			it: "remove code.hljs property",
			req: []*app.ChangeMsg{
				{
					Action: "remove",
					Start: app.ChangeRow{
						Row:    5,
						Column: 0,
					},
					End: app.ChangeRow{
						Row:    9,
						Column: 0,
					},
					Lines: []string{"", "code.hljs {", "    padding: 3px 5px;", "}", ""},
				},
			},

			expectedResponse: `pre code.hljs {
	display: block;
	overflow-x: auto;
	padding: 1em;
}
`,
		},
		{
			it: "remove new lines",
			req: []*app.ChangeMsg{
				{
					Action: "remove",
					Start: app.ChangeRow{
						Row:    4,
						Column: 1,
					},
					End: app.ChangeRow{
						Row:    6,
						Column: 0,
					},
					Lines: []string{"", "", ""},
				},
			},

			expectedResponse: `pre code.hljs {
	display: block;
	overflow-x: auto;
	padding: 1em;
}code.hljs {
	padding: 3px 5px;
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.it, func(t *testing.T) {
			io := &mocks.IOMock{
				ReadFunc: func(ctx context.Context) (string, *app.FileMeta, error) {
					return cssExample, nil, nil
				},
			}

			editor := editor.New(io)

			err := editor.Load(context.Background())
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedError)
			}

			for _, req := range tc.req {
				err = editor.Change(context.Background(), req)
				if err != nil {
					assert.Equal(t, err.Error(), tc.expectedError)
				}
			}

			res, err := editor.Read(context.Background())
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedError)
			}

			assert.Equal(t, strings.ReplaceAll(tc.expectedResponse, "\t", "    "), res)
		})
	}
}

func TestChangesNewLine(t *testing.T) {
	const cssExample = `pre code.hljs {
	display: block;
	overflow-x: auto;
	padding: 1em;
}

code.hljs {
	padding: 3px 5px;
}`

	cases := []struct {
		it string

		req []*app.ChangeMsg

		expectedResponse string
		expectedError    string
	}{
		{
			it: "insert new line #1",
			req: []*app.ChangeMsg{
				{
					Action: "insert",
					Start: app.ChangeRow{
						Row:    5,
						Column: 0,
					},
					End: app.ChangeRow{
						Row:    6,
						Column: 0,
					},
					Lines: []string{"", ""},
				},
			},

			expectedResponse: `pre code.hljs {
	display: block;
	overflow-x: auto;
	padding: 1em;
}


code.hljs {
	padding: 3px 5px;
}
`,
		},
		{
			it: "insert new line #2",
			req: []*app.ChangeMsg{
				{
					Action: "insert",
					Start: app.ChangeRow{
						Row:    2,
						Column: 21,
					},
					End: app.ChangeRow{
						Row:    3,
						Column: 0,
					},
					Lines: []string{"", ""},
				},
			},

			expectedResponse: `pre code.hljs {
	display: block;
	overflow-x: auto;

	padding: 1em;
}

code.hljs {
	padding: 3px 5px;
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.it, func(t *testing.T) {
			io := &mocks.IOMock{
				ReadFunc: func(ctx context.Context) (string, *app.FileMeta, error) {
					return cssExample, nil, nil
				},
			}

			editor := editor.New(io)

			err := editor.Load(context.Background())
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedError)
			}

			for _, req := range tc.req {
				err = editor.Change(context.Background(), req)
				if err != nil {
					assert.Equal(t, err.Error(), tc.expectedError)
				}
			}

			res, err := editor.Read(context.Background())
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedError)
			}

			assert.Equal(t, strings.ReplaceAll(tc.expectedResponse, "\t", "    "), res)
		})
	}
}

func TestChangesRemoveLine(t *testing.T) {
	const cssExample = `pre code.hljs {
	display: block;
	overflow-x: auto;
	margin:1em;
}



code.hljs {
		padding: 3px 5px;

}
`

	cases := []struct {
		it string

		req []*app.ChangeMsg

		expectedResponse string
		expectedError    string
	}{
		{
			it: "remove new line #1",
			req: []*app.ChangeMsg{
				{
					Action: "remove",
					Start: app.ChangeRow{
						Row:    5,
						Column: 0,
					},
					End: app.ChangeRow{
						Row:    7,
						Column: 0,
					},
					Lines: []string{"", "", ""},
				},
			},

			expectedResponse: `pre code.hljs {
	display: block;
	overflow-x: auto;
	margin:1em;
}

code.hljs {
		padding: 3px 5px;

}
`,
		},
		{
			it: "remove new line #2",
			req: []*app.ChangeMsg{
				{
					Action: "remove",
					Start: app.ChangeRow{
						Row:    9,
						Column: 25,
					},
					End: app.ChangeRow{
						Row:    10,
						Column: 0,
					},
					Lines: []string{"", ""},
				},
			},

			expectedResponse: `pre code.hljs {
	display: block;
	overflow-x: auto;
	margin:1em;
}



code.hljs {
		padding: 3px 5px;
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.it, func(t *testing.T) {
			io := &mocks.IOMock{
				ReadFunc: func(ctx context.Context) (string, *app.FileMeta, error) {
					return cssExample, nil, nil
				},
			}

			editor := editor.New(io)

			err := editor.Load(context.Background())
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedError)
			}

			for _, req := range tc.req {
				err = editor.Change(context.Background(), req)
				if err != nil {
					assert.Equal(t, err.Error(), tc.expectedError)
				}
			}

			res, err := editor.Read(context.Background())
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedError)
			}

			assert.Equal(t, strings.ReplaceAll(tc.expectedResponse, "\t", "    "), res)
		})
	}
}
