// Code generated by "stringer -type DeviceType"; DO NOT EDIT.

package main

import "strconv"

const _DeviceType_name = "GT_WT_01GT_WT_01_variant"

var _DeviceType_index = [...]uint8{0, 8, 24}

func (i DeviceType) String() string {
	if i >= DeviceType(len(_DeviceType_index)-1) {
		return "DeviceType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _DeviceType_name[_DeviceType_index[i]:_DeviceType_index[i+1]]
}
