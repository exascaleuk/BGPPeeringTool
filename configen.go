package main

import "fmt"
import "net/http"
import "html/template"
import "os"
import "encoding/json"
import "io/ioutil"
import "strings"
import "flag"

type NetIxLan struct {
	ID      int    `json:"id"`
	IXID    int    `json:"ix_id"`
	Name    string `json:"name"`
	IXLanID int    `json:"ixlan_id"`
	Notes   string `json:"notes"`
	Speed   int    `json:"speed"`
	ASN     int    `json:"asn"`
	IPAddr4 string `json:"ipaddr4"`
	IPAddr6 string `json:"ipaddr6"`
	RSPeer  bool   `json:"is_rs_peer"`
	Created string `json:"created"`
	Updated string `json:"updated"`
	Status  string `json:"status"`
}

type PeeringDBResponse struct {
	Meta string    `json:"meta"`
	Data []Network `json:"data"`
}

type Network struct {
	ID              int        `json:"id"`
	OrgID           int        `json:"org_id"`
	Name            string     `json:"name"`
	AKA             string     `json:"aka"`
	Website         string     `json:"website"`
	ASN             int        `json:"asn"`
	LookingGlass    string     `json:"looking_glass"`
	RouteServer     string     `json:"route_server"`
	IrrAsSet        string     `json:"irr_as_set"`
	InfoType        string     `json:"info_type"`
	Prefixes4       int        `json:"info_prefixes4"`
	Prefixes6       int        `json:"info_prefixes6"`
	Traffic         string     `json:"info_traffic"`
	Ratio           string     `json:"info_ratio"`
	Scope           string     `json:"info_scope"`
	Unicast         bool       `json:"info_unicast"`
	Multicast       bool       `json:"info_multicast"`
	Ipv6            bool       `json:"info_ipv6"`
	Notes           string     `json:"notes"`
	PolicyURL       string     `json:"policy_url"`
	PolicyGeneral   string     `json:"policy_general"`
	PolicyLocations string     `json:"policy_locations"`
	PolicyRatio     string     `json:"policy_ratio"`
	PolicyContracts string     `json:"policy_contracts"`
	NetIxLan        []NetIxLan `json:"netixlan_set"`
	PocSet          string     `json:"poc_set"`
	Created         string     `json:"created"`
	Updated         string     `json:"updated"`
	Status          string     `json:"status"`
}

type Peers struct {
	Peers []NetIxLan
}

func getNetwork(body []byte) (*PeeringDBResponse, error) {
	var s = new(PeeringDBResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return s, err
}

func getNetworkInfo(asn string) Network {

	res, err := http.Get("https://www.peeringdb.com/api/net?asn=" + asn + "&depth=2")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	s, err := getNetwork([]byte(body))

	return s.Data[0]

}

func (p Peers) HasIPv4() bool {
	for _, element := range p.Peers {
		if len(element.IPAddr4) > 0 {
			return true
		}
	}
	return false
}

func (p Peers) HasIPv6() bool {
	for _, element := range p.Peers {
		if len(element.IPAddr6) > 0 {
			return true
		}
	}
	return false
}

func presentOnIX(n Network, ix NetIxLan) bool {
	for _, net := range n.NetIxLan {
		if net.IXID == ix.IXID {
			return true
		}
	}
	return false
}

func createPeeringList(r Network, l Network) Peers {

	var peer []NetIxLan

	for _, net := range r.NetIxLan {
		if presentOnIX(l, net) {
			peer = append(peer, net)
		}
	}

	thePeer := Peers{
		Peers: peer,
	}

	return thePeer
}

func FriendlyIXName(s string) string {
	result := strings.Replace(s, " ", "", -1)
	result = strings.Replace(result, ":", "", -1)
	result = strings.ToLower(result)

	return result
}

func FriendlyNetName(s string) string {

	s = strings.Replace(s, " ", "", -1)

	if len(s) > 6 {
		s = s[0:6]
	}

	return s
}

func main() {

	localPtr := flag.String("local", "61049", "Your Local ASN")
	remotePtr := flag.String("remote", "16509", "A Remote ASN i.e. Amazon")
	maxprefix := flag.Int("maxprefix", 100, "Override IPv4/IPv6 Maximum Prefix")
	MD5 := flag.String("md5", "", "MD5 Password for sessions")
	templatePtr := flag.String("template", "cisco.tpl", "Template to use when generating config")

	flag.Parse()

	l := getNetworkInfo(string(*localPtr))
	r := getNetworkInfo(string(*remotePtr))

	type tVars struct {
		Remote           Network
		Local            Network
		Peers            Peers
		DefaultMaxPrefix int
		MD5 			 string
	}

	theVars := tVars{
		Remote:           r,
		Local:            l,
		Peers:            createPeeringList(r, l),
		DefaultMaxPrefix: int(*maxprefix),
		MD5:			  string(*MD5),
	}

	funcMap := template.FuncMap{
		"FriendlyIXName":  FriendlyIXName,
		"FriendlyNetName": FriendlyNetName,
	}

	strTemplateFilename := "./templates/" + *templatePtr

	if _, err := os.Stat(strTemplateFilename); os.IsNotExist(err) {
		fmt.Println("Whoops: ", err)
	}

	t := template.Must(template.New("main").Funcs(funcMap).ParseFiles("./cisco.tpl"))
	t, _ = t.ParseFiles("./cisco.tpl")
	t.ExecuteTemplate(os.Stdout, "./cisco.tpl", theVars)
}
