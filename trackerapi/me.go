package trackerapi

import (
	"log"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	u "os/user"

	"github.com/dangogh/clirescue/cmdutil"
	"github.com/dangogh/clirescue/user"
)

const url          string     = "https://www.pivotaltracker.com/services/v5/me"
var (
	tokenFile  string     = homeDir() + "/.tracker"
	currentUser  *user.User = user.New()
	Stdout       *os.File   = os.Stdout
)

func Me() {
	tok, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		// assume file not found and use basic auth
		// TODO: exit with error if other than file not found
		setCredentials()
		resp, err := makeRequest()
		if err != nil {
			log.Fatal(err)
		}
		tok, err = getToken(resp)
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile(tokenFile, tok, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	currentUser.APIToken = string(tok)
}

func makeRequest() ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(currentUser.Username, currentUser.Password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\n****\nAPI response: \n%s\n", string(body))
	return body, nil
}

func getToken(body []byte) ([]byte, error) {
	var meResp = new(MeResponse)
	err := json.Unmarshal(body, &meResp)
	if err != nil {
		return nil, err
	}
	return []byte(meResp.APIToken), nil
}

func setCredentials() {
	fmt.Fprint(Stdout, "Username: ")
	var username = cmdutil.ReadLine()
	cmdutil.Silence()
	fmt.Fprint(Stdout, "Password: ")

	var password = cmdutil.ReadLine()
	currentUser.Login(username, password)
	cmdutil.Unsilence()
}

func homeDir() string {
	usr, _ := u.Current()
	return usr.HomeDir
}

type MeResponse struct {
	APIToken string `json:"api_token"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Initials string `json:"initials"`
	Timezone struct {
		Kind      string `json:"kind"`
		Offset    string `json:"offset"`
		OlsonName string `json:"olson_name"`
	} `json:"time_zone"`
}
