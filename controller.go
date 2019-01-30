package main

import (
	"crypto/rand"
	"fmt"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	utils "github.com/kongyi-ibm/k8s-deployment-operator/pkg/utilities"
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

	// Use self-defined DelayWithRateLimiteQueue which support self-defined delay time. if not specify, will use ratelimite delay time.
	workqueue *utils.DelayWithRateLimitQueue

	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

const controllerAgentName = "deploydaemon-controller"

// Random byte reader used for pod name generation.
// var for testing.
var randReader = rand.Reader

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
            workqueue:           utils.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "DeployDaemons"),
		    //delayqueue:          workqueue.NewNamedDelayingQueue("DelyQueue"),
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

	//klog.Info("Starting delay workers")
	//go wait.Until(c.runDeployWorker,time.Second,stopCh)

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
			c.workqueue.AddRateLimited(obj)
			return fmt.Errorf("sync DeployDaemon %s failed: %s", key, err.Error())
		}

		c.workqueue.Forget(obj)
		klog.V(1).Infof("Controller move deploydaemon '%s' out from queue", key)
		klog.Infof("Successfully synced '%s'", key)

		return nil
	}(obj)

	if err != nil {
		//utilruntime.HandleError(err)
		klog.Errorf("Error: %s", err.Error())
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

	//TODO: Move isDone to final steps after all check passed. This can cover other parameter changes

	// All obj in index is ready only
	deploydaemon = deploydaemon.DeepCopy()

	var dp *appsv1.Deployment

	////Don't need to handle this, when decided to add the object in workqueue.
	//
	// TODO:  Here we need to add logic to check if this deploydaemon has set scheduler, if yes, if the scehduler time is after current time, it must be added into deployqueue
	if deploydaemon.Spec.Scheduler != "" {
		klog.Info("Remove scheduler when deploydaemon be handled as duration")
		deploydaemon.Spec.Scheduler = ""
	}

	// If below status is empty, that means haven't create deployment.
	if deploydaemon.Status == nil || deploydaemon.Status.Cluster.DeploymentName == "" {

         dp, err = c.createDeployent(deploydaemon)
         if err !=nil {
         	klog.Errorf("create deployment %s for deploydaemon %s ", deploydaemon.GenerateName, key )
         	// Here we need the to fill in the conditionSpec with Reason that the Deployment can not be ceate
         	// then update the deploydaemon status -- this will trigger the update event and make the deploymentdaemon be add into this controller loop again.
		 }
         klog.Info("Waiting status for deployment %s ", deploydaemon.Status.Cluster.DeploymentName)

	} else {

		dp, err =c.deploymentsLister.Deployments(namespace).Get(deploydaemon.Status.Cluster.DeploymentName)
		if err != nil {
			return fmt.Errorf("can not get deployment %s in namespace %s", deploydaemon.GenerateName, namespace)
			// Here we need the to fill in the conditionSpec with Reason that the Deployment can not be get
			// then we update the deploydaemon
		}
		klog.Infof("check deployment %s status : %s", dp.Name,dp.Status)
	}

	// Else, we need to check deployment status to update deploydaemon status
	// 检查 deployment 的状态是否是ready的状态
	// 如果 deployment 不是ready的状态的话就需要把这个deploydaemon再重新加入到 workqueue中，并更新lastUpdateTime
	// 如果 ready的状态 check 对应的Pod的 expose label是否和 deploydaemon符合，如果不符合就修改

	c.syncDeployDaemon(deploydaemon,dp)

	updateErr := c.updateDeployDaemonStatus(deploydaemon)

	klog.Infof("check deploydaemon after update: \n %s",deploydaemon)

	if sycError:=c.isDone(deploydaemon.Status,updateErr); sycError !=nil {
		klog.Errorf("sync status of deploydaemon %s", deploydaemon.Name)
		return sycError
	}
	return nil
}

func (c *Controller) syncDeployDaemon(deploydaemon *v1alpha1.DeployDaemon, deployment *appsv1.Deployment) {
	klog.Infof("sync deployDaemon status for %s: ", deploydaemon.Name)

	if deployment !=nil {
		deploydaemon.Status.Deployment = deployment.Status

		//1. Check deployment, if necessary, need to delete old and create a new one
		if err := c.syncDeployment(deploydaemon, deployment); err != nil {
			c.remarkSuccessStatus(deploydaemon,false, "Waiting Deployment Ready", err.Error())
			return
		}

		//2. Check pod expose status
		if err := c.syncPodExposeStatus(deploydaemon); err !=nil{
			//c.remarkSuccessStatus(deploydaemon,false, "Waiting Pod Expose Sync Ready",err.Error())
			c.remarkSuccessStatus(deploydaemon,false, "Waiting Pod Expose Sync Ready",err.Error())
			return
		}

		//3. If everything ready, we can mark the DeployDaemon is ready, otherwise, mark is not ready
		c.remarkSuccessStatus(deploydaemon,true, "All Status Synced","Deployment success!")
	}else{
		klog.Infof("deployment %s status is nil, waiting check in next around %s: ", deployment.Name)
	}
}

