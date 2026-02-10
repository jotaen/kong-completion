package kongcompletion

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/posener/complete"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	envLine  = "COMP_LINE"
	envPoint = "COMP_POINT"
)

func TestComplete(t *testing.T) {
	type embed struct {
		Lion string
	}

	predictors := map[string]complete.Predictor{
		"things":      complete.PredictSet("thing1", "thing2"),
		"otherthings": complete.PredictSet("otherthing1", "otherthing2"),
	}

	var cli struct {
		Foo struct {
			Embedded embed  `kong:"embed"`
			Bar      string `kong:"completion-predictor=things"`
			Baz      bool
			Tata     string   `kong:"aliases=titi"`        // one alias
			Xuxu     string   `kong:"aliases='xoxo,xixi'"` // multiple aliases
			Qux      bool     `kong:"hidden"`              // regular hidden
			Quy      bool     `kong:"completion-enabled=false"`
			Quz      bool     `kong:"hidden,completion-enabled=true"`
			Rabbit   struct{} `kong:"cmd"`
			Eagle    struct{} `kong:"cmd,completion-enabled=false"`
			Duck     struct{} `kong:"cmd,aliases=bird"`
		} `kong:"cmd"`
		Bar struct {
			Tiger    string `kong:"arg,completion-predictor=things"`
			Bear     string `kong:"arg,completion-predictor=otherthings"`
			Elephant string `kong:"arg,completion-predictor=${a}${b},set=b=things"`
			OMG      string `kong:"required,enum='oh,my,gizzles'"`
			Number   int    `kong:"required,short=n,enum='1,2,3'"`
			BooFlag  bool   `kong:"name=boofl,short=b"`
		} `kong:"cmd,set=a=other"`
		Baz struct{} `kong:"cmd,hidden"`

		Global string `kong:""`
	}

	for _, td := range []completeTest{
		{
			parser: kong.Must(&cli),
			want:   []string{"foo", "bar"},
			line:   "myApp ",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{"foo"},
			line:   "myApp foo",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{"rabbit", "duck", "bird"},
			line:   "myApp foo ",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{"rabbit"},
			line:   "myApp foo r",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{"--bar", "--baz", "--tata", "--titi", "--xuxu", "--xoxo", "--xixi", "--quz", "--lion", "--help", "-h", "--global"},
			line:   "myApp foo -",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{"--bar", "--baz", "--tata", "--titi", "--xuxu", "--xoxo", "--xixi", "--quz", "--lion", "--help", "--global"},
			line:   "myApp foo --",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{},
			line:   "myApp foo --lion ",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{"rabbit", "duck", "bird"},
			line:   "myApp foo --baz ",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{"--bar", "--baz", "--tata", "--titi", "--xuxu", "--xoxo", "--xixi", "--quz", "--lion", "--help", "-h", "--global"},
			line:   "myApp foo --baz -",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{"thing1", "thing2"},
			line:   "myApp foo --bar ",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{"thing1", "thing2"},
			line:   "myApp bar ",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{"thing1", "thing2"},
			line:   "myApp bar thing",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{"otherthing1", "otherthing2"},
			line:   "myApp bar thing1 ",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{"oh", "my", "gizzles"},
			line:   "myApp bar --omg ",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{"-n", "--number", "--omg", "--help", "-h", "--boofl", "-b", "--global"},
			line:   "myApp bar -",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{"thing1", "thing2"},
			line:   "myApp bar -b ",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{"-n", "--number", "--omg", "--help", "-h", "--boofl", "-b", "--global"},
			line:   "myApp bar -b thing1 -",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{"oh", "my", "gizzles"},
			line:   "myApp bar -b thing1 --omg ",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{"otherthing1", "otherthing2"},
			line:   "myApp bar -b thing1 --omg gizzles ",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{"foo", "bar"},
			line:   "myApp --global=test ",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{"rabbit", "duck", "bird"},
			line:   "myApp foo --global=test ",
		},
		{
			parser: kong.Must(&cli),
			want:   []string{"thing1", "thing2"},
			line:   "myApp bar --global=test ",
		},
	} {
		name := td.name
		if name == "" {
			name = td.line
		}
		t.Run(name, func(t *testing.T) {
			options := []Option{
				WithPredictors(predictors),
			}
			got := runComplete(t, td.parser, td.line, options)
			assert.ElementsMatch(t, td.want, got)
		})
	}
}

