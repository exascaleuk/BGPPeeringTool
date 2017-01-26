{{range $i, $x := $.Peers.Peers}}
route-map peering-{{$x.Name | FriendlyIXName}}-direct-in permit 10
 set metric 0
 set local-preference 200
 set community {{$.Local.ASN}}:{{.IXID}}
{{end}}

router bgp {{.Local.ASN}}
!
  neighbor {{.Remote.Name | FriendlyNetName}}v4 peer-group
  neighbor {{.Remote.Name | FriendlyNetName}}v4 remote-as {{.Remote.ASN}}
  neighbor {{.Remote.Name | FriendlyNetName}}v4 prefix-list bogons-filter-in in
  neighbor {{.Remote.Name | FriendlyNetName}}v4 filter-list 2 out
  neighbor {{.Remote.Name | FriendlyNetName}}v4 remove-private-as
  neighbor {{.Remote.Name | FriendlyNetName}}v4 maximum-prefix {{if gt .Remote.Prefixes4 0}}{{.Remote.Prefixes4}}{{else}}{{.DefaultMaxPrefix}} {{end}}
  !
  address-family ipv4
  {{range $i, $x := $.Peers.Peers}}
    neighbor {{.IPAddr4}} peer-group {{$.Remote.Name | FriendlyNetName}}v4
    neighbor {{.IPAddr4}} description Peering: {{$.Remote.Name | FriendlyNetName}} {{$x.Name | FriendlyIXName}}
    neighbor {{.IPAddr4}} route-map peering-{{$x.Name | FriendlyIXName}}-direct-in in
    neighbor {{.IPAddr4}} activate
  {{end}}
  exit
!
  neighbor {{.Remote.Name | FriendlyNetName}}v6 peer-group
  neighbor {{.Remote.Name | FriendlyNetName}}v6 remote-as {{.Remote.ASN}}
  neighbor {{.Remote.Name | FriendlyNetName}}v6 prefix-list bogons-filter-in in
  neighbor {{.Remote.Name | FriendlyNetName}}v6 filter-list 2 out
  neighbor {{.Remote.Name | FriendlyNetName}}v6 remove-private-as
  neighbor {{.Remote.Name | FriendlyNetName}}v6 maximum-prefix {{if gt .Remote.Prefixes6 0}}{{.Remote.Prefixes6}}{{else}}{{.DefaultMaxPrefix}}{{end}}
  !
  address-family ipv6
   {{range $i, $x := $.Peers.Peers}}
    neighbor {{.IPAddr6}} peer-group {{$.Remote.Name | FriendlyNetName}}v6
    neighbor {{.IPAddr6}} description Peering: {{$.Remote.Name | FriendlyNetName}} {{$x.Name | FriendlyIXName}}
    neighbor {{.IPAddr6}} route-map peering-{{$x.Name | FriendlyIXName}}-direct-in in
    neighbor {{.IPAddr6}} activate
    !
  {{end}}
  exit
  !
exit
!