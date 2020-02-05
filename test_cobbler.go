package main

import (
	"./cobbler"
	"fmt"
)

func main() {
	c := cobbler.NewClient()
	systems := c.GetSystem()
	fmt.Println(systems)
	arg := make(map[string]interface{})
	// arg["hostname"] = "test_name"
	arg["interface"] = "bond1"
	arg["interface_type"] = "bond"
	arg["bonding_opts"] = "miimon=100 mode=1"
	arg["static"] = "1"
	arg["ip_address"] = "192.168.110.110"
	arg["netmask"] = "255.255.255.0"
	arg["gateway"] = "192.168.110.1"
	c.EditSystem("816460142", arg)
}
