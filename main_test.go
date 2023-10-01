package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	base64EncodedHelloWorld       = "SGVsbG8sIFdvcmxkIQ=="
	base64EncodedZippedHelloWorld = "H4sIAAAAAAAAA/NIzcnJ11EIzy/KSVEEANDDSuwNAAAA"
	helloWorldExpected            = "Hello, World!"
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

	mimeMessage = `
H4sICFuIGWUAA3VzZXJkYXRhAL2USXOjOBTH73wKV+5OC7DTISkfzCKMbXAjQCw3IQhgi2UMXuDTtzqdSk1P9aGnpmp0kepJb/m/+ulpbTPkzTD3xy5/mdUXNlQdOQ9f6uqeZ6+ztL00GTmPq4fVr2u5kL8+K8ulAhaKoiyWXyWwWj0ItmUbc5yf+6ptXmbiIxCE+fxPXAXtl0KG/D58oay9ZHPaNm9V8TqjJTn3+bB6uAxv8+ffpfoMcSZN/5af50ZD26xqipdZSvr8afH5Qq/6ru2r4d2VDAOhZc3tr7O3iuUNqfPVw9+TP46kZg+CYI3qMTUVMdHEYyotaxLSSx7CIV13Mq0ZSLDCbXdGp3ZnaeujBbORRIhZG9Qmnvrhk4m0hkAgoXKxzKRLzaDIzOdiL2EQj8oUh+ItNfFbHKm3/a0tLA0U1ISA6PxcKSwz7SuV4EDNO4skeKOaMgmpJLJ9xMo0vP2sSLN+VFBYJr7EknIi4bI5VGqTc83O0f24c65pg1jauE+WUQAKmO6KCRTQKcNuYEx+k1RYxFxWZmbmVnd8nOAJTm4QjIEhoj1wtu7JMQ/es+QBW6YAASRvgXPKbrYhXpC/vgs2YHdPxDtkwirUlDiV1Q5JCrGl5NtuYvoh6o4pKBcJ65AbbWVvygJ0Go4xRgHZlJV7zM426HwCMRb2Ipbsk2MEG3i39ZNEwuEebnBzgEnPK/So3MVpfT9kJ2gjyR1DYwuJATUHWKMPtnqiJ3IsBpJvxkshCcoxj7o2xKqPpFL3ZTSkOBETM9kTyHyiOzd/Kv2MGcCZLJHqaPQDe4kC8eTrpYYmqPpNCWgAOwGHzmYHysrb2KIfLg+pHtyIr9bIHBYcmUsKsRls4ttuVPY7ADuygUdvKjtPLpaBBBcx6/pDAPcZs0fBa8pvLoaq6/6WI5C8IwIBP7/zxu1XS1NEKlscjbKkdfCJiPCTkecfgQZro5aZWTxZmnXNImfcy04bR1vG9zKN1J4DXNLNmjsrQ+KJIIlKsG/QIvvASfjXPEV4i8R/4GRkoWMiXXB8lXkQqt5UyAFEIQ4VP/bVK+/lXz4OFsRgNreZbuAwu062Id6qbgBVrejGLFwe09B90qp1sffUVOBf410/b/YxDm/XuGYXLgtweVyC01oMrP7LMLrP+zJnrKfnqhv+92nUD3woP/YlH0JJ6LTpuK4CGZWUtzvRrD8WNp8LwncQdAAY8gUAAA==`
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

func TestDecodeWithEncrypted(t *testing.T) {
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

	attachments, err := extractMimeAttachments([]byte(mimeMessage))

	assert.Nil(t, err)
	assert.Len(t, attachments, 2)
	assert.Equal(t, attachments[0].ContentType, "text/cloud-config; charset=\"utf-8\"")
	assert.Equal(t, attachments[1].ContentType, "text/x-shellscript; charset=\"utf-8\"")
}

func TestExtractWriteFiles(t *testing.T) {
	attachments, _ := extractMimeAttachments([]byte(mimeMessage))

	writeFiles, err := extractCloudConfigWriteFiles(attachments[0], "./")
	assert.Nil(t, err)
	assert.Len(t, writeFiles, 2)
}
