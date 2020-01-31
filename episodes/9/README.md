# TBS[8]: Would you like some k3sup on that?

In this episode we are going to have some fun with the new Crossplane
[k3sup](https://github.com/alexellis/k3sup) integration! We will:
- Setup a k3s cluster and install Crossplane (using k3sup)
- Provision a second k3s cluster on a Raspberry Pi (using k3sup)
- Expose the kube-apiserver on the Rasberry Pi to public internet using [inlets-pro](https://github.com/inlets/inlets-pro)
- Expose an endpoint for the Raspberry Pi cluster via the
  [inlets-operator](https://github.com/inlets/inlets-operator) and an exit node
  on [Packet](https://www.packet.com/)
- Schedule workloads from the Crossplane cluster to the Raspberry Pi

Host: [@hasheddan](https://twitter.com/hasheddan)

Live Stream: https://youtu.be/RVAFEAnirZA

## Time References

* [Start](https://youtu.be/RVAFEAnirZA?t=78)
* [Create GCP VM for Control Cluster](https://youtu.be/RVAFEAnirZA?t=371)
* [Install k3s on Control Cluster using k3sup](https://youtu.be/RVAFEAnirZA?t=683)
* [Install Crossplane on Control Cluster using k3sup](https://youtu.be/RVAFEAnirZA?t=796)
* [Install k3s on Raspberry Pi using k3sup](https://youtu.be/RVAFEAnirZA?t=1008)
* [Install inlets-operator on Raspberry Pi k3s cluster using k3sup](https://youtu.be/RVAFEAnirZA?t=1341)
* [Setup GCP VM for Rasberry Pi kube-apiserver exit node](https://youtu.be/RVAFEAnirZA?t=1571)
* [Start inlets-pro server on exit node](https://youtu.be/RVAFEAnirZA?t=1787)
* [Start inlets-pro client on Raspberry Pi](https://youtu.be/RVAFEAnirZA?t=2143)
* [Inject Raspberry Pi cluster kubeconfig into Control Cluster Secret](https://youtu.be/RVAFEAnirZA?t=2420)
* [Schedule Namespace to Raspberry Pi using Crossplane in the Control Cluster](https://youtu.be/RVAFEAnirZA?t=2947)

## Media References

* k3sup: <https://github.com/alexellis/k3sup>
* inlets: <https://inlets.dev/>
* inlets-pro: <https://github.com/inlets/inlets-pro>
* Opening and closing theme by [Daniel Suskin](https://soundcloud.com/suskin)