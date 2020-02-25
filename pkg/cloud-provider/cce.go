/*
Copyright 2018 The Kubernetes Authors.

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
package cloud_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/bce"
	bcesdk "icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/bce"
	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/blb"
	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/eip"
	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/vpc"
	cce "icode.baidu.com/baidu/jpaas-caas/cloud-provider-baiducloud/pkg/temp-cce"
)

// ProviderName is the name of this cloud provider.
const ProviderName = "cce"

// CceUserAgent is prefix of http header UserAgent
const (
	// CCEUserAgent is prefix of http header UserAgent
	CCEUserAgent = "cce-k8s"
)

const (
	// How long to wait before retrying the processing of a service change.
	// If this changes, the sleep in hack/jenkins/e2e.sh before downing a cluster
	// should be changed appropriately.
	minRetryDelay = 5 * time.Second
	maxRetryDelay = 300 * time.Second
)

// Baiducloud defines the main struct
type Baiducloud struct {
	CloudConfig
	clientSet        *ClientSet
	kubeClient       kubernetes.Interface
	eventBroadcaster record.EventBroadcaster
	eventRecorder    record.EventRecorder
	// services that need to be synced
	svcQueue workqueue.RateLimitingInterface
}

// CloudConfig is the cloud config
type CloudConfig struct {
	ClusterID       string `json:"ClusterId"`
	ClusterName     string `json:"ClusterName"`
	AccessKeyID     string `json:"AccessKeyID"`
	SecretAccessKey string `json:"SecretAccessKey"`
	Region          string `json:"Region"`
	VpcID           string `json:"VpcId"`
	SubnetID        string `json:"SubnetId"`
	MasterID        string `json:"MasterId"`
	Endpoint        string `json:"Endpoint"`
	NodeName        string `json:"NodeName"`
	Debug           bool   `json:"Debug"`
}

// CCMVersion is the version of CCM
var CCMVersion string

func init() {
	cloudprovider.RegisterCloudProvider(ProviderName, func(configReader io.Reader) (cloudprovider.Interface, error) {
		var cloud Baiducloud
		var cloudConfig CloudConfig
		configContents, err := ioutil.ReadAll(configReader)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(configContents, &cloudConfig)
		if err != nil {
			return nil, err
		}
		klog.Infof("Init CCE cloud with cloudConfig: %v\n", cloudConfig)
		if cloudConfig.MasterID == "" {
			return nil, fmt.Errorf("Cloud config must have a Master ID\n ")
		}
		if cloudConfig.ClusterID == "" {
			return nil, fmt.Errorf("Cloud config must have a ClusterID\n ")
		}
		if cloudConfig.Endpoint == "" {
			return nil, fmt.Errorf("Cloud config must have a Endpoint\n ")
		}

		cloud.CloudConfig = cloudConfig
		cloud.clientSet, err = newClientSet(&cloudConfig)
		if err != nil {
			return nil, err
		}
		return &cloud, nil
	})
}

// ProviderName returns the cloud provider ID.
func (bc *Baiducloud) ProviderName() string {
	return ProviderName
}

// Initialize provides the cloud with a kubernetes client builder and may spawn goroutines
// to perform housekeeping activities within the cloud provider.
func (bc *Baiducloud) Initialize(clientBuilder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {
	bc.kubeClient = clientBuilder.ClientOrDie(ProviderName)
	bc.eventBroadcaster = record.NewBroadcaster()
	bc.eventBroadcaster.StartLogging(klog.Infof)
	bc.eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: bc.kubeClient.CoreV1().Events("")})
	bc.eventRecorder = bc.eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: "CCM"})
	bc.svcQueue = workqueue.NewNamedRateLimitingQueue(workqueue.NewItemExponentialFailureRateLimiter(minRetryDelay, maxRetryDelay), "endpoints")
	bc.runServiceWorker()
}

// SetInformers sets the informer on the cloud object.
func (bc *Baiducloud) SetInformers(informerFactory informers.SharedInformerFactory) {
	klog.Infof("Setting up informers for Baiducloud")
	// node
	nodeInformer := informerFactory.Core().V1().Nodes().Informer()
	nodeInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			node := obj.(*v1.Node)
			klog.Infof("nodeInformer node add: %v", node.Name)
			// TODO: cache some node info
		},
		UpdateFunc: func(prev, obj interface{}) {
			node := obj.(*v1.Node)
			klog.Infof("nodeInformer node update: %v", node.Name)
		},
		DeleteFunc: func(obj interface{}) {
			node := obj.(*v1.Node)
			klog.Infof("nodeInformer node delete: %v", node.Name)
			// TODO: remove node info from cache
		},
	})
	// service
	serviceInformer := informerFactory.Core().V1().Services().Informer()
	serviceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			service := obj.(*v1.Service)
			klog.Infof("serviceInformer service add: %v", service.Name)
			// TODO: cache some service info
		},
		UpdateFunc: func(prev, obj interface{}) {
			service := obj.(*v1.Service)
			klog.Infof("serviceInformer service update: %v", service.Name)
			// TODO:
		},
		DeleteFunc: func(obj interface{}) {
			service := obj.(*v1.Service)
			klog.Infof("serviceInformer service delete: %v", service.Name)
			// TODO: remove service info from cache
		},
	})

	// endpoints
	podInformer := informerFactory.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, newObj interface{}) {
			pod := newObj.(*v1.Pod)

			if len(pod.Status.PodIP) == 0 {
				return
			}

			svcs, err := bc.kubeClient.CoreV1().Services(pod.Namespace).List(metav1.ListOptions{})
			if err != nil {
				klog.Errorf("podInformer failed to get service of pod %s/%s: %s", pod.Namespace, pod.Name, err)
				return
			}
			for _, svc := range svcs.Items {
				if svc.Spec.Selector == nil {
					continue
				}
				if svc.Spec.ExternalTrafficPolicy != v1.ServiceExternalTrafficPolicyTypeLocal {
					klog.Infof("Service %s/%s is not externalTrafficPolicy Local, skipping update rs", svc.Namespace, svc.Name)
					continue
				}
				selector := labels.Set(svc.Spec.Selector).AsSelectorPreValidated()
				if selector.Matches(labels.Set(pod.Labels)) {
					key := fmt.Sprintf("%s/%s", svc.Namespace, svc.Name)
					bc.svcQueue.AddRateLimited(key)
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)

			if len(pod.Status.PodIP) == 0 {
				return
			}

			svcs, err := bc.kubeClient.CoreV1().Services(pod.Namespace).List(metav1.ListOptions{})
			if err != nil {
				klog.Errorf("podInformer failed to get service of pod %s/%s: %s", pod.Namespace, pod.Name, err)
				return
			}
			for _, svc := range svcs.Items {
				if svc.Spec.Selector == nil {
					continue
				}
				if svc.Spec.ExternalTrafficPolicy != v1.ServiceExternalTrafficPolicyTypeLocal {
					klog.Infof("Service %s/%s is not externalTrafficPolicy Local, skipping update rs", svc.Namespace, svc.Name)
					continue
				}
				selector := labels.Set(svc.Spec.Selector).AsSelectorPreValidated()
				if selector.Matches(labels.Set(pod.Labels)) {
					key := fmt.Sprintf("%s/%s", svc.Namespace, svc.Name)
					bc.svcQueue.AddRateLimited(key)
				}
			}
		},
	})
}

// ClientSet contains all the bce product client
type ClientSet struct {
	BLBClient blb.Interface
	EIPClient eip.Interface
	CCEClient cce.Interface
	VPCClient vpc.Interface
}

func newClientSet(config *CloudConfig) (*ClientSet, error) {
	if config == nil {
		return nil, fmt.Errorf("newClientSet failed: config is nil")
	}

	clientset := &ClientSet{}

	// set cce-gateway proxy
	proxyHost, proxyPort := getCCEGatewayHostAndPort(config.Region)
	// BLBClient
	lbClient := blb.NewBLBClient(&blb.Config{
		Config: &bcesdk.Config{
			Credentials: bcesdk.NewCredentials(config.AccessKeyID, config.SecretAccessKey),
			Checksum:    true,
			Timeout:     30 * time.Second,
			Region:      config.Region,
			Endpoint:    blb.Endpoint[config.Region], // notice!
			UserAgent:   fmt.Sprintf("%s:%s", CCEUserAgent, config.ClusterID),
		},
	})
	clientset.BLBClient = lbClient

	// EIPClient
	eipClient := eip.NewClient(&eip.Config{
		Config: &bcesdk.Config{
			Credentials: bcesdk.NewCredentials(config.AccessKeyID, config.SecretAccessKey),
			Checksum:    true,
			Timeout:     30 * time.Second,
			Region:      config.Region,
			Endpoint:    eip.Endpoint[config.Region],
			ProxyHost:   proxyHost,
			ProxyPort:   proxyPort,
		},
	})
	clientset.EIPClient = eipClient

	// CCEClient request Internal API
	cceClient := cce.NewClient(&cce.Config{
		Config: &bcesdk.Config{
			Credentials: bcesdk.NewCredentials(config.AccessKeyID, config.SecretAccessKey),
			Checksum:    true,
			Timeout:     30 * time.Second,
			Region:      config.Region,
			Endpoint:    config.Endpoint,
			UserAgent:   fmt.Sprintf("%s:%s", CCEUserAgent, config.ClusterID), // UserAgent
			ProxyHost:   proxyHost,
			ProxyPort:   proxyPort,
		},
	})
	clientset.CCEClient = cceClient

	// VPCClient
	vpcClient := vpc.NewClient(&vpc.Config{
		Config: &bcesdk.Config{
			Credentials: bcesdk.NewCredentials(config.AccessKeyID, config.SecretAccessKey),
			Checksum:    true,
			Timeout:     30 * time.Second,
			Region:      config.Region,
			Endpoint:    vpc.Endpoint[config.Region],
			ProxyHost:   proxyHost,
			ProxyPort:   proxyPort,
		},
	})
	clientset.VPCClient = vpcClient

	// Set Debug
	config.Debug = true
	if config.Debug == true {
		klog.Info("cce-ingresss-controller set debug = true")
	}
	lbClient.SetDebug(config.Debug)
	eipClient.SetDebug(config.Debug)
	cceClient.SetDebug(config.Debug)
	vpcClient.SetDebug(config.Debug)

	return clientset, nil
}

func getCloudConfig(ctx context.Context) (*CloudConfig, error) {
	ccBytes, err := ioutil.ReadFile("/etc/kubernetes/cloud.config")
	if err != nil {
		return nil, fmt.Errorf("read cloud config error: '%v'", err)
	}

	cc := new(CloudConfig)
	err = json.Unmarshal(ccBytes, cc)
	if err != nil {
		return nil, fmt.Errorf("unmarshal cloud config error: '%v'", err)
	}

	return cc, nil
}

// NewCCEClient for internal cce service
func NewCCEClient(accessKeyID, secretAccessKey, region, endpoint string) *cce.Client {
	return cce.NewClient(&cce.Config{
		&bce.Config{
			Credentials: bce.NewCredentials(accessKeyID, secretAccessKey),
			Checksum:    true,
			Timeout:     30 * time.Second,
			Region:      region,
			Endpoint:    endpoint,
		},
	})
}

func (bc *Baiducloud) runServiceWorker() {
	stopCh := make(chan struct{})
	for i := 0; i < 10; i++ {
		go wait.Until(bc.serviceWorker, time.Second, stopCh)
	}
}

func (bc *Baiducloud) serviceWorker() {
	for bc.processNextService() {

	}
}

func (bc *Baiducloud) processNextService() bool {
	ctx := context.WithValue(context.Background(), RequestID, GetRandom())
	key, quit := bc.svcQueue.Get()
	if quit {
		return false
	}
	defer bc.svcQueue.Done(key)
	klog.Infof(Message(ctx, fmt.Sprintf("Pod changed, begin reconcile backend server for service %s", key)))

	err := func() error {
		namespace, name, err := cache.SplitMetaNamespaceKey(key.(string))
		if err != nil {
			runtime.HandleError(fmt.Errorf("Invalid resource key: %s", key))
			return err
		}
		service, err := bc.kubeClient.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		nodes := make([]*v1.Node, 0)
		return bc.reconcileBackendServers(ctx, bc.ClusterName, service, nodes)
	}()
	if err == nil {
		bc.svcQueue.Forget(key)
		return true
	}

	runtime.HandleError(fmt.Errorf("error processing service %v (will retry): %v", key, err))
	bc.svcQueue.AddRateLimited(key)
	return true
}

func getCCEGatewayHostAndPort(region string) (string, int) {
	// default to bj
	host := "xxxxxxxxxxxxxxx"
	if region != "" {
		host = "xxxxxxxxx" + region + "xxxxxxxxxxxx"
	}
	if env := os.Getenv("CCE_GATEWAY_HOST"); env != "" {
		// host is explicitly set in sandbox
		host = env
	}
	return host, 0
}
