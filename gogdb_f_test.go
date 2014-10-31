package godbg

func prbgtest() {
	Pdbgf("prbgtest content")
}

func prbgtestCustom(pdbg *Pdbg) {
	pdbg.Pdbgf("prbgtest content2")
}

func (pdbg *Pdbg) pdbgTestInstance() {
	pdbg.Pdbgf("pdbgTestInstance content3")
}

func globalPdbgExcludeTest() {
	Pdbgf("calling no")
	globalNo()
}

func globalNo() {
	Pdbgf("gcalled1")
	globalCNo()
}

func globalCNo() {
	Pdbgf("gcalled2")
}

func customPdbgExcludeTest(pdbg *Pdbg) {
	pdbg.Pdbgf("calling cno")
	customNo(pdbg)
}

func customNo(pdbg *Pdbg) {
	pdbg.Pdbgf("ccalled1")
	customCNo(pdbg)
}

func customCNo(pdbg *Pdbg) {
	pdbg.Pdbgf("ccalled2")
}