func Test_tagPredictor(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		got, err := tagPredictor(nil, nil, nil)
		assert.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("no predictor tag", func(t *testing.T) {
		got, err := tagPredictor(testTag{}, nil, nil)
		assert.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("missing predictor", func(t *testing.T) {
		got, err := tagPredictor(testTag{predictorTag: "foo"}, nil, nil)
		assert.Error(t, err)
		assert.ErrorContains(t, err, `no predictor with name "foo"`)
		assert.Nil(t, got)
	})

	t.Run("existing predictor", func(t *testing.T) {
		got, err := tagPredictor(testTag{predictorTag: "foo"}, map[string]complete.Predictor{"foo": complete.PredictAnything}, nil)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})

	t.Run("interpolation", func(t *testing.T) {
		vars := kong.Vars{"VAR": "foo"}
		got, err := tagPredictor(testTag{predictorTag: "${VAR}"}, map[string]complete.Predictor{"foo": complete.PredictAnything}, vars)
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}

func TestCompleteCumulativeFlags(t *testing.T) {
	t.Run("single cumulative arg", func(t *testing.T) {
		predictors := map[string]complete.Predictor{
			"sources": complete.PredictSet("src1", "src2", "src3"),
		}

		var cli struct {
			Sources []string `kong:"arg,completion-predictor=sources"`
		}

		for _, td := range []completeTest{
			{line: "myApp ", want: []string{"src1", "src2", "src3"}},
			{line: "myApp src", want: []string{"src1", "src2", "src3"}},
			{line: "myApp src1 ", want: []string{"src1", "src2", "src3"}},
			{line: "myApp src1 src2 ", want: []string{"src1", "src2", "src3"}},
			{line: "myApp src1 src2 src3 ", want: []string{"src1", "src2", "src3"}},
		} {
			t.Run(td.line, func(t *testing.T) {
				got := runComplete(t, kong.Must(&cli), td.line, []Option{WithPredictors(predictors)})
				assert.ElementsMatch(t, td.want, got)
			})
		}
	})

	t.Run("cumulative and non-cumulative args mixed", func(t *testing.T) {
		predictors := map[string]complete.Predictor{
			"dests":   complete.PredictSet("dest1", "dest2"),
			"sources": complete.PredictSet("src1", "src2", "src3"),
		}

		var cli struct {
			Dest    string   `kong:"arg,completion-predictor=dests"`
			Sources []string `kong:"arg,completion-predictor=sources"`
		}

		for _, td := range []completeTest{
			{line: "myApp ", want: []string{"dest1", "dest2"}},
			{line: "myApp dest1 ", want: []string{"src1", "src2", "src3"}},
			{line: "myApp dest1 src1 ", want: []string{"src1", "src2", "src3"}},
			{line: "myApp dest1 src1 src2 ", want: []string{"src1", "src2", "src3"}},
		} {
			t.Run(td.line, func(t *testing.T) {
				got := runComplete(t, kong.Must(&cli), td.line, []Option{WithPredictors(predictors)})
				assert.ElementsMatch(t, td.want, got)
			})
		}
	})

	t.Run("cumulative and non-cumulative flag mixed", func(t *testing.T) {
		predictors := map[string]complete.Predictor{
			"sources": complete.PredictSet("src1", "src2", "src3"),
		}

		var cli struct {
			Verbose bool     `kong:""`
			Sources []string `kong:"arg,completion-predictor=sources"`
		}

		for _, td := range []completeTest{
			{line: "myApp ", want: []string{"src1", "src2", "src3"}},
			{line: "myApp --verbose ", want: []string{"src1", "src2", "src3"}},
			{line: "myApp src1 ", want: []string{"src1", "src2", "src3"}},
			{line: "myApp --verbose src1 ", want: []string{"src1", "src2", "src3"}},
			{line: "myApp src1 --verbose src2 ", want: []string{"src1", "src2", "src3"}},
		} {
			t.Run(td.line, func(t *testing.T) {
				got := runComplete(t, kong.Must(&cli), td.line, []Option{WithPredictors(predictors)})
				assert.ElementsMatch(t, td.want, got)
			})
		}
	})
}

type testTag map[string]string

func (t testTag) Has(k string) bool {
	_, ok := t[k]
	return ok
}

func (t testTag) Get(k string) string {
	return t[k]
}

type completeTest struct {
	name   string
	parser *kong.Kong
	want   []string
	line   string
}

func setLineAndPoint(t *testing.T, line string) func() {
	t.Helper()
	origLine, hasOrigLine := os.LookupEnv(envLine)
	origPoint, hasOrigPoint := os.LookupEnv(envPoint)
	require.NoError(t, os.Setenv(envLine, line))
	require.NoError(t, os.Setenv(envPoint, strconv.Itoa(len(line))))
	return func() {
		t.Helper()
		require.NoError(t, os.Unsetenv(envLine))
		require.NoError(t, os.Unsetenv(envPoint))
		if hasOrigLine {
			require.NoError(t, os.Setenv(envLine, origLine))
		}
		if hasOrigPoint {
			require.NoError(t, os.Setenv(envPoint, origPoint))
		}
	}
}

func runComplete(t *testing.T, parser *kong.Kong, line string, options []Option) []string {
	t.Helper()
	options = append(options,
		WithErrorHandler(func(err error) {
			t.Helper()
			assert.NoError(t, err)
		}),
		WithExitFunc(func(code int) {
			t.Helper()
			assert.Equal(t, 0, code)
		}),
	)
	cleanup := setLineAndPoint(t, line)
	defer cleanup()
	var buf bytes.Buffer
	if parser != nil {
		parser.Stdout = &buf
	}
	Register(parser, options...)
	return parseOutput(buf.String())
}

func parseOutput(output string) []string {
	lines := strings.Split(output, "\n")
	options := []string{}
	for _, l := range lines {
		if l != "" {
			options = append(options, l)
		}
	}
	return options
}
