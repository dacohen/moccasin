package moccasin

import (
	"errors"
	"runtime"
	"strings"
)

var errInvalidCallerFunction = errors.New("couldn't get caller function")

func parseFnName(fullName string) string {
	lastSlashIdx := strings.LastIndexByte(fullName, '/')
	if lastSlashIdx < 0 {
		lastSlashIdx = 0
	}
	lastDotIdx := strings.LastIndexByte(fullName[lastSlashIdx:], '.') + lastSlashIdx
	return fullName[lastDotIdx+1:]
}

func getFnName() (string, error) {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "", errInvalidCallerFunction
	}
	fun := runtime.FuncForPC(pc)
	return parseFnName(fun.Name()), nil
}
