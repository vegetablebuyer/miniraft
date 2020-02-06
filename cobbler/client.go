package cobbler

import (
	"github.com/kolo/xmlrpc"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

const (
	cobblerRpcUrl = "http://127.0.0.1:25151"
	secretFile    = "/var/lib/cobbler/web.ss"
)

var tokenData string

func init() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:25151")
	log.Println(tcpAddr)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err == nil {
		conn.Close()
	} else {
		log.Println("cobbler server is not running")
		panic(err)
	}
	data, err := ioutil.ReadFile(secretFile)
	if err != nil {
		panic(err)
	}

	tokenData = string(data)
	strings.TrimSpace(tokenData)
}

type Client struct {
	client       *xmlrpc.Client
	sharedSecret string
	token        string
}

func NewClient() *Client {
	client, err := xmlrpc.NewClient(cobblerRpcUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	c := &Client{
		client:       client,
		sharedSecret: tokenData,
	}
	return c
}

func (c *Client) GetSystem() []string {
	systems := make([]string, 0)
	err := c.client.Call("get_item_names", "system", &systems)
	if err != nil {
		panic(err)
	}
	return systems
}

func (c *Client) FindSystem(serialNumber string) bool {
	systems := c.GetSystem()
	for _, system := range systems {
		if system == serialNumber {
			return true
		}
	}
	return false
}

func (c *Client) login() {
	args := []interface{}{"", c.sharedSecret}
	err := c.client.Call("login", args, &c.token)
	if err != nil {
		panic(err)
	}
}

func (c *Client) logout() {
	args := []interface{}{c.token}
	err := c.client.Call("logout", args, nil)
	if err != nil {
		panic(err)
	}
}

func (c *Client) EditSystem(serialNumber string, arg map[string]interface{}) error {
	c.login()
	defer c.logout()
	args := []interface{}{"system", serialNumber, "edit", arg, c.token}
	err := c.client.Call("xapi_object_edit", args, nil)
	return err
}

func (c *Client) AddSystem(serialNumber string, arg []map[string]interface{}) error {
	c.login()
	defer c.logout()
	args := []interface{}{"system", serialNumber, "add", c.token}
	err := c.client.Call("xapi_object_edit", args, nil)
	return err
}

func (c *Client) RemoveSystem(serialNumber string) error {
	c.login()
	defer c.logout()
	args := []interface{}{"system", serialNumber, "remove", c.token}
	err := c.client.Call("remove_system", args, nil)
	return err
}
