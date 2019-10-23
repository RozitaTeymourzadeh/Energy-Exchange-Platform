package data

import (
	"encoding/json"
	"fmt"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"strings"
)

type DeviceProfiles struct {
	Profiles []models.DeviceProfile
}

func NewDeviceProfiles() DeviceProfiles {
	return DeviceProfiles{
		Profiles: make([]models.DeviceProfile, 0),
	}
}

func DeviceProfilesFromJson(jsonBytes []byte) DeviceProfiles {
	dps := NewDeviceProfiles()
	err := json.Unmarshal(jsonBytes, &dps.Profiles)
	if err != nil {
		fmt.Println("Error in getting device list from JSON")
	}
	return dps
}

func (dps *DeviceProfiles) ShowDeviceProfiles() string {
	sb := strings.Builder{}
	sb.WriteString("Device Profiles: \n")
	for _, val := range dps.Profiles {
		sb.WriteString(val.Name + " : " + val.Id + "\n")
	}
	return sb.String()
}
