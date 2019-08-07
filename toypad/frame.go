package toypad

func checkSum(buf []byte) (sum byte) {
	for _, c := range buf {
		sum += c
	}
	// log.Printf("Checksum [% X]: %X", buf, sum)
	return
}
