# Slack Message Kubernetes Controller

This is a simple Kubernetes controller built using [Kubebuilder](https://github.com/crossplaneio/crossplane-runtime) and [Crossplane-Runtime](https://github.com/crossplaneio/crossplane-runtime).

## Running Locally

Execute the following steps to run the Slack controller locally on a KIND cluster:

1. Make sure you have [installed KIND](https://github.com/kubernetes-sigs/kind#installation-and-usage).
2. Create a new cluster: `kind create cluster`.
3. Set your `kubectl` configuration to talk to your KIND cluster: `export KUBECONFIG="$(kind get kubeconfig-path --name="kind")"`.
4. Install the Slack CRD's: `make install`.
5. Create `token.txt` in `config/tbs/slack`.
6. Go to [https://api.slack.com/custom-integrations/legacy-tokens](https://api.slack.com/custom-integrations/legacy-tokens) and create an API token for your chosen Workspace.
7. Paste the API token into `token.txt`.
8. Start the Slack controller: `make run`. This will run in the foreground, so you may want to open a new terminal window. If you do, make sure to `export KUBECONFIG` in that session as well.
9. Create a `Secret` and `Provider` object populated with your API token: `./config/tbs/slack/provider.sh`.
10. Edit `config/tbs/slack/message.yaml` with your desired `message` and `channel` to post to. You can find your channel ID by navigating to the channel in your browser and it will be the last part of the URL path.
11. Create your `Message` object: `kubectl apply -f config/tbs/slack/message.yaml`

You should see a message posted to your Slack channel, and the controller should report successful reconciliation in your logs!

## Create Your Own

This project was bootstrapped using Kubebuilder. You can copy this project and modify it, or setup your own using the [Kubebuilder book](https://book.kubebuilder.io/quick-start.html).