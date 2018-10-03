package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

/**
* Init is called right after main each time
 */
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	fmt.Println("threadtest.main() : starting exectuion")
	fmt.Println()
	start := time.Now()

	// simpletest()
	// mediumtest()
	// raceGoogleVersusBing()
	// timeoutTest()
	// waitGroupTest()
	showFirst3Threads()

	elapsed := time.Since(start)
	fmt.Println("elapsed time to execute: ", elapsed)
	fmt.Println()
	fmt.Println("threadtest.main() : ending exectuion")

}

func simpletest() {
	a := 1
	b := 2

	operationDone := make(chan bool)
	go func() {
		b = a * b

		operationDone <- true
	}()

	<-operationDone

	a = b * b

	fmt.Println("Hit Enter when you want to see the answer")
	fmt.Scanln()

	fmt.Printf("a = %d, b = %d\n", a, b)
}

func mediumtest() {
	defer timeTrack(time.Now(), funcName())

	query := "Our Query"
	respond := make(chan string)

	go googleItFake(respond, query)

	queryResp := <-respond

	fmt.Printf("Sent query:\t\t %s\n", query)
	fmt.Printf("Got Response:\t\t %s\n", queryResp)
}

func timeTrack(pStart time.Time, pName string) {
	elapsed := time.Since(pStart)
	fmt.Printf("%s took %s\n", pName, elapsed)
}

func funcName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func raceGoogleVersusBing() {

	query := "Our Query"
	respond := make(chan string, 2)

	go googleItFake(respond, query)
	go bingItFake(respond, query)

	queryResp := <-respond

	fmt.Printf("Sent query:\t\t %s\n", query)
	fmt.Printf("Got Response:\t\t %s\n", queryResp)
}

func googleItFake(respond chan<- string, query string) {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)

	respond <- "A Google Response"
}

func bingItFake(respond chan<- string, query string) {
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)

	respond <- "A Bing Response"
}

func timeoutTest() {

	query := "Our Query"
	respond := make(chan string, 1)

	go googleItFake(respond, query)

	select {
	case queryResp := <-respond:
		fmt.Printf("Sent query:\t\t %s\n", query)
		fmt.Printf("Got Response:\t\t %s\n", queryResp)

	case <-time.After(5 * time.Second):
		fmt.Printf("A timeout occurred for query:\t\t %s\n", query)
	}
}

func waitGroupTest() {

	respond := make(chan string, 5)
	var wg sync.WaitGroup

	wg.Add(5)
	go checkDNS(respond, &wg, "pragmacoders.com", "ns1.nameserver.com")
	go checkDNS(respond, &wg, "pragmacoders.com", "ns2.nameserver.com")
	go checkDNS(respond, &wg, "pragmacoders.com", "ns3.nameserver.com")
	go checkDNS(respond, &wg, "pragmacoders.com", "ns4.nameserver.com")
	go checkDNS(respond, &wg, "pragmacoders.com", "ns5.nameserver.com")

	wg.Wait()
	close(respond)

	for queryResp := range respond {
		fmt.Printf("Got Response:\t %s\n", queryResp)
	}
}

func checkDNS(respond chan<- string, wg *sync.WaitGroup, query string, ns string) {
	defer wg.Done()

	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	respond <- fmt.Sprintf("%s responded to query: %s", ns, query)
}

func showFirst3Threads() {
	rand.Seed(time.Now().UTC().UnixNano())

	respond := make(chan string, 5)

	go checkDNSNoWG(respond, "pragmacoders.com", "ns1.nameserver.com")
	go checkDNSNoWG(respond, "pragmacoders.com", "ns2.nameserver.com")
	go checkDNSNoWG(respond, "pragmacoders.com", "ns3.nameserver.com")
	go checkDNSNoWG(respond, "pragmacoders.com", "ns4.nameserver.com")
	go checkDNSNoWG(respond, "pragmacoders.com", "ns5.nameserver.com")

	for i := 1; i <= 3; i++ {
		queryResp := <-respond
		fmt.Printf("Got Response:\t %s\n", queryResp)
	}
}

func checkDNSNoWG(respond chan<- string, query string, ns string) {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	respond <- fmt.Sprintf("%s responded to query: %s", ns, query)
}
