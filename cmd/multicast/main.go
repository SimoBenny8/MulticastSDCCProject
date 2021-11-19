package main

import (
	"flag"
	"fmt"
	server2 "github.com/SimoBenny8/MulticastSDCCProject/pkg/ServiceRegistry/server"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/endToEnd/server"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/restApi"
	"github.com/SimoBenny8/MulticastSDCCProject/pkg/util"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

func main() {

	delay := flag.Uint("delay", uint(util.GetEnvIntWithDefault("DELAY", 0)), "delay for sending operations (ms)")
	grpcPort := flag.Uint("rpcPort", uint(util.GetEnvIntWithDefault("GRPC_PORT", 90)), "port number of the grpc server")
	restPort := flag.Uint("restPort", uint(util.GetEnvIntWithDefault("REST_PORT", 80)), "port number of the rest server")
	restPath := flag.String("restPath", util.GetEnvStringWithDefault("REST_PATH", "/multicast/v1"), "path of the rest api")
	numThreads := flag.Uint("nThreads", uint(util.GetEnvIntWithDefault("NUM_THREADS", 1)), "number of threads used to multicast messages")
	verb := flag.Bool("verbose", util.GetEnvBoolWithDefault("VERBOSE", true), "Turn verbose mode on or off.")
	registryAddr := flag.String("registryAddr", "registry:90", "service registry address")
	r := flag.Bool("registry", util.GetEnvBoolWithDefault("REGISTRY", false), "start multicast registry")
	application := flag.Bool("application", util.GetEnvBoolWithDefault("APP", false), "start multicast application")
	//myId := flag.Int("id", util.GetEnvIntWithDefault("ID", 0), "number id of member")

	//utils.Myid = *myId
	flag.Parse()
	services := make([]func(registrar grpc.ServiceRegistrar) error, 0)

	if *application {
		log.Println("Adding basic communication service to gRPC server")
		services = append(services, server.Register)
	}
	if *r {
		log.Println("Adding multicast registry service to gRPC server")
		services = append(services, server2.Registration)
	}
	log.Println("start services")
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
		if err != nil {
			return
		}
		log.Printf("Grpc-Server started at %v", lis.Addr().String())

		s := grpc.NewServer()
		for _, grpcService := range services {
			err = grpcService(s)
			if err != nil {
				return
			}

		}

		if err = s.Serve(lis); err != nil {
			return
		}
		wg.Done()
	}()

	if *application {

		wg.Add(1)
		go func() {
			err := restApi.Run(*grpcPort, *restPort, *registryAddr, *restPath, int(*numThreads), *delay, *verb)
			if err != nil {
				log.Println("Error in running application", err.Error())
				return
			}
			wg.Done()
		}()
	}

	wgChan := make(chan bool)

	go func() {
		wg.Wait()
		wgChan <- true
	}()

	log.Println("App started")
	select {
	case <-wgChan:
	}

	//seleziono il nodo che effettuerà il multicasting
	//salvo l'interfaccia della procedura grpc
	//utilizzo l'interfaccia per mandare i messaggi

	//numNodes := 3
	//per avere un buon grado di reliability posso utilizzare waitgroup che mi permette di attendere che il messaggio
	//sia consegnato a livello applicativo
	//posso utilizzare i canali, un canale fittizio true e uno di errore per sapere quando si riceve un errore che tipo
	//di errore si è ricevuto (o altri metodi)
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//for i := 0; i < numNodes; i++ {
	//	go func() {
	//	err :=   funzione BMulticast(..)
	//response chann : posso aspettarmi il  canale di risposta da ricordare che è un semaforo e deve consumarsi
	//deve essere letto
	//	}()
	//reliability : attendo l'ack
	//wg.Done() solo quando ho finito di attendere la chiamata

}
