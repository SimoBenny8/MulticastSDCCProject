package main

import (
	"MulticastSDCCProject/pkg/endToEnd/client"
	"MulticastSDCCProject/pkg/endToEnd/server"
	"log"
	"sync"
)

func main() {

	//connetto il server
	nodeAddress := []string{"node1:8090"}
	addressGroups := []string{"node1:8090", "node2:8091", "node3:8092"} //"node1:8090", "node2:8091",

	go func() {
		err := server.RunServer(8090)
		if err != nil {
			log.Println("Error in connection")
			return
		}
	}()
	Connections := [3]*client.Client{}
	//creo il cluster di nodi (docker compose)
	// ogni connessione avrà la propria goroutine
	//il numero di nodi connessi e concorrenti è limitato dalla size del threadpool
	wg := sync.WaitGroup{}
	numNodes := 3
	index := 0
	for i := 0; i < numNodes; i++ {
		wg.Add(1)
		c := client.ConnectWithWaitGroup(addressGroups[i], &wg)
		wg.Wait()
		// deve tornare il riferimento alla connessione del client
		Connections[i] = c
		index = index + 1
	}

	//creazione della pool
	//pool, err := multicast.CreatePool("1", 3, Connections)
	var err error
	//invio del messaggio per ogni goroutine in PARALLELO
	log.Println("Multicaster selected : ", nodeAddress)
	multicast := func(metadata map[string]string, payload []byte) {
		wg = sync.WaitGroup{}
		for j := 0; j < len(Connections)-1; j++ {
			wg.Add(1)
			go func() {
				var localErr error
				localErr = Connections[j].Send(metadata, payload, nil)

				if localErr != nil {
					err = localErr
				}

				wg.Done()
			}()

		}
		wg.Wait()

		if err != nil {
			return
		}
	}
	multicast(make(map[string]string), []byte("PROVA"))
	//sendToOne := func(client *grpc.Client, metadata map[string]string, payload []byte, wg *sync.WaitGroup) {
	//	log.Println("Sending from ..", client.Connection.Target())
	//	err := client.Send(make(map[string]string), []byte(""))
	//	if err != nil {
	//		return
	//	}
	//	wg.Done()
	//	wg.Wait()
	//}
	//sendToAll := func(client *grpc.Client, metadata map[string]string, payload []byte, poolId string) {
	//	wg := sync.WaitGroup{}
	//	for i := 0; i < numNodes; i++ {
	//		wg.Add(1)
	//		go sendToOne(Connections[i], metadata, payload, &wg)
	//		wg.Wait()
	//	}
	//}
	//
	//multicast := func(metadata map[string]string, payload []byte, poolId string, client *grpc.Client) {
	//	sendToAll(client, make(map[string]string), []byte(""), poolId)
	//}
	//
	//multicast(make(map[string]string), []byte("PROVA"), "1", Connections[0])

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
