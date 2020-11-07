package main

import (
	"dnscop/dnsmsg"
	"flag"
	"log"
	"net"
	"regexp"
	"strings"
)

const maxDnsPacketSize = 512
const VERSION = "0.0.1"

func main() {
	listen := flag.String("listen", ":53", "listen IP:Port. ex. 127.0.0.1:53")
	resolver := flag.String("resolver", "8.8.8.8:53", "Resolver IP:Port. ex. 8.8.8.8:53")
	flag.Parse()

	log.Printf("==== DNSCOP (ver %s) ====\n", VERSION)
	log.Println("Listen: " + *listen)
	log.Println("Resolver: " + *resolver)

	packet, err := net.ListenPacket("udp", *listen)
	if err != nil {
		log.Fatal(err)
	}
	defer packet.Close()

	for {
		buf := make([]byte, maxDnsPacketSize)
		readbyte, address, err := packet.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}
		go handleDnsRequest(packet, address, buf[:readbyte])
	}
}

func handleDnsRequest(packet net.PacketConn, address net.Addr, data []byte) {
	//log.Println(data)
	name, err := dnsmsg.GetQuestionName(data)
	name = strings.TrimRight(name, ".")
	log.Println(name)

	r,_ := regexp.Compile("www.youtube.com|youtube.com|i.ytimg.com|.+.googlevideo.com")
	if r.MatchString(name) {
		log.Println("  ** block youtube **")
		return
	}

	response, err := dnsmsg.Send("8.8.8.8:53", data)
	if err != nil {
		log.Fatal(err)
	}
	packet.WriteTo(response, address)
}
