package sdk

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

var BLS_TEST_DATA = map[string][]string{
	"message": {
		"hello foo",
		"Brother, I love you",
		"Jo Jo",
		"Hello Kitty",
		"Ha Haha Hahaha Hahahaha Hahahahaha Hahahahahaha Hahahahahahaha",
	},
	"ikm": {
		"0000000000000000000000000000000000000000000000000000000000000000",
		"15923F9D2FFFB11D771818E1F7D7DDCD363913933264D58533CB3A5DD2DAA66A",
		"A9323CEF24497AB770516EA572A0A2645EE2D5A75BC72E78DE534C0A03BC328E",
		"D7D816DA0566878EE739EDE2131CD64201BCCC27F88FA51BA5815BCB0FE33CC8",
		"27FC9998454848B987FAD89296558A34DEED4358D1517B953572F3E0AAA0A22D",
	},
	"privateKey": {
		"3562DBB3987D4FEB5B898633CC9C812AEC49F2C64CAD5B34F5A086DF199A124D",
		"B49CAEC8626B832103C7B999D001B911A215E8942EB4C7FED235B6DBC3E48E69",
		"D5D8B9ABF20D6493BD21162415C1B756F19188BD36FE93C5339E8F2E6E5C614E",
		"426AD2746B63CB08EBE6C16B4E2B9C47BEBD8648135FC720EAAA886A0323AA36",
		"A364026A6F25E375553646C4A29A51C635F669DD68B6188E2D72CEED95516E37",
	},
	"publicKey": {
		"A695AD325DFC7E1191FBC9F186F58EFF42A634029731B18380FF89BF42C464A42CB8CA55B200F051F57F1E1893C68759",
		"A389E43A21A4F2C2CC465F2CB666FD8C5BEEBCBB05547DA36121035DF1D0FF9BDE2583B5F1886A0A66CD729BC619E770",
		"8504A48E1116F51D5857C5E281CD4EACF196C7A288ED55546C3E0B16FADFEFC95D0E947DF2D483310CDCE5836DD5DCB9",
		"B0900B51CD5FE877EA91248537D787E90883C85361DE79A25A460741EFB29DFEC86B06D028C0FFE9B02671FDA7113538",
		"AF6B154FC92D1EFDE3A75D297F0119581E24D67F895DB7CB252B8BEF5BA93054B69881688CB35FFFB86993251E4C88EA",
	},
	"signature": {
		"864CBD674C5657256671C3BBDA0EF50D077FFE50EB03AC2C9C1F435B5BF5D872603BCD0A337665A0EF867397F3985EAD0EB347E6CAE7E0BA4F0F04355D8886899393AC82F1F895414ADDFD683054316EECEEAEA85F9C309666BBCFC922F2728A",
		"B440DAFDABB1BD03EAFD5B230E60468D10D295370A6E1EE918D2231684D22B74780835AD6EB7E7F7734557662923860C183D1BE43BEF1998BAC2175D14D4782110E4E943EAAD8DFC57CCF26BE84B443548CEC761C33894CA20263DCD793B1FAE",
		"986089A075412DAE9522AE2AEC1DA072D89A685848B0A18DA62C9FCFAA4B9706E139ED8EE99D7AC786E8AAC938CAE3FD15687CC666113AEF43AB69234488B646A1AEDDA60D7C1B253E4955ED51572FABCF6394C31C82BD00F794662E1863F27E",
		"8E1831E42FA1C3BFD7A2A346D84356F08D5AFBA01AB2CB472BEB7AFA4012B4113015A8BEB1A04E506D5976C118216156086CE1836AB4BC1AE14CFCC7FF5D32BD11B1E15C5B6CAA11EC57B1C40BE657A8724889CA87E353C8701902EE9FDD8656",
		"A5649669AB0162C7219CC91A5B3BB0DB1DF2CF2B597702C404AEBF3DDCA9CBD4081375C585F38B392BAA443D48F9DE6703D3516D3408B4ABEC2062C5B64CBF611C046A2FE6CED0788D1F7FDE6AE7277FB25257790D5C9713CF0F37A661862576",
	},
}

func TestGeneratePrivateKeyFromIKM(t *testing.T) {
	ikms := BLS_TEST_DATA["ikm"]
	sks := BLS_TEST_DATA["privateKey"]

	for i, _ := range ikms {
		b, _ := hex.DecodeString(ikms[i])
		var ikm [32]byte
		copy(ikm[:], b[:])
		sk := GeneratePrivateKeyFromIKM(ikm)

		assert.Equal(t, sk.HexString(), sks[i])
	}
}

type ZeroReader struct {
	t []byte
}

func (z *ZeroReader) Read(p []byte) (n int, err error) {
	for i, _ := range p {
		p[i] = z.t[i]
	}
	return len(p), nil
}

