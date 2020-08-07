package pkg

import (
	"net"
	"syscall"
)

type job struct {
	targetIP   net.IP
	targetPort int
	resolver   net.IP
}

func worker(jobsChan <-chan job) {
	handle := openSocket()
	defer syscall.Close(handle)
	for job := range jobsChan {
		packet := createPacket(job.targetIP, job.targetPort, job.resolver)
		sendPacket(handle, packet, job.resolver)
	}
}
