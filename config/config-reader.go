package config

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"os"
)

type Network struct {
	Count int `json:"count"`
	Data  []struct {
		SubnetRange           string `json:"SubnetRange"`
		VNetName              string `json:"VNetName"`
		VNetRange             string `json:"VNetRange"`
		RouteAddressPrefix    string `json:"routeAddressPrefix"`
		RouteName             string `json:"routeName"`
		RouteNextHopIPAddress string `json:"routeNextHopIpAddress",omitempty`
		RouteNextHopType      string `json:"routeNextHopType"`
	} `json:"data"`
	SkipToken    interface{} `json:"skip_token"`
	TotalRecords int         `json:"total_records"`
}

func ReadInput(filePath string) *Network {
	config := Network{}

	// read yaml config
	_, err := os.Stat(filePath)
	if err != nil {
		log.Error().Err(err).Str("filePath", filePath).Msg("File does not exists")
		os.Exit(1)
	}

	yamlBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Error().Err(err).Str("filePath", filePath).Msg("Cannot read config file")
		os.Exit(1)
	}

	if os.IsNotExist(err) {
		log.Error().Err(err).Str("filePath", filePath).Msg("Config file does not exist")
		os.Exit(1)
	}

	err = json.Unmarshal(yamlBytes, &config)
	if err != nil {
		log.Error().Err(err).Str("filePath", filePath).Msg("Cannot parse config file")
		os.Exit(1)
	}
	return &config

}
