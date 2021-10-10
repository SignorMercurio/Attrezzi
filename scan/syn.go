package scan

import (
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/routing"
	"github.com/mostlygeek/arp"
	"github.com/phayes/freeport"
	"github.com/pkg/errors"
)

type SYNScanner struct {
	scan          chan hostScan
	targets       []net.IP
	timeout       time.Duration
	routines      int
	serialOptions gopacket.SerializeOptions
}

func NewSYNScanner(targets []net.IP, timeout time.Duration, routines int) *SYNScanner {
	return &SYNScanner{
		scan:     make(chan hostScan, routines),
		targets:  targets,
		timeout:  timeout,
		routines: routines,
		serialOptions: gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		},
	}
}

func (ss *SYNScanner) Start() {
	for i := 0; i < ss.routines; i++ {
		go func() {
			for {
				scan := <-ss.scan
				if scan.ports == nil || len(scan.ports) == 0 {
					break
				}

				if result, err := ss.hostScan(scan); err == nil {
					scan.result <- &result
				}
				close(scan.done)
			}
		}()
	}
}

func (ss *SYNScanner) Scan(ctx context.Context, ports []int) ([]Result, error) {
	var (
		wg         = &sync.WaitGroup{}
		resultChan = make(chan *Result)
		results    = []Result{}
		done       = make(chan struct{})
	)

	go func() {
		for {
			result := <-resultChan
			if result == nil {
				close(done)
				break
			}
			results = append(results, *result)
		}
	}()

	for _, ip := range ss.targets {
		wg.Add(1)
		go func(ip net.IP, ports []int, wg *sync.WaitGroup) {
			done := make(chan struct{})
			ss.scan <- hostScan{
				result: resultChan,
				ip:     ip,
				ports:  ports,
				done:   done,
				ctx:    ctx,
			}
			<-done
			wg.Done()
		}(ip, ports, wg)
	}
	wg.Wait()
	close(resultChan)
	close(ss.scan)
	<-done

	return results, nil
}

