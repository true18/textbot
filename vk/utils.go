package vk

// CheckError паникует, если err != nil
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
