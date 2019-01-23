package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
	"time"

	clientset "github.com/kongyi-ibm/k8s-deployment-operator/pkg/client/clientset/versioned"
	deploycontrschema "github.com/kongyi-ibm/k8s-deployment-operator/pkg/client/clientset/versioned/scheme"
	deploycontrinformer "github.com/kongyi-ibm/k8s-deployment-operator/pkg/client/informers/externalversions/deploycontrol/v1alpha1"
	deploycontrlisters "github.com/kongyi-ibm/k8s-deployment-operator/pkg/client/listers/deploycontrol/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	appsinformer "k8s.io/client-go/informers/apps/v1"
	corev1informer "k8s.io/client-go/informers/core/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	"github.com/kongyi-ibm/k8s-deployment-operator/pkg/apis/deploycontrol/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Controller struct {

	kubeclientset kubernetes.Interface

	extclientset  clientset.Interface

	deploymentsLister appslisters.DeploymentLister

	deploymentsSynced cache.InformerSynced

	deploydaemonLister deploycontrlisters.DeployDaemonLister

	deploydaemonSynced cache.InformerSynced

	podslister corev1listers.PodLister

	podsSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

const controllerAgentName = "deploydaemon-controller"

func MyNewController(name string){
	fmt.Println("Hello %s", name)
}

func NewController(
	   kubeclientset kubernetes.Interface,
	   extclientset clientset.Interface,
	   deploymentInformer appsinformer.DeploymentInformer,
	   podInformer corev1informer.PodInformer,
	   deploydaemonInformer deploycontrinformer.DeployDaemonInformer) *Controller {

    // Create event broadcaster
    // Add deploycontrol types to the default Kubernetes Scheme so Events can be
    // logged for deploycontrol types.
    utilruntime.Must(deploycontrschema.AddToScheme(scheme.Scheme))

    klog.V(4).Info("Creating event broadcaster")
    eventBroadcaster := record.NewBroadcaster()
    eventBroadcaster.StartLogging(klog.Infof)
    eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
    recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

    controller := &Controller{
            kubeclientset:       kubeclientset,
		    extclientset:        extclientset,
		    deploymentsLister:   deploymentInformer.Lister(),
            deploymentsSynced:   deploymentInformer.Informer().HasSynced,
		    podslister:          podInformer.Lister(),
		    podsSynced:          podInformer.Informer().HasSynced,
		    deploydaemonLister:  deploydaemonInformer.Lister(),
		    deploydaemonSynced:  deploydaemonInformer.Informer().HasSynced,
            workqueue:           workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "DeployDaemons"),
            recorder:            recorder,
    }

	klog.Info("Setting up event handlers for deploydaemon")

	// Set up an event handler for when DeployDaemon resources change
	deploydaemonInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueDeployDaemon,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueDeployDaemon(new)
		},
	})

    return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func ( c *Controller ) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting Foo controller")

	// Waiting for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")

	if ok:=cache.WaitForCacheSync(stopCh, c.deploydaemonSynced,c.podsSynced,c.deploymentsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch two workers to process Foo resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error{

		defer c.workqueue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}

		if err := c.reconcile(key); err !=nil {
			return fmt.Errorf("sync DeployDaemon %s failed: %s", key, err.Error())
		}

		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		//utilruntime.HandleError(err)
		klog.Errorf("Error: %s", err.Error())
		return true
	}
    return true
}

//TODO: Create new package for Reconciler. So far, wrap most logic in this function
func (c *Controller ) reconcile(key string) error {

	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		klog.Errorf("invalid resource key: %s", key)
		return nil
	}

	deploydaemon, err:=c.deploydaemonLister.DeployDaemons(namespace).Get(name)
	if err != nil {
		klog.Errorf("get crd object %s failed", key)
		return nil
	}

	// If the build's done, then ignore it.
	if c.isDone(deploydaemon.Status) {
		return nil
	}

	// All obj in index is ready only
	deploydaemon = deploydaemon.DeepCopy()

	var dp *appsv1.Deployment

	// If below status is empty, that means haven't create deployment.
	if deploydaemon.Status == nil || deploydaemon.Status.Cluster.DeploymentName == "" {

         dp, err = c.createDeployent(deploydaemon)
         if err !=nil {
         	klog.Errorf("create deployment %s for deploydaemon %s ", deploydaemon.GenerateName, key )
         	// Here we need the to fill in the conditionSpec with Reason that the Deployment can not be ceate
         	// then update the deploydaemon status -- this will trigger the update event and make the deploymentdaemon be add into this controller loop again.
		 }

	} else {

		// Else, we need to check deployment status to update deploydaemon status
		// 检查 deployment 的状态是否是ready的状态
		// 如果 deployment 不是ready的状态的话就需要把这个deploydaemon再重新加入到 workqueue中，并更新lastUpdateTime
		// 如果 ready的状态 check 对应的Pod的 expose label是否和 deploydaemon符合，如果不符合就修改
		dp, err =c.deploymentsLister.Deployments(namespace).Get(deploydaemon.GenerateName)
		if err != nil {
			klog.Errorf("can not get deployment %s in namespace %s", deploydaemon.GenerateName, namespace)
			// Here we need the to fill in the conditionSpec with Reason that the Deployment can not be get
			// then we update the deploydaemon
		}

	}

	c.updateDeployDaemonStatus(deploydaemon, dp)

	klog.Infof("Check object info: .........")
	klog.Info(deploydaemon)
	return nil
}


