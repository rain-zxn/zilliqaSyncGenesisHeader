package tool

func (t *Tool) RegMethods() {
	t.RegMethod("SignatureData", SignatureData)
	t.RegMethod("SyncZilGenesisHeader", SyncZilGenesisHeader)
	t.RegMethod("GetZilGenesisHeaderMain", GetZilGenesisHeaderMain)
}
