package lxo

import (
	"fmt"

	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
)

// StopContainer will try to stop the container with provided name.
// It will retry for half a minute and return success when it's stopped.
// It will also return success when the container does not exist.
func StopContainer(server lxd.ContainerServer, id string, timeout, retries int) error {
	ETag := ""
	var lastErr error
	for i := 1; i <= retries; i++ {
		lxdReq := api.ContainerStatePut{
			Action:  "stop",
			Timeout: timeout,
			Force:   i == retries,
		}
		op, err := server.UpdateContainerState(id, lxdReq, ETag)
		if err != nil {
			if err.Error() == "not found" { // it's not around, that's ok with us
				return nil
			}
			return fmt.Errorf("failed to stop container %v, %v", id, err)
		}

		err = op.Wait()
		if err != nil && err.Error() == "The container is already stopped" {
			return nil
		}
		lastErr = err
		// we try again with or without err
	}
	return lastErr
}

// StartContainer will start the container and wait till operation is done or
// return an error
func StartContainer(server lxd.ContainerServer, id string) error {
	ETag := ""
	lxdReq := api.ContainerStatePut{
		Action:  "start",
		Timeout: -1,
	}
	op, err := server.UpdateContainerState(id, lxdReq, ETag)
	if err != nil {
		return err
	}

	return op.Wait()
}

// CreateContainer will create the container and wait till operation is done or
// return an error
func CreateContainer(server lxd.ContainerServer, container api.ContainersPost) error {
	op, err := server.CreateContainer(container)
	if err != nil {
		return err
	}
	return op.Wait()
}

// UpdateContainer will create the container and wait till operation is done or
// return an error
func UpdateContainer(server lxd.ContainerServer, id string, container api.ContainerPut, ETag string) error {
	op, err := server.UpdateContainer(id, container, ETag)
	if err != nil {
		return err
	}
	return op.Wait()
}

// DeleteContainer will delete the container and wait till operation is done or
// return an error
func DeleteContainer(server lxd.ContainerServer, id string) error {
	op, err := server.DeleteContainer(id)
	if err != nil {
		return err
	}
	return op.Wait()
}

// ExecContainer runs a command on a container and wait till operation is done or
// return an error
func ExecContainer(server lxd.ContainerServer, id string, containerExec api.ContainerExecPost, execArgs *lxd.ContainerExecArgs) (lxd.Operation, error) {
	op, err := server.ExecContainer(id, containerExec, execArgs)
	if err != nil {
		return op, err
	}
	err = op.Wait()
	return op, err
}
