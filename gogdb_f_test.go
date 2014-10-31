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
}