func (ss *SYNScanner) hostScan(scan hostScan) (Result, error) {
	var (
		result       = NewResult(scan.ip)
		openChan     = make(chan int)
		filteredChan = make(chan int)
		closedChan   = make(chan int)
		done         = make(chan struct{})
	)

	select {
	case <-scan.ctx.Done():
		return result, nil
	default:
	}

	router, err := routing.New()
	if err != nil {
		return result, err
	}

	iface, gateway, src, err := router.Route(scan.ip)
	if err != nil {
		return result, err
	}

	handle, err := pcap.OpenLive(iface.Name, 65535, true, pcap.BlockForever)
	if err != nil {
		return result, err
	}
	defer handle.Close()

	start := time.Now()
	go func() {
		for {
			select {
			case open := <-openChan:
				if open == 0 {
					close(done)
					return
				}
				if result.Latency <= 0 {
					result.Latency = time.Since(start)
				}
				for _, existing := range result.Open {
					if existing == open {
						continue
					}
				}
				result.Open = append(result.Open, open)
			case filtered := <-filteredChan:
				if result.Latency <= 0 {
					result.Latency = time.Since(start)
				}
				for _, existing := range result.Filtered {
					if existing == filtered {
						continue
					}
				}
				result.Filtered = append(result.Filtered, filtered)
			case closed := <-closedChan:
				if result.Latency <= 0 {
					result.Latency = time.Since(start)
				}
				for _, existing := range result.Closed {
					if existing == closed {
						continue
					}
				}
				result.Closed = append(result.Closed, closed)
			}
		}
	}()

	fPort, err := freeport.GetFreePort()
	if err != nil {
		return result, err
	}

	MACaddr, err := ss.getMACaddr(scan.ip, gateway, src, iface)
	if err != nil {
		return result, err
	}

	eth := layers.Ethernet{
		SrcMAC:       iface.HardwareAddr,
		DstMAC:       MACaddr,
		EthernetType: layers.EthernetTypeIPv4,
	}
	ip4 := layers.IPv4{
		SrcIP:    src,
		DstIP:    scan.ip,
		Version:  4,
		TTL:      255,
		Protocol: layers.IPProtocolTCP,
	}
	tcp := layers.TCP{
		SrcPort: layers.TCPPort(fPort),
		DstPort: 0,
		SYN:     true,
	}
	tcp.SetNetworkLayerForChecksum(&ip4)

	listen := make(chan struct{})
	flow := gopacket.NewFlow(layers.EndpointIPv4, scan.ip, src)
	go func() {
		eth := &layers.Ethernet{}
		ip4 := &layers.IPv4{}
		tcp := &layers.TCP{}

		parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, eth, ip4, tcp)

		for {
			data, _, err := handle.ReadPacketData()
			if err == pcap.NextErrorTimeoutExpired {
				break
			} else if err == io.EOF {
				break
			} else if err != nil {
				continue
			}

			decoded := []gopacket.LayerType{}
			if err := parser.DecodeLayers(data, &decoded); err != nil {
				continue
			}
			for _, layerType := range decoded {
				switch layerType {
				case layers.LayerTypeIPv4:
					if ip4.NetworkFlow() != flow {
						continue
					}
				case layers.LayerTypeTCP:
					if tcp.DstPort != layers.TCPPort(fPort) {
						continue
					} else if tcp.SYN && tcp.ACK {
						openChan <- int(tcp.SrcPort)
					} else if tcp.RST {
						closedChan <- int(tcp.SrcPort)
					}
				}
			}
		}
		close(listen)
	}()

	for _, port := range scan.ports {
		tcp.DstPort = layers.TCPPort(port)
		ss.send(handle, &eth, &ip4, &tcp)
	}

	timer := time.AfterFunc(ss.timeout, func() {
		handle.Close()
	})
	defer timer.Stop()

	<-listen
	close(openChan)
	<-done

	return result, nil
}

func (ss *SYNScanner) getMACaddr(ip net.IP, gateway net.IP, src net.IP, iface *net.Interface) (net.HardwareAddr, error) {
	mac := arp.Search(ip.String())
	if mac != "00:00:00:00:00:00" {
		if macAddr, err := net.ParseMAC(mac); err == nil {
			return macAddr, nil
		}
	}

	arpDst := ip
	if gateway != nil {
		arpDst = gateway
	}

	handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		return nil, errors.Wrap(err, "open device")
	}
	defer handle.Close()

	start := time.Now()

	eth := layers.Ethernet{
		SrcMAC:       iface.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(iface.HardwareAddr),
		SourceProtAddress: []byte(src),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
		DstProtAddress:    []byte(arpDst),
	}

	if err = ss.send(handle, &eth, &arp); err != nil {
		return nil, err
	}

	for {
		if time.Since(start) > ss.timeout {
			return nil, errors.New("get ARP reply(timeout)")
		}
		data, _, err := handle.ReadPacketData()
		if err == pcap.NextErrorTimeoutExpired {
			continue
		} else if err != nil {
			return nil, err
		}

		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)
		if arpLayer := packet.Layer(layers.LayerTypeARP); arpLayer != nil {
			arp := arpLayer.(*layers.ARP)
			if net.IP(arp.SourceProtAddress).Equal(arpDst) {
				return net.HardwareAddr(arp.SourceHwAddress), nil
			}
		}
	}
}

func (ss *SYNScanner) send(handle *pcap.Handle, l ...gopacket.SerializableLayer) error {
	buf := gopacket.NewSerializeBuffer()
	if err := gopacket.SerializeLayers(buf, ss.serialOptions, l...); err != nil {
		return errors.Wrap(err, "serialize layers")
	}
	if err := handle.WritePacketData(buf.Bytes()); err != nil {
		return errors.Wrap(err, "write packets")
	}
	return nil
}
