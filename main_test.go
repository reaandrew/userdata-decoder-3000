package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	base64EncodedHelloWorld       = "SGVsbG8sIFdvcmxkIQ=="
	base64EncodedZippedHelloWorld = "H4sIAAAAAAAAA/NIzcnJ11EIzy/KSVEEANDDSuwNAAAA"
	helloWorldExpected            = "Hello, World!"
)

func TestDecodeBase64(t *testing.T) {
	result, err := decodeBase64([]byte(base64EncodedHelloWorld))
	assert.Nil(t, err)
	assert.Equal(t, helloWorldExpected, string(result))
}

func TestIsGzipped(t *testing.T) {
	encryptedHelloWorld := base64EncodedZippedHelloWorld
	assert.True(t, isGzipped([]byte(encryptedHelloWorld)))
}

func TestUnzip(t *testing.T) {
	encryptedHelloWorld := base64EncodedZippedHelloWorld

	result, err := unzipData([]byte(encryptedHelloWorld))
	assert.Nil(t, err)
	assert.Equal(t, helloWorldExpected, string(result))
}

func TestDecode_With_Encrypted(t *testing.T) {
	encryptedHelloWorld := base64EncodedZippedHelloWorld
	result, err := decode([]byte(encryptedHelloWorld))
	assert.Nil(t, err)
	assert.Equal(t, helloWorldExpected, string(result))
}

func TestDecode(t *testing.T) {
	encryptedHelloWorld := base64EncodedHelloWorld
	result, err := decode([]byte(encryptedHelloWorld))
	assert.Nil(t, err)
	assert.Equal(t, helloWorldExpected, string(result))
}

