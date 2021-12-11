package regDepatcher

var regList map[int]bool

func RegInit() {
	regList = make(map[int]bool)
	for i := 0; i < 32; i++ {
		regList[i] = true
	}
}

func NextAvailReg() int {
	for i := 0; i < 32; i++ {
		if regList[i] {
			regList[i] = false
			return i
		}
	}
	return -1
}

func OccupyReg(regId int) {
	regList[regId] = false
}

func ReleaseReg(regId int) {
	regList[regId] = true
}

var printExist, printlnExist, scanExist bool

func IOInit() {
	printExist = false
	printlnExist = false
	scanExist = false
}

func SetPrint() {
	printExist = true
}

func GetPrint() bool {
	return printExist
}

func SetPrintln() {
	printlnExist = true
}

func GetPrintln() bool {
	return printlnExist
}

func SetScan() {
	scanExist = true
}

func GetScan() bool {
	return scanExist
}
