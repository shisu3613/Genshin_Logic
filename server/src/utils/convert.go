package utils

const UIDPIDDIFF = 100000000

func PidToUid(pid int) int {
	return pid + UIDPIDDIFF
}

func UidToPid(uid int) int {
	return uid - UIDPIDDIFF
}
