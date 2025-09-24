package main

import (
	"context"
	"fmt"
	"log"
	"syscall"
	"time"
        containerd "github.com/containerd/containerd/v2/client"
	"github.com/containerd/containerd/v2/pkg/oci"
	"github.com/containerd/containerd/v2/pkg/cio"
	"github.com/containerd/containerd/v2/pkg/namespaces"
	"github.com/google/uuid"
)

type Pod struct {
	Id       string
	client   *containerd.Client
	ctx      context.Context
	container containerd.Container
	task     *containerd.Task
  statusCh  <-chan containerd.ExitStatus
}

func PullImage(client *containerd.Client, ctx context.Context, registryImage string) (containerd.Image, error) {
    image, err := client.Pull(ctx, registryImage, containerd.WithPullUnpack)
    if err != nil {
        return nil, err
    }
    log.Printf("Successfully pulled image %s\n", image.Name())
    return image, nil
}

func NewPod(client *containerd.Client, ctx context.Context, image containerd.Image, name string) (*Pod, error) {
	id := generateNewID(name)

	container, err := client.NewContainer(
		ctx,
		id,
		containerd.WithImage(image),
		containerd.WithNewSnapshot(id+"-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)
	if err != nil {
		return nil, err
	}

	return &Pod{
		Id:        id,
		client:    client,
		ctx:       ctx,
		container: container,
	}, nil
}

func generateNewID(name string) string {
	id := uuid.New()
	return fmt.Sprintf("%s-%s", name, id)
}

func (pod *Pod) Run() error {
	task, err := pod.container.NewTask(pod.ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return err
	}
	pod.task = &task

	statusCh, err := task.Wait(pod.ctx)
	if err != nil {
		return err
	}
	pod.statusCh = statusCh

	if err := task.Start(pod.ctx); err != nil {
		return err
	}

	log.Println("task started")
	return nil
}

func (pod *Pod) Kill() (uint32, error) {
	err := (*pod.task).Kill(pod.ctx, syscall.SIGTERM)
	if err != nil {
		fmt.Println("Failed sigterm")
		return 0, err
	}
	fmt.Println("Killing")

	select{
  case status := <-pod.statusCh:
  code, _, err := status.Result()
	if err != nil {
		return 0, err
		fmt.Println("Failed")
	}
	fmt.Println("Killed")
	return code, nil

	case <-time.After(10 * time.Second):
	if err := (*pod.task).Kill(pod.ctx, syscall.SIGKILL); err != nil {
		return 0, fmt.Errorf("failed to SIGKILL task: %w", err)
	}
	status := <-pod.statusCh
	code, _, err := status.Result()
	return code, err
	}
}

func main()  {
	cli, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx := namespaces.WithNamespace(context.Background(), "own-kube")
  image, err := PullImage(cli, ctx, "docker.io/chinomandarin/demo:pythonapp.1.0")
  if err != nil {
    log.Fatal(err)
  }
 	AllPods := []*Pod{}
 	for i := 0; i < 5; i++ {
    pod, err := NewPod(cli, ctx, image, fmt.Sprintf("PythonPod-%d", i))
    if err != nil {
        log.Fatal(err)
    }
    AllPods = append(AllPods, pod)
	}
	for _, pod := range AllPods {
    task, err := pod.container.NewTask(pod.ctx, cio.NewCreator(cio.WithStdio))
    if err != nil {
        log.Fatalf("failed to create task: %v", err)
    }
    pod.task = &task
    exitStatusC, err := task.Wait(pod.ctx)
    if err != nil {
        log.Fatalf("failed to wait for task: %v", err)
    }
    pod.statusCh = exitStatusC
    if err := task.Start(pod.ctx); err != nil {
        log.Fatalf("failed to start task: %v", err)
    }

    log.Printf("Pod %s started", pod.Id)
	}
}