func ( c *Controller) enqueueDeployDaemon(obj interface{}){
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

func ( c *Controller ) updateDeployDaemonStatus(deploydaemon *v1alpha1.DeployDaemon, deployment *appsv1.Deployment) error {

	klog.Infof("updateDeployDaemonStatus for deploydaemon %s: ", deploydaemon.Name)

	return nil
}

func ( c *Controller ) createDeployent(deploydaemon *v1alpha1.DeployDaemon) (*appsv1.Deployment, error) {

	dp ,err := c.makeDeployment(deploydaemon)
	if err !=nil {
		klog.Errorf("Faile to create deployment: %s", err.Error())
	}

	klog.Infof("createDeployent for deploydaemon %s: ", deploydaemon.Name)

	return c.kubeclientset.AppsV1().Deployments(deploydaemon.Namespace).Create(dp)
}

func (c *Controller ) makeDeployment(deploydaemon *v1alpha1.DeployDaemon) (*appsv1.Deployment, error){
	klog.Infof("make deployment for deploydaemon: %s", deploydaemon.Name)
	//&corev1.Deployment{
	//	ObjectMeta: metav1.ObjectMeta{
	//		// We execute the build's pod in the same namespace as where the build was
	//		// created so that it can access colocated resources.
	//		Namespace: deploydaemon.Namespace,
	//		// Generate a unique name based on the build's name.
	//		// Add a unique suffix to avoid confusion when a build
	//		// is deleted and re-created with the same name.
	//		// We don't use GenerateName here because k8s fakes don't support it.
	//		Name: deploydaemon.GetDeploymentName(),
	//		// If our parent Build is deleted, then we should be as well.
	//		OwnerReferences: []metav1.OwnerReference{
	//			*metav1.NewControllerRef(deploydaemon, schema.GroupVersionKind{
	//				Group:   v1alpha1.SchemeGroupVersion.Group,
	//				Version: v1alpha1.SchemeGroupVersion.Version,
	//				Kind:    "DeployDaemon",
	//			}),
	//		},
	//		Labels: map[string]string{
	//			"expose": deploydaemon.Spec.Expose,
	//			"app": deploydaemon.GetDeploymentName(),
	//		},
	//	},
	//	Spec: corev1.PodSpec{
	//		// If the build fails, don't restart it.
	//		RestartPolicy:  corev1.RestartPolicyNever,
	//
	//		Containers: []corev1.Container{{
	//			Name:  "Nginx-Deployment",
	//			Image: "docker.io/nginx:latest",
	//		}},
	//	},
	//}, nil
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
					// We execute the build's pod in the same namespace as where the build was
					// created so that it can access colocated resources.
					Namespace: deploydaemon.Namespace,
					// Generate a unique name based on the build's name.
					// Add a unique suffix to avoid confusion when a build
					// is deleted and re-created with the same name.
					// We don't use GenerateName here because k8s fakes don't support it.
					Name: deploydaemon.GetDeploymentName(),
					// If our parent Build is deleted, then we should be as well.
					OwnerReferences: []metav1.OwnerReference{
						*metav1.NewControllerRef(deploydaemon, schema.GroupVersionKind{
							Group:   v1alpha1.SchemeGroupVersion.Group,
							Version: v1alpha1.SchemeGroupVersion.Version,
							Kind:    "DeployDaemon",
						}),
					},
					Labels: map[string]string{
						"expose": deploydaemon.Spec.Expose,
						"app": deploydaemon.GetDeploymentName(),
					},
				},
		Spec: appsv1.DeploymentSpec{
			Replicas: deploydaemon.Spec.Replica,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"expose": deploydaemon.Spec.Expose,
					"app": deploydaemon.GetDeploymentName(),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"expose": deploydaemon.Spec.Expose,
						"app": deploydaemon.GetDeploymentName(),
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
						},
					},
				},
			},
		},
	},nil
}

func ( c *Controller ) isDone(status *v1alpha1.DeploydaemonStatus) bool {
	// judge is deploydaemon is done
	if status == nil || status.Conditions.Status == false {
		return false
	}
	return true
}