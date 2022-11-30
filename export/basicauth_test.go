package export

//Mat-Nr. 8689159
import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createServer(auth AuthenticatorFunc) *httptest.Server {
	return httptest.NewServer(
		WrapperAuth(auth,
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "Hello client")
			}))
}

func TestWithoutPW(t *testing.T) {
	ts := createServer(func(name, pwd string) (bool, string) {
		return true, "test"
	})
	defer ts.Close()
	res, err := http.Get(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode, "wrong status")
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t,
		http.StatusText(http.StatusUnauthorized)+"\n",
		string(body), "wrong message")
}

func doRequestWithPassword(t *testing.T, url string) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)
	req.SetBasicAuth("<username>", "<password>")
	res, err := client.Do(req)
	assert.NoError(t, err)
	return res
}

func TestWithWrongPW(t *testing.T) {
	var receivedName, receivedPwd string
	ts := createServer(func(name, pwd string) (bool, string) {
		receivedName = name
		receivedPwd = pwd
		return false, "test" // <--- deny any request
	})
	defer ts.Close()
	res := doRequestWithPassword(t, ts.URL)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode, "wrong status")
	assert.Equal(t, "<username>", receivedName, "wrong username")
	assert.Equal(t, "<password>", receivedPwd, "wrong password")
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t,
		http.StatusText(http.StatusUnauthorized)+"\n",
		string(body), "wrong message")
}

func TestWithCorrectPW(t *testing.T) {
	var receivedName, receivedPwd string
	ts := createServer(func(name, pwd string) (bool, string) {
		receivedName = name
		receivedPwd = pwd
		return true, "test" // <--- accept any request
	})
	defer ts.Close()
	res := doRequestWithPassword(t, ts.URL)
	assert.Equal(t, http.StatusOK, res.StatusCode, "wrong status code")
	assert.Equal(t, "<username>", receivedName, "wrong username")
	assert.Equal(t, "<password>", receivedPwd, "wrong password")
	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Hello client\n", string(body), "wrong message")
}
