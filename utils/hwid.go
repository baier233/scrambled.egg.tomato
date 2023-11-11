package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/super-l/machine-code/machine"
)

func GetHWID() string {
	machineData := machine.GetMachineData()
	myMD5 := md5.New()
	myMD5.Write([]byte(machineData.CpuId + machineData.PlatformUUID + machineData.SerialNumber))
	return hex.EncodeToString(myMD5.Sum(nil))
}
