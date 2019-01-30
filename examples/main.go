package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/nomasters/pinata"
)

func main() {
	key := os.Getenv("PINATA_API_KEY")
	secret := os.Getenv("PINATA_SECRET_KEY")
	client := pinata.NewClient(key, secret)

	meta := pinata.NewMetadataWithName("such_wow")
	meta.SetKeyValue("time_thing", time.Now())
	meta.SetKeyValue("string_thing", "much_awesome")
	meta.SetKeyValue("int_thing", 123)
	meta.SetKeyValue("float_thing", 123.456)

	hash := "QmZULkCELmmk5XNfCgTnCyFgAVxBRBXyDHGGMVoLFLiXEN"
	resp, err := client.PinHashToIPFSWithMetadata(hash, meta)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))

}
