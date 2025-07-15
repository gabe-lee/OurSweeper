package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"

	App "github.com/gabe-lee/OurSweeper/internal"
	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {
	sWorld := App.ServerWorld{}
	sWorld.InitNew(1)
	clientToServer := make(chan []byte, 4)
	serverToClient := make(chan []byte, 64)
	closeServer := make(chan struct{})
	go runTempServer(&sWorld, serverToClient, clientToServer, closeServer)
	client := App.GameClient{
		ErrorWriter:     os.Stderr, //FIXME
		SendMessages:    clientToServer,
		RecieveMessages: serverToClient,
		DebugServer:     &sWorld,
	}
	client.Init(&sWorld, clientToServer, serverToClient)
	ebiten.SetWindowSize(App.WINDOW_WIDTH, App.WINDOW_HEIGHT)
	ebiten.SetWindowTitle("OurSweeper")
	if err := ebiten.RunGame(&client); err != nil {
		log.Fatal(err)
	}
	closeServer <- struct{}{} //FIXME
	_, ok := <-closeServer
	if !ok {
		fmt.Printf("closed server")
	}
}

// FIXME
func runTempServer(world *App.ServerWorld, serverToClient chan<- []byte, clientToServer <-chan []byte, closeServer chan struct{}) {
	keepAlive := true
	for keepAlive {
		select {
		case <-closeServer:
			keepAlive = false
		case msg := <-clientToServer:
			var code uint32
			r := bytes.NewReader(msg)
			err := binary.Read(r, App.BYTE_ORDER, &code)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to read msg code: %s", err)
			} else {
				switch code {
				case App.CLIENT_SWEEP:
					var request App.SweepRequest
					err = request.Deserialize(r, App.BYTE_ORDER)
					if err != nil {
						fmt.Fprintf(os.Stderr, "unable to deserialize SweepRequest: %s", err)
					} else {
						result := world.SweepTile(request.Pos.ToCoord())

						if result.Len > 0 {

							response, err := result.CodedSerialize(App.BYTE_ORDER)
							if err != nil {
								fmt.Fprintf(os.Stderr, "unable to serialize SweepResult: %s", err)
							} else {
								serverToClient <- response
							}
						}
					}
				}
			}
		default:
		}
	}
	close(closeServer)
}
