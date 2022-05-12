package main

// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-api-machinery/controllers.md

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

func newKubeClient() *kubernetes.Clientset {
	// Load kubeconfig
	// https://github.com/kubernetes-sigs/controller-runtime/blob/13f1400cd4fc0f15b6877453c50c95d35d146e47/pkg/client/config/config.go
	loader := clientcmd.NewDefaultClientConfigLoadingRules()
	loader.Precedence = append(loader.Precedence, path.Join(os.Getenv("HOME"), ".kube", "config"))
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loader, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		panic("Failed to load kubeconfig: " + err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic("Failed to create K8S client: " + err.Error())
	}

	return client
}

type podAddedWorkItem struct {
	key string
	pod *v1.Pod
}

type podUpdatedWorkItem struct {
	key string
	old *v1.Pod
	new *v1.Pod
}

type podDeletedWorkItem struct {
	key string
	pod *v1.Pod
}

func processWorkItem(item interface{}) error {
	if added, ok := item.(podAddedWorkItem); ok {
		fmt.Printf("Pod added: %s\n", added.key)
		return nil
	}
	if updated, ok := item.(podUpdatedWorkItem); ok {
		fmt.Printf("Pod updated: %s\n", updated.key)
		return nil
	}
	if deleted, ok := item.(podDeletedWorkItem); ok {
		fmt.Printf("Pod deleted: %s\n", deleted.key)
		return nil
	}
	return fmt.Errorf("Unknown work item: %v", item)
}

func processWorkQueue(wq workqueue.RateLimitingInterface) bool {
	// Pull the next work item from queue.
	// It should be a key we use to lookup something in a cache.
	// But we have here a workitem itself, is it acceptable or not?
	item, shutdown := wq.Get()
	if shutdown {
		fmt.Println("Workqueue is shutdown")
		return false
	}

	// You always have to indicate to the queue that you've completed a piece of work
	defer wq.Done(item)

	//err := processWorkItem(item.(myWorkItem))
	err := processWorkItem(item)
	if err == nil {
		// When all is good, tell the queue to stop tracking history for your key.
		// This will reset things like failure counts for per-item rate limiting
		wq.Forget(item)
		return true
	}

	// This method allows for pluggable error handling
	fmt.Printf("Unable to process %v : %v", item, err)
	runtime.HandleError(fmt.Errorf("Failed with %v : %v", item, err))

	// Requeue the item to work on later (the method adds backoff)
	wq.AddRateLimited(item)
	return true
}

func main() {
	client := newKubeClient()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	fmt.Println("Initializing workqueue...")
	wq := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	defer wq.ShutDown() // when queue is shutdown it will trigger workers to end

	stopper := make(chan struct{})

	// SharedInformers provide hooks to receive notifications of adds, updates, and deletes for a particular resource.
	informerFactory := informers.NewSharedInformerFactory(client, 0)

	fmt.Println("Adding event handlers...")
	podInformer := informerFactory.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			// Add resource handler
			AddFunc: func(obj interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(obj)
				if err != nil {
					fmt.Printf("Failed to get key in AddFunc: %v\n", err)
					return
				}
				// TODO: is it ok to store a pointer or we always have to search cache by key?
				wq.Add(podAddedWorkItem{key: key, pod: obj.(*v1.Pod)})
			},

			// Update resource handler
			UpdateFunc: func(old interface{}, new interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(new)
				if err != nil {
					fmt.Printf("Failed to get key in UpdateFunc: %v\n", err)
					return
				}
				// TODO: is it ok to store a pointer or we always have to search cache by key?
				wq.Add(podUpdatedWorkItem{key: key, old: old.(*v1.Pod), new: new.(*v1.Pod)})
			},

			// Delete resource handler
			DeleteFunc: func(obj interface{}) {
				key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
				if err != nil {
					fmt.Printf("Failed to get key in DeleteFunc: %v\n", err)
					return
				}
				// TODO: is it ok to store a pointer or we always have to search cache by key?
				wq.Add(podDeletedWorkItem{key: key, pod: obj.(*v1.Pod)})
			},
		},
	)

	go informerFactory.Start(stopper)
	fmt.Println("Informer factory started")

	go podInformer.Run(stopper)
	fmt.Println("Pod informer started")

	// Wait for secondary caches to fill before starting your work.
	// Secondary means the resource we will watch (created/deleted).
	// Primary means the resources we will updated status for.
	fmt.Println("Waiting for secondary caches...")
	if !cache.WaitForCacheSync(stopper, podInformer.HasSynced) {
		fmt.Println("Exiting after WaitForCacheSync")
		return
	}

	// Start controller
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		fmt.Println("Controller started")
		for {
			for processWorkQueue(wq) {
				fmt.Println("Work item processed")
			}

			select {
			case <-stopper:
				fmt.Println("Controller stopped")
				wg.Done()
				return

			case <-time.After(1 * time.Second):
				continue
			}
		}
	}()

	// Wait for term signal
	<-ctx.Done()

	// Stop the stuff
	close(stopper)
	wq.ShutDown()

	// Wait for stopping
	wg.Wait()
	fmt.Println("All stopped")
}
