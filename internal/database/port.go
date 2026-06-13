package database

func UsedPorts(excludeID ...uint) []int {
	var ports []int

	var config Config
	DB.First(&config)
	ports = append(ports, config.Port, config.SubPort)

	var inbounds []Inbound
	if len(excludeID) > 0 {
		DB.Where("id != ?", excludeID[0]).Find(&inbounds)
	} else {
		DB.Find(&inbounds)
	}
	for _, ib := range inbounds {
		ports = append(ports, ib.Port)
	}

	return ports
}
