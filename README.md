# CCE Cloud Controller Manager(CCM)

**cce-cloud-controller-manager** is a Kubernetes Cloud Controller Manager implementation (or out-of-tree cloud-provider) for Cloud-Container-Engine(CCE) on BCE.

## Introduction

External cloud providers were introduced as an Alpha feature in Kubernetes 1.6 with the addition of the Cloud Controller Manager binary. External cloud providers are Kubernetes (master) controllers that implement the cloud-provider specific control loops required for Kubernetes to function.

This functionality is implemented in-tree in the kube-controller-manger binary for existing cloud-providers (e.g. AWS, GCE, etc.), however, in-tree cloud-providers have entered maintenance mode and no additional providers will be accepted. Furthermore, there is an ongoing effort to remove all existing cloud-provider specific code out of the Kubernetes codebase.

## Setup and Installation

See docs/cce-cloud-contoller-manager.md

## Example

See docs/tutorial.md

## License

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.