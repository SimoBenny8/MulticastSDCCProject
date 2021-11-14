package main

import (
	server2 "MulticastSDCCProject/pkg/ServiceRegistry/server"
	"MulticastSDCCProject/pkg/endToEnd/server"
	"MulticastSDCCProject/pkg/restApi"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

func main() {

	delay := flag.Uint("DELAY", uint(utils.GetEnvIntWithDefault("DELAY", 0)), "delay for sending operations (ms)")
	grpcPort := flag.Uint("GRPC_PORT", uint(utils.GetEnvIntWithDefault("GRPC_PORT", 90)), "port number of the grpc server")
	restPort := flag.Uint("REST_PORT", uint(utils.GetEnvIntWithDefault("REST_PORT", 80)), "port number of the rest server")
	restPath := flag.String("restPath", utils.GetEnvStringWithDefault("REST_PATH", "/multicast/v1"), "path of the rest api")
	numThreads := flag.Uint("NUM_THREADS", uint(utils.GetEnvIntWithDefault("NUM_THREADS", 1)), "number of threads used to multicast messages")
	verb := flag.Bool("VERBOSE", utils.GetEnvBoolWithDefault("VERBOSE", true), "Turn verbose mode on or off.")
	registry_addr := flag.String("REGISTRY_ADDR", "registry:90", "service registry address")
	r := flag.Bool("REGISTRY", utils.GetEnvBoolWithDefault("REGISTRY", false), "start multicast registry")
	application := flag.Bool("APP", utils.GetEnvBoolWithDefault("APP", false), "start multicast application")
	myId := flag.Int("ID", utils.GetEnvIntWithDefault("ID", 0), "number id of member")

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
	log.Println("start")
	var wg sync.WaitGroup
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
			err := restApi.Run(*grpcPort, *restPort, *registry_addr, *restPath, int(*numThreads), *delay, *verb)
			if err != nil {
				log.Println("Error in running applicatioon", err.Error())
				return
			}
			wg.Done()
		}()

	}

	log.Println("App started")

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
