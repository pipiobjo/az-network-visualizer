package service

import (
	"fmt"
	"github.com/emicklei/dot"
	"github.com/pipiobjo/az-network-visualizer/config"
	"net"
)

var networkClusterByVnetRange = make(map[string]*dot.Graph)
var vnet = make(map[string]*dot.Graph)
var subnet = make(map[string]*dot.Graph)
var routes = make(map[string]*dot.Graph)
var ipNodes = make(map[string]*dot.Node)

func FindVnetByAddress(ip string) *dot.Graph {
	parseIP := net.ParseIP(ip)

	for k, v := range subnet {
		_, ipNet, err := net.ParseCIDR(k)
		if err == nil {
			if ipNet.Contains(parseIP) {
				//fmt.Println("ip=" + ip + " is contained by network=" + ipNet.String() + " with cidr=" + cidr.String())
				return v
			}
		}
	}

	return nil
}

func DotService(network *config.Network) {
	const ATTR_VNET_RANGE = "vnet-range"

	g := dot.NewGraph(dot.Directed)
	internet := g.Node("Internet").Box()
	internet.Attr("label", dot.HTML("<b> Internet</b>"))

	// prepare clusters
	for _, v := range network.Data {

		//Network
		//fmt.Println(v.VNetName)
		//fmt.Println(v.VNetRange)
		var networkCluster *dot.Graph
		if _, ok := vnet[v.VNetName]; !ok {
			//fmt.Println("Creating cluster for vnet" + v.VNetName)
			networkCluster = g.Subgraph(v.VNetName, dot.ClusterOption{})
			networkClusterByVnetRange[v.VNetRange] = networkCluster

			networkCluster.Attr("label", dot.HTML("<b>"+v.VNetName+":</b>"+v.VNetRange+""))
			vnet[v.VNetName] = networkCluster

		} else {
			networkCluster = vnet[v.VNetName]
			//value := networkCluster.Value(ATTR_VNET_RANGE)
		}

		//Subnet
		//fmt.Println(v.SubnetRange)
		if _, ok := subnet[v.SubnetRange]; !ok {
			var subnetCluster *dot.Graph
			//fmt.Println("Creating cluster for subnet" + v.SubnetRange)
			//subnetCluster := networkCluster.Subgraph(v.SubnetRange, dot.ClusterOption{})
			subnetCluster = networkCluster.Subgraph(v.SubnetRange, dot.ClusterOption{})
			subnetCluster.Attr("label", dot.HTML("<b>Subnet: </b>"+v.SubnetRange+""))
			subnet[v.SubnetRange] = subnetCluster
			//networkCluster.Edge(subnetCluster)
		}

		//fmt.Println(v.RouteAddressPrefix)
		//fmt.Println(v.RouteName)
		//fmt.Println(v.RouteNextHopIPAddress)
		//fmt.Println(v.RouteNextHopType)
	}

	// creating edges
	for _, v := range network.Data {

		var subnetCluster *dot.Graph = subnet[v.SubnetRange]
		//Route
		if v.RouteName != "" {
			route := subnetCluster.Node(v.RouteName)
			route.Attr("label", dot.HTML("<b>routeName: </b>"+v.RouteName))

			routeTargetNodeCluster := networkClusterByVnetRange[v.RouteAddressPrefix]
			if routeTargetNodeCluster != nil {

				targetNode := ipNodes[v.RouteNextHopIPAddress]
				if targetNode == nil {
					targetNodeCluster := FindVnetByAddress(v.RouteNextHopIPAddress)
					newIpNode := targetNodeCluster.Node(v.RouteNextHopIPAddress)
					ipNodes[v.RouteNextHopIPAddress] = &newIpNode
					targetNode = &newIpNode

				}
				edge := route.Edge(*targetNode, v.RouteAddressPrefix)
				edge.Attr("label", dot.HTML("<b>route: </b>"+v.RouteAddressPrefix))
			} else if v.RouteNextHopType == "Internet" {
				route := subnetCluster.Node(v.RouteName)
				route.Attr("label", dot.HTML("<b>InternetRoute: </b>"+v.RouteName))
				route.Edge(internet, "Internet")

			} else if v.RouteAddressPrefix == "0.0.0.0/0" {

				targetNodeCluster := FindVnetByAddress(v.RouteNextHopIPAddress)
				if targetNodeCluster != nil {

					ipNode := ipNodes[v.RouteNextHopIPAddress]
					if ipNode == nil {
						newIpNode := targetNodeCluster.Node(v.RouteNextHopIPAddress)
						ipNodes[v.RouteNextHopIPAddress] = &newIpNode
						ipNode = &newIpNode

					}
					//fmt.Println("targetNode=" + ipNode.ID())
					edge := route.Edge(*ipNode, v.RouteName)
					edge.Attr("label", dot.HTML("<b>default: </b>"+v.RouteAddressPrefix))
				}
			}

		}
	}

	//n1 := g.Node("coding")
	//n2 := g.Node("testing a little").Box()
	//
	//g.Edge(n1, n2)
	//g.Edge(n2, n1, "back").Attr("color", "red")

	fmt.Println(g.String())
}
