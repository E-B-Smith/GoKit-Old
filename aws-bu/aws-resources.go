package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

func aws_resources_install_sql() ([]byte, error) {
	return bindata_read([]byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x00, 0xff, 0xac, 0x57,
		0x6d, 0x6f, 0xdb, 0x36, 0x17, 0xfd, 0x2c, 0xfd, 0x0a, 0x22, 0x28, 0x20,
		0xe9, 0x81, 0xdb, 0xda, 0x41, 0x1b, 0xb4, 0x71, 0x83, 0x07, 0x0d, 0xbc,
		0x15, 0x41, 0xbb, 0xad, 0x68, 0x5c, 0xf4, 0xc3, 0x30, 0x14, 0xb4, 0xc4,
		0xd8, 0xdc, 0x64, 0x52, 0x25, 0xa9, 0x24, 0x2e, 0xfa, 0xe3, 0x77, 0x48,
		0xea, 0xd5, 0xb6, 0x6c, 0x17, 0x58, 0x10, 0x10, 0x16, 0xef, 0x0b, 0x0f,
		0xcf, 0xb9, 0xbc, 0xa2, 0xc2, 0x4c, 0xc9, 0x82, 0x64, 0xd4, 0xd0, 0x05,
		0xd5, 0x8c, 0xf0, 0x3b, 0xc2, 0x1e, 0xb9, 0x36, 0x9a, 0xbc, 0xfd, 0x72,
		0xfb, 0xc7, 0xe2, 0x6f, 0x96, 0x1a, 0x3d, 0x0d, 0x53, 0xc5, 0xa8, 0x61,
		0xad, 0x57, 0x6b, 0x23, 0x0f, 0xdc, 0xac, 0x08, 0x13, 0xa9, 0xcc, 0xb8,
		0x58, 0x92, 0xe8, 0xf3, 0xfc, 0xd7, 0x57, 0xd1, 0x34, 0xf4, 0x59, 0xcd,
		0xa6, 0xd8, 0x9b, 0xf1, 0xd6, 0xd8, 0x6c, 0x29, 0xd5, 0x29, 0xcd, 0x58,
		0x93, 0xdd, 0x79, 0x6f, 0xf9, 0x50, 0x8d, 0xdc, 0xe5, 0x3a, 0x0c, 0xe2,
		0x30, 0x20, 0x11, 0x8c, 0x6e, 0xfa, 0x5d, 0x4e, 0x53, 0xce, 0x54, 0x14,
		0x06, 0xa3, 0x66, 0xee, 0x13, 0xd3, 0x46, 0x2a, 0x60, 0xd8, 0x33, 0xcb,
		0xb2, 0xde, 0x24, 0x06, 0x91, 0x51, 0x65, 0x27, 0x93, 0x06, 0x2a, 0x5d,
		0xe4, 0x7b, 0xb1, 0xce, 0xad, 0xa1, 0xc5, 0xe8, 0xdc, 0xfa, 0x46, 0x0f,
		0xee, 0x1f, 0xb6, 0x09, 0xf0, 0x77, 0x4f, 0x55, 0xba, 0xa2, 0x2a, 0x3e,
		0x7f, 0x79, 0x91, 0x90, 0x52, 0xf0, 0x6f, 0x25, 0x23, 0x42, 0x1a, 0x22,
		0xca, 0x3c, 0x27, 0x85, 0xe2, 0x6b, 0xaa, 0x36, 0xd6, 0x37, 0x04, 0x20,
		0xfa, 0xa0, 0x67, 0x48, 0x8a, 0x28, 0xc3, 0xd7, 0x00, 0x4a, 0xd7, 0x85,
		0xf9, 0xee, 0xe7, 0xaf, 0x37, 0x86, 0x69, 0x18, 0xb8, 0x30, 0x6c, 0xc9,
		0x94, 0x9f, 0x74, 0xe0, 0x31, 0xd9, 0xe7, 0xc8, 0x66, 0xca, 0x65, 0x4a,
		0xf3, 0xbd, 0xb9, 0x9c, 0x65, 0x27, 0x5b, 0x32, 0xac, 0xd0, 0x07, 0xb9,
		0xfc, 0xc0, 0xee, 0x59, 0x3e, 0x28, 0x4f, 0xe3, 0xb0, 0xa3, 0x0d, 0x2c,
		0x33, 0xb6, 0x28, 0x1b, 0x05, 0xf0, 0x7c, 0x23, 0xee, 0x64, 0xe7, 0x11,
		0x80, 0x95, 0xe9, 0x3c, 0xff, 0xf2, 0xc8, 0xbb, 0x8f, 0x5f, 0xa8, 0x12,
		0x1d, 0x05, 0xad, 0x83, 0x52, 0x52, 0x1d, 0xd1, 0x09, 0x6e, 0x03, 0x22,
		0xd5, 0x16, 0x0f, 0x91, 0x09, 0xa3, 0xac, 0x46, 0x9a, 0x29, 0x4e, 0x73,
		0xfc, 0x38, 0x24, 0x4f, 0x30, 0xb2, 0x3c, 0xf6, 0xe9, 0x0c, 0x82, 0xda,
		0x15, 0xe6, 0x42, 0xc9, 0x94, 0x69, 0x2d, 0x28, 0xbc, 0x6a, 0xd1, 0x27,
		0x17, 0x49, 0xdf, 0x87, 0x67, 0xa4, 0xa5, 0x3d, 0xe8, 0xd9, 0x72, 0xcb,
		0xa1, 0x17, 0xb3, 0x66, 0xb4, 0x67, 0xc7, 0xaa, 0x9a, 0x2e, 0x59, 0x5b,
		0x51, 0x2f, 0x27, 0xe7, 0x09, 0xe9, 0x78, 0x34, 0x94, 0x70, 0x91, 0xb1,
		0xc7, 0x5d, 0x4a, 0x80, 0xfb, 0xc6, 0x5a, 0x1a, 0x5a, 0xbc, 0xdf, 0x96,
		0x95, 0x48, 0xd1, 0x65, 0x2a, 0xb6, 0xdb, 0x1d, 0x79, 0xaa, 0x0e, 0x71,
		0x6e, 0x2b, 0xaf, 0xd4, 0x03, 0xb4, 0x77, 0x8c, 0x7d, 0xe6, 0x7f, 0x92,
		0xfa, 0x61, 0xee, 0x79, 0x16, 0x0c, 0xf1, 0xba, 0x28, 0x45, 0x96, 0xb3,
		0xfe, 0x51, 0x7c, 0x99, 0x04, 0xb0, 0x94, 0x05, 0xfa, 0x17, 0xcb, 0xaa,
		0xc3, 0xb0, 0xe0, 0x4b, 0xc4, 0xdb, 0xc5, 0xa4, 0xa1, 0xf9, 0x67, 0x67,
		0xf3, 0xa6, 0xc6, 0x92, 0xb1, 0x9c, 0x0d, 0x05, 0xcc, 0x9c, 0x6d, 0x2b,
		0xc0, 0x59, 0x3e, 0xb4, 0x27, 0xae, 0x31, 0x34, 0x62, 0xf6, 0xd4, 0x3c,
		0xac, 0xa1, 0x67, 0xf1, 0xe3, 0xcd, 0x6c, 0x40, 0xc5, 0xbe, 0xbd, 0xd2,
		0xb1, 0x43, 0x7d, 0x0c, 0x96, 0x6c, 0x7e, 0xbf, 0xc0, 0x1a, 0xb1, 0x96,
		0x7b, 0xfe, 0x9d, 0x65, 0xe4, 0x9e, 0xb3, 0x87, 0xc1, 0x5e, 0x37, 0xb7,
		0x9b, 0x68, 0x7b, 0xfe, 0x6e, 0xe0, 0x3e, 0x77, 0x34, 0x83, 0x50, 0x83,
		0x92, 0x14, 0xbb, 0xd5, 0x45, 0xce, 0xcd, 0xd7, 0x02, 0x27, 0x3d, 0x86,
		0x9c, 0x23, 0x12, 0x3d, 0x8f, 0x46, 0x64, 0x92, 0xd8, 0x7e, 0x71, 0xed,
		0xc4, 0x19, 0xc1, 0xa7, 0x5c, 0xc7, 0xa9, 0x7d, 0x93, 0x3c, 0xac, 0x98,
		0x20, 0xe8, 0x6d, 0x96, 0x7e, 0xf2, 0x86, 0xb8, 0x76, 0xe5, 0x7e, 0x1b,
		0x6b, 0x98, 0x10, 0x96, 0xc3, 0x69, 0x8c, 0x0a, 0xca, 0x5c, 0x02, 0xaf,
		0x53, 0xf6, 0x11, 0xc9, 0xf5, 0xe9, 0x69, 0xdc, 0xe3, 0xc2, 0x4a, 0x32,
		0x90, 0xcf, 0xc9, 0x35, 0x98, 0x8f, 0x6b, 0x5f, 0xa4, 0xfb, 0x21, 0xfd,
		0xce, 0x1e, 0x0e, 0xc3, 0xe9, 0x85, 0x0f, 0x43, 0x41, 0x9e, 0xfd, 0x30,
		0xda, 0xcd, 0x1c, 0x06, 0xe2, 0x4b, 0x72, 0x3f, 0x96, 0x81, 0x1c, 0x80,
		0xb8, 0x17, 0x4b, 0xa7, 0xba, 0x8f, 0xa4, 0xaa, 0x0f, 0xf0, 0x7e, 0x48,
		0xf3, 0xe6, 0x40, 0x9c, 0x02, 0xab, 0x97, 0x6b, 0x98, 0xa8, 0x79, 0xff,
		0x94, 0x1d, 0xe4, 0xfd, 0x04, 0x78, 0x28, 0xe7, 0xe3, 0xfa, 0xf5, 0xf2,
		0x0c, 0xb1, 0x56, 0xa7, 0x1b, 0x84, 0xa5, 0xdd, 0x9d, 0xe6, 0x6a, 0xf7,
		0x22, 0x33, 0x80, 0xcf, 0xf5, 0x93, 0xca, 0x67, 0x10, 0xe3, 0x91, 0xa4,
		0x43, 0x60, 0xbb, 0xb9, 0x1d, 0xe0, 0x30, 0x08, 0xee, 0x94, 0x5c, 0xdb,
		0x00, 0xe9, 0x8e, 0xb7, 0xf1, 0x0d, 0x3c, 0x58, 0x2a, 0x59, 0x16, 0x64,
		0xb1, 0x21, 0x03, 0x07, 0xdb, 0xf6, 0x98, 0xa7, 0x4f, 0xc9, 0xe9, 0x6d,
		0x66, 0x86, 0x47, 0xc5, 0x17, 0xa5, 0xe1, 0x52, 0x54, 0xaf, 0x10, 0xc4,
		0x1f, 0xef, 0x37, 0x3b, 0x71, 0xb6, 0xeb, 0x20, 0xb2, 0x6a, 0x3c, 0xf8,
		0x65, 0xff, 0xaa, 0x5b, 0xd5, 0xe5, 0x65, 0x56, 0x5d, 0x20, 0xcf, 0x66,
		0x74, 0x73, 0x36, 0x0a, 0xac, 0xfd, 0x94, 0xde, 0xe4, 0xdc, 0x40, 0x72,
		0x7d, 0x0b, 0xbb, 0xbc, 0xf4, 0x6d, 0xdc, 0x39, 0x9d, 0xb9, 0xa9, 0xb3,
		0x7a, 0x2d, 0x47, 0xd8, 0xd6, 0x6d, 0xd0, 0xc6, 0x43, 0x1b, 0xc5, 0x6a,
		0x24, 0xdd, 0x02, 0x72, 0xd6, 0x96, 0xd1, 0x3e, 0x56, 0x6f, 0x94, 0x2a,
		0x63, 0x6a, 0xd7, 0x98, 0x31, 0x9d, 0x4e, 0x3d, 0xd1, 0x78, 0xfc, 0x6a,
		0x54, 0x29, 0xd2, 0x38, 0xca, 0xe8, 0x06, 0x1b, 0xa8, 0x3c, 0xf1, 0x96,
		0x83, 0xbd, 0x61, 0xe3, 0xb4, 0xcd, 0xb6, 0xaa, 0xd7, 0xaf, 0x8a, 0x3b,
		0x64, 0xb6, 0x14, 0x77, 0xa4, 0x5b, 0x95, 0x6b, 0x2a, 0x3e, 0x31, 0x9a,
		0xd9, 0x1d, 0x3a, 0x0a, 0x62, 0x0d, 0x81, 0x48, 0xcd, 0xcc, 0xf6, 0x5d,
		0xb1, 0x49, 0x71, 0x2c, 0x50, 0x31, 0x53, 0x2a, 0x81, 0x6a, 0x64, 0x8f,
		0xc6, 0xaa, 0x19, 0x3c, 0x79, 0x12, 0x06, 0x19, 0x4b, 0x73, 0x0a, 0xfe,
		0xc2, 0x00, 0xd7, 0x4b, 0x92, 0xc9, 0xd2, 0x4a, 0x5d, 0x28, 0x96, 0x72,
		0x8d, 0x9c, 0xd3, 0xd0, 0xde, 0x1e, 0x8c, 0x0b, 0xc1, 0xef, 0x05, 0x43,
		0xa6, 0x10, 0xef, 0x77, 0xa0, 0x75, 0xa9, 0x7b, 0x0d, 0x0e, 0xe5, 0xeb,
		0x97, 0x70, 0x73, 0x70, 0x0f, 0x50, 0xff, 0xd8, 0xd8, 0xd4, 0x46, 0xd8,
		0xec, 0x57, 0x2e, 0x08, 0x1c, 0xef, 0xae, 0x72, 0x87, 0x0f, 0x07, 0x9a,
		0xae, 0x88, 0x5b, 0x8d, 0xe3, 0x24, 0x29, 0x45, 0x37, 0x7e, 0xfc, 0x33,
		0x22, 0xee, 0x50, 0x45, 0xa3, 0x88, 0xbc, 0xbf, 0xb6, 0xe3, 0x6f, 0x6e,
		0x7c, 0xe7, 0xc6, 0xb9, 0x1b, 0x3f, 0x5e, 0x47, 0x7f, 0xa1, 0x8b, 0xc9,
		0x02, 0x8b, 0x02, 0x9b, 0x5d, 0xec, 0x0d, 0x99, 0x8c, 0xcf, 0x5f, 0x3c,
		0x1b, 0x3b, 0x6c, 0x98, 0xae, 0xb1, 0x19, 0xf9, 0xd5, 0x5d, 0x0a, 0xe0,
		0x03, 0x95, 0x5e, 0xe3, 0x6f, 0xf6, 0x3a, 0x4a, 0x7e, 0xfc, 0xb0, 0x2b,
		0x5b, 0xcc, 0x2d, 0xe8, 0x0a, 0xb3, 0x1d, 0x9f, 0x57, 0xc9, 0xea, 0x4d,
		0xd9, 0xa5, 0x50, 0x20, 0x75, 0xce, 0xbc, 0xa0, 0x59, 0x1c, 0x81, 0xfa,
		0x3c, 0xdf, 0xfc, 0x1f, 0xba, 0xbf, 0x4a, 0xdc, 0x9e, 0xe1, 0x39, 0x25,
		0x96, 0x66, 0x0c, 0x39, 0x15, 0xcb, 0x12, 0xd7, 0x12, 0x52, 0xe4, 0xc5,
		0x52, 0x7f, 0xcb, 0x09, 0x5f, 0xaf, 0xcb, 0xea, 0xd4, 0xd7, 0xc2, 0x38,
		0x2a, 0xa5, 0xa7, 0x0f, 0x24, 0x14, 0xa5, 0x99, 0x1e, 0xa8, 0x13, 0x94,
		0xe3, 0x6d, 0x41, 0x45, 0x6c, 0x4b, 0xca, 0x95, 0x7e, 0xe7, 0x06, 0x37,
		0xc2, 0x36, 0xb7, 0xe7, 0x0e, 0x54, 0xce, 0xcf, 0xa7, 0xaa, 0x21, 0x57,
		0xd7, 0xc2, 0x9d, 0x72, 0xfa, 0x56, 0x52, 0x61, 0xa8, 0x33, 0xab, 0x7b,
		0x6a, 0x8b, 0xa1, 0xb6, 0x28, 0xa6, 0xcb, 0xdc, 0xd4, 0x81, 0x4d, 0x51,
		0xb9, 0x9a, 0x6a, 0xd6, 0xaf, 0xeb, 0x4a, 0xaa, 0x7a, 0xf5, 0x6e, 0xa5,
		0x1d, 0x28, 0x34, 0xb0, 0x5e, 0x2d, 0x7d, 0xd5, 0x66, 0x7b, 0x5a, 0x25,
		0x99, 0x5a, 0xaa, 0xdd, 0xea, 0x57, 0x20, 0xd1, 0x28, 0x9a, 0x9a, 0x98,
		0x15, 0x12, 0x65, 0xe7, 0x8e, 0xa5, 0x0f, 0x4c, 0x2e, 0x2f, 0x5b, 0x6c,
		0xd5, 0x32, 0x71, 0x15, 0xf6, 0x9c, 0xc4, 0x17, 0x63, 0xf2, 0x3f, 0xe2,
		0x86, 0xf3, 0x17, 0x49, 0xcf, 0xd7, 0xeb, 0xfd, 0xdf, 0xc8, 0xbd, 0xe7,
		0x2b, 0x00, 0xef, 0x24, 0x7c, 0xfe, 0x80, 0xce, 0x81, 0x0f, 0x81, 0xbe,
		0xdd, 0x7f, 0x0b, 0xdc, 0x33, 0x65, 0xcf, 0x17, 0xe9, 0xdc, 0x85, 0x27,
		0xe3, 0xa4, 0xb9, 0xb8, 0x17, 0xd4, 0xac, 0xb6, 0x3e, 0xa3, 0x61, 0x43,
		0xae, 0xb7, 0xa9, 0xfd, 0xdc, 0x7a, 0xcf, 0x36, 0x37, 0xb3, 0xd6, 0x7c,
		0xf1, 0xa2, 0x67, 0xbd, 0x65, 0x40, 0x60, 0x76, 0xad, 0x9f, 0x20, 0xa7,
		0x14, 0xdd, 0x15, 0x2f, 0xaa, 0xcb, 0x37, 0x17, 0xf8, 0x24, 0x71, 0xca,
		0xcb, 0x5d, 0xc0, 0x24, 0xae, 0xc0, 0x26, 0x04, 0xf5, 0x52, 0xe2, 0x1d,
		0x1a, 0x47, 0x93, 0x67, 0xe3, 0x31, 0xfe, 0x27, 0x91, 0x8d, 0xfe, 0x37,
		0x00, 0x00, 0xff, 0xff, 0xa4, 0xec, 0x13, 0x44, 0x3c, 0x11, 0x00, 0x00,
	},
		"aws-resources/install.sql",
	)
}

