// +build !darwin

package udockerutilstils

import "os/exec"

type Docker struct {
}

//获取docker容器ID号
func (r *Docker) GetDockerContainerId() (res string, err error) {
	//cat /proc/self/cgroup | head -1 /proc/self/cgroup|cut -d/ -f3|cut -c1-12

	command := `cat /proc/self/cgroup | head -1 /proc/self/cgroup|cut -d/ -f3|cut -c1-12`
	cmd := exec.Command("/bin/bash", "-c", command)
	bytes, err := cmd.Output()
	if err != nil {
		return
	}
	res = string(bytes)
	return
}
