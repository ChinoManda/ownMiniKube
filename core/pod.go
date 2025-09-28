package core

import (
	"context"
	"fmt"
	"log"
	"syscall"
	"time"
        containerd "github.com/containerd/containerd/v2/client"
	"github.com/containerd/containerd/v2/pkg/oci"
	"github.com/containerd/containerd/v2/pkg/cio"
	"github.com/containerd/containerd/namespaces"
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
func InitClient()(*containerd.Client) {
    var err error
		Client, err := containerd.New("/run/containerd/containerd.sock")
    if err != nil {
        log.Fatalf("No se pudo crear el client: %v", err)
    }
		return Client
}
func InitNamespace()(context.Context)  {	
	ctx := namespaces.WithNamespace(context.Background(), "ownkube")
	return ctx
}

func GetPodByID(ctx context.Context, client *containerd.Client, id string) (*Pod, error) {
	c, err := client.LoadContainer(ctx, id)
		if err != nil {
		return nil, err
	}
	var (
		task    containerd.Task
		statusC <-chan containerd.ExitStatus
	)
	t, err := c.Task(ctx, cio.Load)
	if err == nil {
		task = t
		statusC, _ = t.Wait(ctx)
	} else {
		task = nil
		statusC = nil
	}
	return &Pod{
		Id:        id,
		client:    client,
		ctx:       ctx,
		container: c,
		task:      &task,    
		statusCh:  statusC,  
	}, nil
}

func PullImage(client *containerd.Client, ctx context.Context, registryImage string) (containerd.Image, error) {
    image, err := client.Pull(ctx, registryImage, containerd.WithPullUnpack)
    if err != nil {
        return nil, err
    }
    log.Printf("Successfully pulled image %s\n", image.Name())
    return image, nil
}

func ListPods(client *containerd.Client, ctx context.Context) []containerd.Container {
	containers, _:= client.Containers(ctx)
	return containers
}

func ListRunningPods(containers []containerd.Container, ctx context.Context) ([]containerd.Container)  {
	var running []containerd.Container
			for _,c := range containers {
			task, err:= c.Task(ctx, nil)
			if err != nil {
				continue
			}
			stat, _ := task.Status(ctx)
      if (stat.Status == "running"){
				running = append(running, c)
			}
		}
		return running
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

	fmt.Println("task started")
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
	fmt.Println("sigkill killed")
	if err := (*pod.task).Kill(pod.ctx, syscall.SIGKILL); err != nil {
		return 0, fmt.Errorf("failed to SIGKILL task: %w", err)
	}
	status := <-pod.statusCh
	code, _, err := status.Result()
	return code, err
	}
}
func (pod *Pod) Delete() error {
	status, err := (*pod.task).Status(pod.ctx)
	fmt.Println(status)
	if err != nil {
		return err
	}
	if status.Status != "stopped" {
		return fmt.Errorf("task still running, cannot delete container")
	}
  if _, err := (*pod.task).Delete(pod.ctx); err != nil {
    return fmt.Errorf("failed to delete task: %w", err)
  }

	return pod.container.Delete(pod.ctx)
}
