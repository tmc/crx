package crx

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/tmc/crx/ziputil"
	bufio "gopkg.in/bufio.v1"
)

const CRXMagic = "Cr24"

// CRXFromPath returns a reader that contains a CRX file.
//
// If walkable is nil it will fall back to the local filesystem.
func CRXFromPath(path string, rsaKey io.Reader, walkable ziputil.Walkable) (io.Reader, error) {
	zipBuf := bufio.NewBuffer([]byte{})
	if err := ziputil.ZipPaths(zipBuf, []string{path}, walkable); err != nil {
		return nil, err
	}
	outBuf := bufio.NewBuffer([]byte{})
	return outBuf, WriteCRXFromZip(outBuf, zipBuf, rsaKey)
}

// WriteCRXFromZip writes to the given writer a crx file that is described by zipContents.
func WriteCRXFromZip(w io.Writer, zipContents io.Reader, rsaKey io.Reader) error {
	pkey, err := privKey(rsaKey)
	if err != nil {
		return err
	}
	zipBytes, err := ioutil.ReadAll(zipContents)
	if err != nil {
		return err
	}
	sig, err := sig(bytes.NewBuffer(zipBytes), pkey)
	if err != nil {
		return err
	}
	pubBytes, err := x509.MarshalPKIXPublicKey(pkey.Public())
	if err != nil {
		return err
	}
	header := make([]byte, 16)
	copy(header, []byte(CRXMagic))
	binary.LittleEndian.PutUint32(header[4:], uint32(2))
	binary.LittleEndian.PutUint32(header[8:], uint32(len(pubBytes)))
	binary.LittleEndian.PutUint32(header[12:], uint32(len(sig)))
	fmt.Printf("%x	\n", header)
	buf := bytes.NewBuffer(header)
	if _, err := buf.Write(pubBytes); err != nil {
		return err
	}
	if _, err := buf.Write(sig); err != nil {
		return err
	}
	if _, err := buf.Write(zipBytes); err != nil {
		return err
	}
	_, err = io.Copy(w, buf)
	return err
}

func sig(r io.Reader, key *rsa.PrivateKey) ([]byte, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	sha := sha1.New()
	sha.Write(buf)
	shaBytes := sha.Sum(nil)
	return rsa.SignPKCS1v15(nil, key, crypto.SHA1, shaBytes)
}

func privKey(r io.Reader) (*rsa.PrivateKey, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(buf)
	if block == nil {
		return nil, fmt.Errorf("key: issue decoding pem block")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