func TestGeneratePrivateKey_Seed(t *testing.T) {
	ikms := BLS_TEST_DATA["ikm"]
	sks := BLS_TEST_DATA["privateKey"]

	for i, _ := range ikms {
		b, _ := hex.DecodeString(ikms[i])
		sk := GeneratePrivateKey(&ZeroReader{b})

		assert.Equal(t, sk.HexString(), sks[i])
	}
}

func TestGeneratePrivateKey_Rand(t *testing.T) {
	sk1 := GeneratePrivateKey(&ZeroReader{make([]byte, 32)})
	sk2 := GeneratePrivateKey(nil)
	assert.NotEqual(t, sk1.HexString(), sk2.HexString())
}

func TestPrivateKey_Public(t *testing.T) {
	ikms := BLS_TEST_DATA["ikm"]
	sks := BLS_TEST_DATA["publicKey"]

	for i, _ := range ikms {
		b, _ := hex.DecodeString(ikms[i])
		sk := GeneratePrivateKey(&ZeroReader{b})

		assert.Equal(t, sk.Public().HexString(), sks[i])
	}
}

func TestPrivateKey_Sign(t *testing.T) {
	ikms := BLS_TEST_DATA["ikm"]
	msgs := BLS_TEST_DATA["message"]
	sigs := BLS_TEST_DATA["signature"]

	for i, _ := range ikms {
		b, _ := hex.DecodeString(ikms[i])
		sk := GeneratePrivateKey(&ZeroReader{b})

		assert.Equal(t, sk.Sign(msgs[i]).HexString(), sigs[i])
	}
}

func TestKeyPair_Verify(t *testing.T) {
	ikms := BLS_TEST_DATA["ikm"]
	keys := BLS_TEST_DATA["publicKey"]
	msgs := BLS_TEST_DATA["message"]
	sigs := BLS_TEST_DATA["signature"]

	for i, _ := range ikms {
		b, _ := hex.DecodeString(ikms[i])
		kp := GenerateKeyPair(&ZeroReader{b})
		sig := kp.Sign(msgs[i])

		assert.Equal(t, sig.HexString(), sigs[i])
		assert.Equal(t, kp.PublicKey.HexString(), keys[i])
		assert.True(t, sig.Verify(msgs[i], kp.PublicKey))
		assert.True(t, kp.PublicKey.Verify(msgs[i], sig))
	}
}

func TestAggregateVerify(t *testing.T) {
	ikms := BLS_TEST_DATA["ikm"]
	msgs := BLS_TEST_DATA["message"]

	signatures := make([]BLSSignature, len(ikms))
	publicKeys := make([]BLSPublicKey, len(ikms))
	for i, _ := range ikms {
		b, _ := hex.DecodeString(ikms[i])
		kp := GenerateKeyPair(&ZeroReader{b})
		publicKeys[i] = kp.PublicKey
		signatures[i] = kp.Sign(msgs[i])
	}

	aggregateSig, err := AggregateSignatures(signatures...)
	assert.Nil(t, err)
	assert.Equal(t, aggregateSig.HexString(), "91C5891B27541EF99BB2441BF9138C59103B78EDBA6BE72BF28576D343B75B95E14749950653FFF86816C9F1654E9D100521023EBB2A664530E6674EA6AE35E7F5184FBFA81436C10C437D50DD1460F0FDCC91160605748164AD21763D0E1462")

	assert.True(t, AggregateVerify(publicKeys, msgs, aggregateSig))
}

func TestFastAggregateVerify(t *testing.T) {
	ikms := BLS_TEST_DATA["ikm"]
	message := "It is same message for all signers to verify that fast aggregate verify works properly"

	signatures := make([]BLSSignature, len(ikms))
	publicKeys := make([]BLSPublicKey, len(ikms))
	for i, _ := range ikms {
		b, _ := hex.DecodeString(ikms[i])
		kp := GenerateKeyPair(&ZeroReader{b})
		publicKeys[i] = kp.PublicKey
		signatures[i] = kp.Sign(message)
	}

	aggregateSig, err := AggregateSignatures(signatures...)
	assert.Nil(t, err)
	assert.Equal(t, aggregateSig.HexString(), "B98B4E81319EAA7BB4BF3BC123697EDDFC888A52D7CD8FC4F260639D448C6409E4BC4437A3DA69A72D99CA153FB466FF15F2EE567DDF667F4F88A6BCA92C5CAA5B2C359DBF9367703D88D1974501D0CE0A47012AB4269498CBE076B763BE0B83")

	aggregatedPub, err := AggregatePublicKeys(publicKeys...)
	assert.Nil(t, err)
	assert.Equal(t, aggregatedPub.HexString(), "B1BB94B73381FE39C7D25C3E1353274D34D09CE698E04548B7F6DD49C062DAF17DA13405C92CB61C04508B991576183D")

	assert.True(t, FastAggregateVerify(publicKeys, message, aggregateSig))
}
