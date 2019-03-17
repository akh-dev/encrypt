package engine

import (
	"bytes"
	"testing"

	"github.com/pkg/errors"
)

type testcase struct {
	plaintext []byte
	key       *[32]byte
}

func TestEncryptDecrypt(t *testing.T) {
	aes, err := newAESEngine()
	if err != nil {
		t.Fatalf("failed to create new aes engine : %s", err.Error())
	}

	testCases, err := setupTestData(aes)
	if err != nil {
		t.Fatalf("failed to setup test data : %s", err.Error())
	}

	for i, data := range testCases {
		ciphertext, err := aes.Encrypt(data.plaintext, data.key)
		if err != nil {
			t.Errorf("test case %d failed : failed to encode :%s", i, err.Error())
			continue
		}

		plaintext, err := aes.Decrypt(ciphertext, data.key)
		if err != nil {
			t.Errorf("test case %d failed : failed to decode :%s", i, err.Error())
			continue
		}

		if !bytes.Equal(plaintext, data.plaintext) {
			t.Errorf("test case %d failed : texts don't match. expected %s, got %s", i, data.plaintext, plaintext)
		}
	}

	incorrectKey, err := aes.GenerateNewKey()
	if err != nil {
		t.Fatalf("failed to generate new key : %s", err.Error())
	}

	ciphertext, err := aes.Encrypt(testCases[2].plaintext, testCases[2].key)
	if err != nil {
		t.Fatalf("test case %d failed : failed to encode :%s", 3, err.Error())
	}

	_, err = aes.Decrypt(ciphertext, incorrectKey)
	//t.Error(err)
	if err == nil {
		t.Error("test case failed : expected to receive error \"cipher: message authentication failed\" but got success")
	}

	//cipherStr := "eExBjrIiqouen6Mfy5BjItJv+CDotFikcotWCOlQxVHazDyQEzCB+HXt8B0OIXGk9Cdw+EPrMEMHjmc="
	//keyStr := "7pdenu5EBuR3RNqt9Poty0TypaJttOL3kJ9Zei8ebAA="
	//cipherBytes, _ := base64.StdEncoding.DecodeString(cipherStr)
	//keyBytes, _ := base64.StdEncoding.DecodeString(keyStr)
	//keyArray := [32]byte{}
	//for i, b := range keyBytes {
	//	keyArray[i] = b
	//}
	//plaintext, err := aes.Decrypt(cipherBytes, &keyArray)
	//t.Error(string(plaintext))
}

func setupTestData(aes *AESEngine) ([]testcase, error) {
	key1, err := aes.GenerateNewKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate new key")
	}

	key2, err := aes.GenerateNewKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate new key : %s")
	}

	key3, err := aes.GenerateNewKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate new key : %s")
	}

	testCases := []testcase{
		{
			plaintext: []byte("foo bar"),
			key:       key1,
		},
		{
			plaintext: []byte(""),
			key:       key2,
		},
		{
			plaintext: []byte(`Lorem ipsum dolor sit amet, consectetur adipiscing elit. Praesent ultrices libero lectus. Morbi ornare sed augue ut semper. Aliquam sed diam odio. Sed orci enim, bibendum euismod luctus et, dapibus vel neque. Vivamus non maximus risus. Fusce vel libero feugiat, iaculis nunc vitae, gravida leo. Nullam tellus elit, elementum fermentum sapien a, ullamcorper tempus leo. Nulla quis cursus nibh, ac condimentum justo. Duis commodo dui nisi, in consectetur orci suscipit eu. Nunc vitae enim a ligula dapibus dapibus. Ut feugiat dolor quis convallis egestas. Suspendisse malesuada, risus vitae ullamcorper malesuada, dolor justo ullamcorper risus, sit amet lobortis erat enim et odio. Proin non risus erat. Sed sapien mi, tempor nec fringilla eu, semper ut magna. Aenean lobortis velit sed nisi maximus mattis.
Quisque convallis, mauris et gravida volutpat, risus diam venenatis enim, a vestibulum massa metus vitae libero. Vivamus lobortis eu arcu a dignissim. Nulla nec maximus sapien, nec accumsan orci. Nulla condimentum nisl sed turpis pulvinar sagittis. Vivamus non ipsum eu nisi tempus elementum nec et lorem. Sed a maximus mauris. Sed nisl quam, interdum a ante sit amet, placerat vestibulum massa. Sed vehicula congue ipsum, et hendrerit ipsum viverra ac. Nunc quis leo porta, egestas eros commodo, consequat felis.
Vivamus luctus tellus et eleifend rhoncus. Etiam et dui volutpat, posuere justo in, consectetur dui. Phasellus non volutpat risus. Quisque volutpat erat velit, et aliquet urna malesuada vel. Donec elementum, nisl vel semper rutrum, felis lectus viverra tortor, nec cursus ligula mi at risus. Mauris consequat dapibus aliquet. Aenean efficitur tempor blandit. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
Donec imperdiet quam eu eros dictum luctus ac ac ex. Nullam suscipit tristique libero, sit amet convallis metus laoreet volutpat. Sed placerat vestibulum augue a pulvinar. Nullam eu nulla ante. Suspendisse potenti. Suspendisse placerat tincidunt urna, at lobortis ex rhoncus sit amet. Fusce in ligula id risus tincidunt vulputate. Quisque sapien quam, fermentum eget rhoncus sed, dictum vel dolor. Nullam venenatis libero nibh, eget efficitur sapien tincidunt non. Morbi fringilla velit sapien, eu facilisis risus varius in. Aliquam maximus lorem id leo lacinia euismod. Nam tristique, lacus at convallis sodales, nibh mauris feugiat leo, a consectetur sem dolor non quam. Duis non consectetur ante.
Vestibulum magna nisi, ultricies vitae erat non, congue porttitor diam. Fusce laoreet placerat accumsan. Aenean dictum neque quis mi iaculis, eu condimentum dui convallis. Curabitur eget urna lectus. Nulla in convallis quam, sit amet accumsan sapien. Proin et mi libero. Integer quis risus eu neque interdum aliquet. Nam placerat nisi vel est cursus, vel volutpat leo rutrum. Suspendisse in lectus interdum, volutpat odio sed, commodo magna. Ut venenatis nulla quis purus hendrerit sollicitudin. Nunc ultricies, enim vel auctor euismod, neque tortor finibus purus, et suscipit arcu dui at odio.`),
			key: key3,
		},
	}

	return testCases, nil
}

func newAESEngine() (*AESEngine, error) {
	return NewAESEngine()
}