func ( c *Controller) syncDeployment(deploydaemon *v1alpha1.DeployDaemon, deployment *appsv1.Deployment) error {
	klog.Infof("sync deployment status for %s: ", deploydaemon.Name)
	var err error = nil

	//1. Check replica / image / version change and call createDeployment
	// delete deployment first
	// call createDeployment method
	// err = fmt.Errorf("Re-genereate deployment, since significant change on deployment")

	var deploymentUpdated bool = false

	if *deployment.Spec.Replicas != *deploydaemon.Spec.Replica {
		klog.Infof("deployment replica %s not sync with deploydaemon replica %s", *deployment.Spec.Replicas, *deploydaemon.Spec.Replica)
		deployment.Spec.Replicas = deploydaemon.Spec.Replica
		deploymentUpdated = true
		c.kubeclientset.AppsV1().Deployments(deploydaemon.Namespace).Update(deployment)
		err = fmt.Errorf("Waiting Pod Scale Ready")
	}

	//2. Check deployment ready
	if ! deploymentUpdated {
		if deployment.Status.AvailableReplicas != deployment.Status.ReadyReplicas {
			err = fmt.Errorf("Waiting Pod Status Ready")
		}
	}

	return err
}


func ( c *Controller) syncPodExposeStatus(deploydaemon *v1alpha1.DeployDaemon) error {

	klog.Infof("Sync PodExposeStatus for deploydaemon: %s", deploydaemon.Name)

	selector := labels.SelectorFromSet(map[string]string{
		"app": deploydaemon.GetDeploymentName(),
		"version": deploydaemon.Spec.Version,
	})

	podlist, err:= c.podslister.Pods(deploydaemon.Namespace).List(selector)

	completeSync := true

	if err !=nil {
		klog.Errorf("list pod for deployment %s with version failed", deploydaemon.GetDeploymentName())
		return err
	}
	for _, pod := range podlist  {

		if pod.Labels["expose"] != deploydaemon.Spec.Expose {
			klog.Infof("Sync Pod %s Expose To %s Caused By DeployDaemon %s Expose Change ", pod.Name,deploydaemon.Spec.Expose,deploydaemon.Name )
			pod.Labels["expose"]=deploydaemon.Spec.Expose
			_, err =c.kubeclientset.CoreV1().Pods(deploydaemon.Namespace).Update(pod)
			if err !=nil {
				klog.Errorf("Sync Pod %s Expose To %s Failed", pod.Name, deploydaemon.Spec.Expose)
			}
		}
	}

	if ! completeSync {

	}

	return err
}

func (c *Controller) remarkSuccessStatus(deploydaemon *v1alpha1.DeployDaemon, status bool, reason, message string) {

	deploydaemon.Status.Conditions.Status = status
	deploydaemon.Status.Conditions.Reason = reason
	deploydaemon.Status.Conditions.Message = message
}


func ( c *Controller) enqueueDeployDaemon(obj interface{}){
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}

	//check if the scheduler exist, if yes, add to AddDelayDefined(item interface{})
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		klog.Errorf("invalid resource key: %s", key)
	}

	deploydaemon, err:=c.deploydaemonLister.DeployDaemons(namespace).Get(name)
	if err != nil {
		klog.Errorf("get crd object %s failed", key)
	}

	if deploydaemon.Spec.Scheduler != "" {
		klog.Infof("get deploydaemon duration %s", deploydaemon.Spec.Scheduler )

		durationT,err := time.ParseDuration(deploydaemon.Spec.Scheduler)
		if err !=nil {
			klog.Errorf("Parse deploydaemon duration failed!")
		}
		klog.Infof("deploydaemon duration is %s ", durationT.String())
		c.workqueue.AddDelayDefined(key,durationT)

	}else {
		c.workqueue.AddRateLimited(key)
	}
}

func ( c *Controller ) updateDeployDaemonStatus(deploydaemon *v1alpha1.DeployDaemon) error {
	_, err := c.extclientset.DeploycontrolV1alpha1().DeployDaemons(deploydaemon.Namespace).Update(deploydaemon)
	return err
}

