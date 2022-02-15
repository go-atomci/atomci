package main

import (
	"flag"
	"fmt"
	"github.com/astaxie/beego"
	"os"
)

const usageK = "AES Private Key"
const usageI = "AES IV"
const usageD = "Check IF integrate service config encrypted with specified KEY AND IV"
const usageC = "Encrypt CURRENT integrate service config (NOT Encrypted) with specifid KEY AND IV"

func main() {

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("-k %s\n", usageK)
		fmt.Printf("-i %s\n", usageI)
		fmt.Printf("-d %s\n", usageD)
		fmt.Printf("-c %s\n", usageC)
	}

	var key string
	var iv string
	var check bool
	var encrypt bool

	flag.StringVar(&key, "k", "", usageK)
	flag.StringVar(&iv, "i", "", usageI)
	flag.BoolVar(&check, "d", false, usageD)
	flag.BoolVar(&encrypt, "c", false, usageC)

	flag.Parse()

	if (!check && !encrypt) || (check && encrypt) {
		flag.Usage()
		os.Exit(1)
	}

	if check {
		fmt.Printf("%q", beego.AppConfig.String("notification::smtpHost"))
	}
}
