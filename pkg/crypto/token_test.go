package sparklecrypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type A struct {
	ID string `json:"id"`
}

func TestBase64TokenSigner(t *testing.T) {
	signer := NewBase64TokenSigner()
	var aa = A{
		ID: "foo",
	}
	output, err := signer.Sign(aa)
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
	t.Logf("token: %s", output)

	var aaa A
	err = signer.Parse(output, &aaa)
	assert.EqualValues(t, aaa.ID, aaa.ID)
}

func TestJWTTokenSigner(t *testing.T) {

	signer := NewJWT("testest")
	var aa = A{
		ID: "foo",
	}
	output, err := signer.Sign(aa)
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
	t.Logf("token: %s", output)

	var aaa A
	err = signer.Parse(output, &aaa)
	assert.EqualValues(t, aaa.ID, aaa.ID)

}