//func ( c *Controller ) updateDeployDaemonStatus(deploydaemon *v1alpha1.DeployDaemon, deployment *appsv1.Deployment) error {
//
//	klog.Infof("updateDeployDaemonStatus for deploydaemon %s: ", deploydaemon.Name)
//
//	if deployment !=nil {
//
//		deploydaemon.Status.Deployment = deployment.Status
//
//		// First check the current deployment replica / image / version -- if anyone of them different, that means the we need to delete current deployment and create a new one
//		// If any condition above matched, that means, we need to delete deployment and create a new one.
//		// Below change will need to delete deployment and create a new one
//		// TODO: Check if docker image changed, need to remove the deployment and create new one by call createDeployent method
//		if err := c.syncDeployment(deploydaemon, deployment); err != nil {
//			// createDeployent
//			return err
//		}
//
//		if deployment.Status.AvailableReplicas == deployment.Status.ReadyReplicas {
//			deploydaemon.Status.Conditions.Reason = "Waiting Pod Expose status sync"
//			deploydaemon.Status.Conditions.Message = "Deployment success!"
//		}else{
//
//			//updateDeployment -- it will call update deploydaemon status and update deployment
//			return fmt.Errorf("Waiting Pod Ready")
//		}
//
//	    // Second check expose service when deployment finish, if not match, we need to sync the expose status.
//	    if deploydaemon.Status.Conditions.Reason == "Waiting Pod Expose status sync"{
//			if err :=c.syncPodExposeStatus(deploydaemon); err !=nil {
//				deploydaemon.Status.Conditions.Message = fmt.Sprintf("Pod Sync Expose Failed: %s",err.Error())
//				klog.Errorf("Pod Sync For DeployDaemon %s faild",deploydaemon.Name)
//				// Here Need to Update DeployDamon Status
//				return err
//			}else {
//				deploydaemon.Status.Conditions.Reason = "All Status Synced"
//				deploydaemon.Status.Conditions.Status = true
//				deploydaemon.Status.Conditions.Message = "Deployment success!"
//			}
//		}
//	}
//
//	_, err :=c.extclientset.DeploycontrolV1alpha1().DeployDaemons(deploydaemon.Namespace).Update(deploydaemon)
//
//	return err
//}

//func (c *Controller ) updateDeployDaemonStatus(deploydaemon *v1alpha1.DeployDaemon, deployment *appsv1.Deployment) error{
//
//	 syncErr :=c.syncDeployDaemonStatus(deploydaemon, deployment)
//
//	 if syncErr != nil {
//		 deploydaemon.Status.Conditions.Message = syncErr.Error()
//
//		 if err !=nil {
//		 	klog.Errorf("Update DeployDaemon Status failed %s", err)
//		 }
//	 }
//
//	return syncErr
//}

func ( c *Controller ) createDeployent(deploydaemon *v1alpha1.DeployDaemon) (*appsv1.Deployment, error) {

	dp ,err := c.makeDeployment(deploydaemon)
	klog.Infof("deployment template is %s: ", dp)
	if err !=nil {
		klog.Errorf("Faile to create deployment: %s", err.Error())
	}

	klog.Infof("create deployment %s for deploydaemon %s: ", deploydaemon.GetDeploymentName()+"-"+deploydaemon.Spec.Version,deploydaemon.Name)

	deploydaemon.Status.Conditions = v1alpha1.ConditionsSpec{
		Type: "Successful",
		Status: false,
		Reason: "Waiting deployment ready",
		Message: "Waiting deployment ready",

	}

	return c.kubeclientset.AppsV1().Deployments(deploydaemon.ObjectMeta.Namespace).Create(dp)
}

func (c *Controller ) makeDeployment(deploydaemon *v1alpha1.DeployDaemon) (*appsv1.Deployment, error){
	klog.Infof("make deployment for deploydaemon: %s", deploydaemon.Name)

	// Here need to update deploydaemon status.cluster to store the deployment name

	//// Generate a short random hex string.
	//b, err := ioutil.ReadAll(io.LimitReader(randReader, 3))
	//if err != nil {
	//	return nil, err
	//}
	//gibberish := hex.EncodeToString(b)

	deploydaemon.Status = &v1alpha1.DeploydaemonStatus{
		Cluster: &v1alpha1.ClusterSpec{
                Name: deploydaemon.Name,
                NameSpace: deploydaemon.Namespace,
                DeploymentName: deploydaemon.GetDeploymentName()+"-"+deploydaemon.Spec.Version,
		},
	}

	klog.Infof("new deployment name is : %s", deploydaemon.Status.Cluster.DeploymentName)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
					// We execute the build's pod in the same namespace as where the build was
					// created so that it can access colocated resources.
					Namespace: deploydaemon.Namespace,
					// Generate a unique name based on the build's name.
					// Add a unique suffix to avoid confusion when a build
					// is deleted and re-created with the same name.
					// We don't use GenerateName here because k8s fakes don't support it.
					Name: deploydaemon.Status.Cluster.DeploymentName,
					// If our parent Build is deleted, then we should be as well.
					OwnerReferences: []metav1.OwnerReference{
						*metav1.NewControllerRef(deploydaemon, schema.GroupVersionKind{
							Group:   v1alpha1.SchemeGroupVersion.Group,
							Version: v1alpha1.SchemeGroupVersion.Version,
							Kind:    "DeployDaemon",
						}),
					},
					Labels: map[string]string{
						"app": deploydaemon.GetDeploymentName(),
						"version": deploydaemon.Spec.Version,
					},
				},
		Spec: appsv1.DeploymentSpec{
			Replicas: deploydaemon.Spec.Replica,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deploydaemon.GetDeploymentName(),
					"version": deploydaemon.Spec.Version,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deploydaemon.GetDeploymentName(),
						"version": deploydaemon.Spec.Version,
						"expose": deploydaemon.Spec.Expose,
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

func ( c *Controller ) isDone(status *v1alpha1.DeploydaemonStatus, err error) error {
	// judge is deploydaemon is done

	if status == nil || status.Conditions.Status == false || err !=nil {
		return fmt.Errorf("Sync is ongoing!")
	}
	return nil
}