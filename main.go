package main

import (
	"flag"
	clientset "github.com/kongyi-ibm/k8s-deployment-operator/pkg/client/clientset/versioned"
	extInformers "github.com/kongyi-ibm/k8s-deployment-operator/pkg/client/informers/externalversions"
	"github.com/kongyi-ibm/k8s-deployment-operator/pkg/signals"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"time"
)



//import (
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	appsv1 "k8s.io/api/apps/v1"
//)
//
//func main() {
//	JsonParse := NewJsonStruct()
//	deployment := appsv1.Deployment{}
//	//下面使用的是相对路径，config.json文件和main.go文件处于同一目录下
//	JsonParse.Load("./deployment.json", &deployment)
//	fmt.Println("start output parsed deployment")
//	fmt.Println(deployment.Spec.Replicas)
//	//fmt.Println(deployment.Spec.Replicas)
//	//fmt.Println(deployment.Name)
//}
//
//type JsonStruct struct {
//}
//
//func NewJsonStruct() *JsonStruct {
//	return &JsonStruct{}
//}
//
//func (jst *JsonStruct) Load(filename string, v interface{}) {
//	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
//	data, err := ioutil.ReadFile("/Users/kongyi/Workspace/GOlang/deploydaemon/src/github.com/kongyi-ibm/k8s-deployment-operator/deployment.json")
//	fmt.Println(data)
//	if err != nil {
//		return
//	}
//
//	//读取的数据为json格式，需要进行解码
//	err = json.Unmarshal(data, v)
//	if err != nil {
//		return
//	}
//}


var (
	masterURL string
	kubeconfig string
)

func main() {

	//TODO:  make those paramter as the input parameter for Debug flag
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "/tmp")
	flag.Set("v", "4")

	flag.Parse()

	klog.V(0).Info("Start deploy daemon server .....")

	stopCh := signals.SetupSignalHandler()
	// Based on input parameter to generate kube configuration
	cfg,err := clientcmd.BuildConfigFromFlags(masterURL,kubeconfig)

	if err !=nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	// Based on kube configuration to generate kube client
	kubeClient, err := kubernetes.NewForConfig(cfg)

	if err !=nil {
		klog.Fatalf("Error generate kubeclient: %s",err.Error())
	}

	// Based on extension object api client set to generate extension client
	extClient, err := clientset.NewForConfig(cfg)

	if err !=nil {
		klog.Fatalf("Error build deploycontrol client: %s", err.Error())
	}

	// Generate SharedIndexInformerFactory based on the clientset
	kubeInformerFactory :=kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	extInformerFactory := extInformers.NewSharedInformerFactory(extClient, time.Second*30)


	//NewController(
	//	kubeclientset kubernetes.Interface,
	//	extclientset clientset.Interface,
	//	deploymentInformer appsinformer.DeploymentInformer,
	//	podInformer corev1informer.PodInformer,
	//	deploydaemonInformer deploycontrinformer.DeployDaemonInformer) *Controller

	controller := NewController(kubeClient, extClient,
		kubeInformerFactory.Apps().V1().Deployments(),
		kubeInformerFactory.Core().V1().Pods(),
		extInformerFactory.Deploycontrol().V1alpha1().DeployDaemons())


	kubeInformerFactory.Start(stopCh)
	extInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}

	<-stopCh
}


func init() {
	klog.InitFlags(nil)
	flag.StringVar(&kubeconfig, "kubeconfig", "/Users/kongyi/.kube/config", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}