func aws_resources_uninstall_sql() ([]byte, error) {
	return bindata_read([]byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x00, 0xff, 0x84, 0x91,
		0xc1, 0x6a, 0xc3, 0x30, 0x10, 0x44, 0xef, 0xfd, 0x8a, 0x3d, 0xa6, 0xd0,
		0x3f, 0xe8, 0xa9, 0x25, 0x97, 0x40, 0xa0, 0xa1, 0x0e, 0xf4, 0x3c, 0xb6,
		0xd6, 0xe9, 0x16, 0x4b, 0x32, 0xde, 0x75, 0x1a, 0xe7, 0xeb, 0x2b, 0xb5,
		0x0d, 0x18, 0xc7, 0xc6, 0x57, 0xcd, 0xe3, 0xa1, 0x99, 0x75, 0x5d, 0x6c,
		0xc9, 0xc1, 0x50, 0x42, 0x99, 0xa4, 0x26, 0xbe, 0x88, 0x9a, 0xd2, 0xcb,
		0x47, 0xf1, 0x56, 0x7e, 0x71, 0x65, 0xfa, 0xfc, 0xf0, 0xcb, 0xd8, 0xd0,
		0xce, 0xe6, 0x85, 0xc1, 0x98, 0x2a, 0x68, 0x05, 0xc7, 0x37, 0x16, 0x65,
		0x33, 0x0b, 0x1f, 0x73, 0xb0, 0x28, 0xdc, 0xc7, 0xd3, 0x9e, 0xcf, 0xdc,
		0xac, 0xdb, 0x12, 0x39, 0x56, 0x49, 0x70, 0x7c, 0xb9, 0x27, 0xc4, 0xf3,
		0x2e, 0x27, 0xcb, 0x9e, 0xfc, 0xf9, 0x5e, 0x57, 0x54, 0x7f, 0xd0, 0x61,
		0xb7, 0x1d, 0xcb, 0x7c, 0x6a, 0xdd, 0x09, 0x1a, 0xb9, 0xb2, 0xa3, 0xb3,
		0xf0, 0xf7, 0x62, 0xdd, 0x63, 0x34, 0x34, 0xb7, 0x15, 0xeb, 0x3e, 0x54,
		0x26, 0x31, 0x8c, 0xe8, 0xcf, 0xde, 0x23, 0xbc, 0x33, 0x5c, 0x86, 0x5f,
		0x07, 0x63, 0xdd, 0x68, 0xb2, 0x52, 0x29, 0x27, 0x09, 0xf6, 0x38, 0x19,
		0x63, 0x46, 0xe0, 0x30, 0x14, 0x2d, 0xc2, 0xa6, 0xee, 0xa2, 0xdf, 0xe6,
		0x63, 0x58, 0x2a, 0xae, 0x06, 0xdf, 0xda, 0xf5, 0x89, 0x2c, 0x4e, 0xdf,
		0xa6, 0xca, 0x99, 0x5d, 0x0e, 0xe8, 0xe0, 0x39, 0x35, 0xfc, 0x9f, 0xe6,
		0x27, 0x00, 0x00, 0xff, 0xff, 0x02, 0xd3, 0xcb, 0xf0, 0x26, 0x02, 0x00,
		0x00,
	},
		"aws-resources/uninstall.sql",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"aws-resources/install.sql": aws_resources_install_sql,
	"aws-resources/uninstall.sql": aws_resources_uninstall_sql,
}
// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"aws-resources": &_bintree_t{nil, map[string]*_bintree_t{
		"uninstall.sql": &_bintree_t{aws_resources_uninstall_sql, map[string]*_bintree_t{
		}},
		"install.sql": &_bintree_t{aws_resources_install_sql, map[string]*_bintree_t{
		}},
	}},
}}
