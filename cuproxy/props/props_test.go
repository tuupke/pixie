package props

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseToCallString(t *testing.T) {
	res, err := parseToCallString(strings.Join([]string{
		`team;GET;https://domjudge.org/demoweb/api/|user;DELETE;https://domjudge.org/demoweb/api/`,
		`pixie;POST;https://localhost:9000/`,
	}, "&&"))

	require.Nil(t, err)
	require.EqualValues(t, endpointsSet{
		endpoints{
			endpoint{
				method: http.MethodGet,
				name:   "team",
				url:    "https://domjudge.org/demoweb/api/",
			},
			endpoint{
				method: http.MethodDelete,
				name:   "user",
				url:    "https://domjudge.org/demoweb/api/",
			},
		},
		endpoints{
			endpoint{
				method: http.MethodPost,
				name:   "pixie",
				url:    "https://localhost:9000/",
			}},
	}, res)
}

func TestLoad(t *testing.T) {
	prependPropsWithName = true
	toCall = endpointsSet{
		endpoints{
			endpoint{
				method: http.MethodGet,
				name:   "user",
				url:    "https://jury:jury@www.domjudge.org/demoweb/api/user?strict=false",
			},
			endpoint{
				method: http.MethodGet,
				name:   "user",
				url:    "https://jury:jury@www.domjudge.org/demoweb/api/teams/{{user_team_id}}?strict=false",
			},
		},
		endpoints{
			endpoint{
				method: http.MethodGet,
				name:   imageKey,
				url:    "https://www.w3.org/MarkUp/Test/xhtml-print/20050519/tests/jpeg444.jpg",
			},
		},
	}

	ip := net.ParseIP("127.0.0.1")
	firstTime := time.Now()
	Load(ip)
	first := time.Since(firstTime)
	secondTime := time.Now()
	p := Load(ip)
	second := time.Since(secondTime)

	assert.Greater(t, first, second)
	fmt.Println(firstTime, secondTime, first, second)
	bts, err := io.ReadAll(p.json())
	fmt.Println(err, string(bts))
}

func BenchmarkRefresh(b *testing.B) {
	b.StopTimer()
	prependPropsWithName = true
	toCall = endpointsSet{
		endpoints{
			endpoint{
				method: http.MethodGet,
				name:   "user",
				url:    "https://jury:jury@www.domjudge.org/demoweb/api/user?strict=false",
			},
			endpoint{
				method: http.MethodGet,
				name:   "user",
				url:    "https://jury:jury@www.domjudge.org/demoweb/api/teams/{{user_team_id}}?strict=false",
			},
		},
		endpoints{
			endpoint{
				method: http.MethodGet,
				name:   imageKey,
				url:    "https://www.w3.org/MarkUp/Test/xhtml-print/20050519/tests/jpeg444.jpg",
			},
		},
	}

	ip := net.ParseIP("127.0.0.1")
	Load(ip)
	for i := 0; i < b.N; i++ {

	}

}

func TestToJson(t *testing.T) {
	p := Load(nil)

	doCheck := func(expected map[string]string) {
		var found map[string]string
		err := json.NewDecoder(p.json()).Decode(&found)

		require.Nil(t, err)
		require.NotNil(t, found)
		assert.EqualValues(t, expected, found)
	}

	doCheck(make(map[string]string))

	// Set a key
	p.data.Store("foo", "bar")
	doCheck(map[string]string{"foo": "bar"})

	// Set another key, note the second check only serves to illustrate that order
	// does not matter.
	p.data.Store("foobar", "baz")
	doCheck(map[string]string{"foo": "bar", "foobar": "baz"})
	doCheck(map[string]string{"foobar": "baz", "foo": "bar"})
}