func TestExtractMimeAttachments(t *testing.T) {
	/*  EXAMPLE MESSAGE
	Content-Type: multipart/mixed; boundary="===============6758256748430075474=="
	MIME-Version: 1.0

	--===============6758256748430075474==
	Content-Type: text/cloud-config; charset="utf-8"
	MIME-Version: 1.0
	Content-Transfer-Encoding: base64
	Content-Disposition: attachment; filename="cloud-config.yaml"

	IyBjbG91ZC1jb25maWcueWFtbAp3cml0ZV9maWxlczoKICAjIFdyaXRlIHRoZSBjb25maWd1cmF0
	aW9uIGZpbGUgdG8gL2V0Yy9zYW1wbGVfYXBwLwogIC0gcGF0aDogIi9ldGMvc2FtcGxlX2FwcC9z
	b21lLXlhbWwueWFtbCIKICAgIGVuY29kaW5nOiBneitiNjQKICAgIGNvbnRlbnQ6IEg0c0lDQ1ZF
	RkdVQUEzTnZiV1V0ZVdGdGJDNTVZVzFzQUUyUE1RL0NJQkNGOS82S0M3c0R0R3J0NkdwME1uRTAx
	M0lxS1VKRGFiWC9Yb3BpR29aM2ZPKzlDOXpjb0h4ZlpRQXJ3SzdURktjYVRUaHhiQjdrM0pTaFVV
	L1V2MkNEUHFxMDk2aWtxWHVnOFZsdGJSc3pYbmxOdkFMR2QyWEJFaEFCN0IyT0JDZDZ3Y1U2TGY5
	ZUhyeXpoWVBTR2hDT3RtbVZ1ZGZLaFlTaDNwTzhTdlE0NzI1cDRyTUM5RU1kTDhCRzFBTnh0cUFp
	VWNHK0hiSHM1TW5ObDUwaTBmRGt4bG9ubFVGUHYwKy9LK0FpaHFjSzhpS3g5U2F4YlpsOUFLdlMy
	SnhPQVFBQQoKICAjIFdyaXRlIHRoZSB0ZW1wbGF0ZSBmaWxlIHRvIC91c3Ivc2hhcmUvc2FtcGxl
	X2FwcC8KICAtIHBhdGg6ICIvdXNyL3NoYXJlL3NhbXBsZV9hcHAvc29tZS10ZXh0LnR4dCIKICAg
	IGVuY29kaW5nOiBneitiNjQKICAgIGNvbnRlbnQ6IEg0c0lDSFZFRkdVQUEzTnZiV1V0ZEdWNGRD
	NTBlSFFBQXdBQUFBQUFBQUFBQUE9PQoKcnVuY21kOgogIC0gWyIvdXNyL2xvY2FsL2Jpbi9zdGFy
	dC5zaCJd

	--===============6758256748430075474==
	Content-Type: text/x-shellscript; charset="utf-8"
	MIME-Version: 1.0
	Content-Transfer-Encoding: base64
	Content-Disposition: attachment; filename="start.sh"

	ZWNobyAiU3RhcnRlZCI=

	--===============6758256748430075474==--
	*/

	mimeMessage := `
H4sICAtIFGUAA3VzZXJkYXRhAL2U226jOhSG73mKqPfpGEja0CoX5WBCEkgxZ+6MTYHEHBRIAjz90IOqGe25mK0tbSQLa9nL6/+tz0upqy6turk7NOnTrLywrmjwuftRFn1Kn2dJfakoPg/ru/Xv38PjciUsHx4Xq4UIwONy8bhYr+840zC1uZ+e26Kunmb8PeC4+fxvUjnlNyFd2nc/CKsvdE7q6q3Inmckx+c27dZ3l+5tvvpTqe8jzrhq39LzXKtITYsqe5oluE0fFt871KJt6rboPlJx12GSl1P8efZWsLTCZbq++7X4/YBLdsdxxiAfE13iY4U/JsKyxAG5pAHskpdGJCUDsS9NsZ6Rsd4ZysvRgHTAIWLGBtWxI3/lUJ6UEHA4kC6GHjeJ7mVUX2V7wQfRII1RwN8S3X+LQvm2v9WZoYCM6BBgdZoXEqO6eSUC7Ijes1CAN6JII5cIPNuHLE+C26cixXhXkBm6f4kE6YSDZXUo5CqdPFtH+2vNuiYVYkllPxhaBghgqs3HkEMn6tueNrpVXPi8P9miOtW3quX6sT/C0fa8wdN4tAfW1j5Z+sFZCQ4wRQIQQOIWWCd6MzX+gtyXnjMB6x3e3yEdFoEiRYkoN0iQsCnEr7uRqYewOSYgX8SsQXa4FZ2ReujUHSMfeXiTF/aRnk3QuBj6PrfnfcE8WZq3gb2pngQcdH2w8asDjNtJoUPEJkrK/kBP0ESCPQTaFmINKhYwBhds1ViNxYj3BFePllzs5UMaNnXgyy4SctUVUZf4MR/r8R5D5mLVurlj7lKmAWs0eKKiwfXMJfL4k6vmChqh7FY5IB5sOD+wNjuQF87G5N1geUhU74ZduUR6t5iQuSTQ171NdNsN0n4HYIM38OiMeeOI2dIT4CJiTXvw4J4yc+CcKn+1fSjb9h85AvEHIhBM8w/epvjVUCSeiMaERp6T0vtGhPtkZPV+UGds5Jzq2YOhGFcaWsNetOoo3LLpnyeh3E4A52TzMiVLXezwIA5zsK/Qgn7hxP1bnhwYw3/gpNHA0pHKWa7MHDjZDKlse/CXoUmvk3VSvRfjT4fs8xUEw5dqob9GAmz3wrZJCmmkOhw4qixHrGzpf2k4/bzNU8Zaci6a7n/vOG03Nd77Np8aTRxYdTK8FJ6IcjJdaawY6781Np9z3E+mf/Qm1gUAAA==
`
	attachments, err := extractMimeAttachments([]byte(mimeMessage))

	assert.Nil(t, err)
	assert.Len(t, attachments, 2)
	assert.Equal(t, attachments[0].ContentType, "text/cloud-config; charset=\"utf-8\"")
	assert.Equal(t, attachments[1].ContentType, "text/x-shellscript; charset=\"utf-8\"")
}
