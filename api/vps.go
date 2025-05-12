package api

import (
	"fmt"
	"github.com/edvin/oh/cache"
)

func ExecuteVirtualServerAction(vpsId int, action VirtualServerAction, body any) (VirtualServerActionResponse, error) {
	path := fmt.Sprintf("servers/%d/%s", vpsId, action.String())
	return Fetch[VirtualServerActionResponse]("POST", path, body, cache.NoCache)
}

func OrderVps(order CloudServerOrder) (CloudServerOrderResponse, error) {
	return Fetch[CloudServerOrderResponse]("POST", "servers/order", order, cache.NoCache)
}

func ListVpsFlavours(serverId int) ([]CloudServerFlavour, error) {
	path := fmt.Sprintf("servers/%d/possible-flavours", serverId)
	return Fetch[[]CloudServerFlavour]("GET", path, nil, cache.KeyFlavours.WithArg(serverId))
}

func ListVpsProducts() ([]Product, error) {
	return Fetch[[]Product]("GET", "products", nil, cache.KeyVpsProducts)
}

func ChangeVpsFlavour(serverId int, flavourId int) (ChangeFlavourResponse, error) {
	path := fmt.Sprintf("servers/%d/change-flavour", serverId)
	request := ChangeFlavourRequest{FlavourId: flavourId}
	return Fetch[ChangeFlavourResponse]("POST", path, request, cache.NoCache)
}

func ListVpsImages() ([]CloudServerImage, error) {
	return Fetch[[]CloudServerImage]("GET", "images", nil, cache.KeyVpsImages)
}

func GetVpsImage(imageId int) (CloudServerImage, error) {
	url := fmt.Sprintf("images/%d", imageId)
	return Fetch[CloudServerImage]("GET", url, nil, cache.NoCache)
}

func ListCloudServers() ([]CloudServer, error) {
	return Fetch[[]CloudServer]("GET", "servers", nil, cache.KeyCloudServers)
}

func GetVirtualServer(serverId int) (CloudServer, error) {
	url := fmt.Sprintf("servers/%d", serverId)
	return Fetch[CloudServer]("GET", url, nil, cache.NoCache)
}

func ListVirtualNetworks() ([]VirtualNetwork, error) {
	return Fetch[[]VirtualNetwork]("GET", "virtual-networks", nil, cache.KeyVirtualNetworks)
}

func ListAttachedVirtualNetworks(serverId int) ([]AttachedNetwork, error) {
	path := fmt.Sprintf("servers/%d/networks", serverId)
	networks, err := Fetch[[]AttachedNetwork]("GET", path, nil, cache.KeyAttachedNetworks.WithArg(serverId))
	return networks, err
}

func DetachVirtualNetwork(vpsId int, networkId string) (DetachVirtualNetworkResponse, error) {
	path := fmt.Sprintf("servers/%d/detach-network", vpsId)
	request := DetachVirtualNetworkRequest{NetworkId: networkId}
	return Fetch[DetachVirtualNetworkResponse]("POST", path, request, cache.NoCache)
}

func AttachVirtualNetwork(vpsId int, networkId string, ipv4 string, ipv6 string) (AttachVirtualNetworkResponse, error) {
	path := fmt.Sprintf("servers/%d/attach-network", vpsId)
	request := AttachVirtualNetworkRequest{
		NetworkId: networkId,
		IPv4:      ipv4,
		IPv6:      ipv6,
	}
	return Fetch[AttachVirtualNetworkResponse]("POST", path, request, cache.NoCache)
}
