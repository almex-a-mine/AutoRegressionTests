package pc

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
	"unsafe"
)

type AppFileInfo struct {
	DwSignature        uint32
	DwStrucVersion     uint32
	DwFileVersionMS    uint32
	DwFileVersionLS    uint32
	DwProductVersionMS uint32
	DwProductVersionLS uint32
	DwFileFlagsMask    uint32
	DwFileFlags        uint32
	DwFileOS           uint32
	DwFileType         uint32
	DwFileSubtype      uint32
	DwFileDateMS       uint32
	DwFileDateLS       uint32
}

var (
	modVersion                  = syscall.NewLazyDLL("version.dll")
	procGetFileVersionInfoSizeW = modVersion.NewProc("GetFileVersionInfoSizeW")
	procGetFileVersionInfoW     = modVersion.NewProc("GetFileVersionInfoW")
	procVerQueryValueW          = modVersion.NewProc("VerQueryValueW")
)

const GETLOCALIP_WAIT_SEC = 10000 // ローカルIPが取得できるまでの最大待機時間
const GETLOCALIP_WAIT_RETRY = 500 // ローカルIPが取得できるまでのリトライ間隔

// ローカルIPアドレス取得
// return: IPアドレスリスト,デバイスNo用IPbin
func GetLocalIpAddrInfo() ([]string, string) {
	var ipListTbl []string
	var ipbinDeviceNo string
	var addrs []net.Addr
	var err error

	// IPが取得できるまでリトライ
	for i := 0; i < GETLOCALIP_WAIT_SEC; i += GETLOCALIP_WAIT_RETRY {
		_, err = net.InterfaceAddrs()

		if nil == err {
			break
		} else {
			time.Sleep(GETLOCALIP_WAIT_RETRY * time.Millisecond)
		}
	}

	for {
		_, err = net.InterfaceAddrs()
		if err != nil {
			time.Sleep(500 * time.Millisecond)
		} else {
			break
		}
	}

	if addrs, err = net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				//if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				ipv4 := ipnet.IP.To4()
				if ipv4 == nil {
					continue
				}
				if ipv4[0] == 127 && ipv4[1] == 0 && ipv4[2] == 0 && ipv4[3] == 1 {
					continue
				}
				if ipv4[0] == 169 && ipv4[1] == 254 {
					continue
				}
				if len(ipbinDeviceNo) == 0 {
					ipbinDeviceNo = fmt.Sprintf("%02x%02x%02x%02x", ipv4[0], ipv4[1], ipv4[2], ipv4[3])
				}

				ipaddr := fmt.Sprintf("%v.%v.%v.%v", ipv4[0], ipv4[1], ipv4[2], ipv4[3])
				ipListTbl = append(ipListTbl, ipaddr)
			}
		}
	}
	if len(ipbinDeviceNo) == 0 {
		ipbinDeviceNo = fmt.Sprintf("%02x%02x%02x%02x", 0, 0, 0, 0)
		ipaddr := fmt.Sprintf("%v.%v.%v.%v", 0, 0, 0, 0)
		ipListTbl = append(ipListTbl, ipaddr)
	}
	return ipListTbl, ipbinDeviceNo
}

// FilePath,FileVersion,ProductVersion,ok
func GetAppVersion() (string, bool) {
	filename, err := os.Executable()
	if err != nil {
		return "", false
	}
	x, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return "", false
	}
	y, err := syscall.UTF16PtrFromString(`\`)
	if err != nil {
		return "", false
	}
	size, _, _ := procGetFileVersionInfoSizeW.Call(uintptr(unsafe.Pointer(x)), 0)
	// errを取得すると、成功時もerrにログがセットされる為利用する応答値でチェック
	if size <= 0 {
		return "", false
	}

	pBlock := make([]byte, size)

	ret, _, _ := procGetFileVersionInfoW.Call(
		uintptr(unsafe.Pointer(x)),
		0,
		uintptr(len(pBlock)),
		uintptr(unsafe.Pointer(&pBlock[0])),
	)
	// errを取得すると、成功時もerrにログがセットされる為利用する応答値でチェック
	if ret <= 0 {
		return "", false
	}

	var pSubBlock *uint16
	var puLen uint32
	r, _, _ := procVerQueryValueW.Call(
		uintptr(unsafe.Pointer(&pBlock[0])),
		uintptr(unsafe.Pointer(y)),
		uintptr(unsafe.Pointer(&pSubBlock)),
		uintptr(unsafe.Pointer(&puLen)),
	)
	// errを取得すると、成功時もerrにログがセットされる為利用する応答値でチェック
	if r <= 0 {
		return "", false
	}

	fi := *(*AppFileInfo)(unsafe.Pointer(pSubBlock))
	fVersion := fmt.Sprintf("%v.%v.%v.%v", fi.DwFileVersionMS>>16, fi.DwFileVersionMS&0xFFFF, fi.DwFileVersionLS>>16, fi.DwFileVersionLS&0xFFFF)

	fileInfo, err := os.Stat(filename)
	if err != nil {
		return "", false
	}

	return fmt.Sprintf("ファイルパス:%s ファイルバージョン:%s 最終更新日時:%v", filename, fVersion, fileInfo.ModTime().Format(time.RFC3339)), true
}
