/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package clients

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/linode/linodego"
	"golang.org/x/oauth2"
)

// NewClient returns a new Client
func NewClient(credentials []byte) linodego.Client {
	var apiKey string
	if credentials == nil {
		var ok bool
		apiKey, ok = os.LookupEnv("LINODE_TOKEN")
		if !ok {
			log.Fatal("No credentials provided and LINODE_TOKEN not set.")
		}
	} else {
		apiKey = strings.TrimSpace(string(credentials))
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiKey})
	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
		},
	}

	client := linodego.NewClient(oauth2Client)

	return client
}
