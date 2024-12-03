package sdk

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"net/http"
	"time"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	digits  = "0123456789"
)

func generateRandomString(tpl string, length int) string {
	rand.NewSource(time.Now().UnixNano())
	result := make([]byte, length)
	for i := range result {
		result[i] = tpl[rand.Intn(len(tpl))]
	}
	return string(result)
}

func sign(r *http.Request, body []byte) {
	s := r.Header.Get(Wsign)
	broker := r.Header.Get(Wbroker)
	ts := r.Header.Get(Wts)
	nonce := r.Header.Get(Wnonce)

	tmp := make([]byte, 0, len(body)+len(s)+len(broker)+len(ts)+len(nonce))
	tmp = append(tmp, body...)
	tmp = append(tmp, broker...)
	tmp = append(tmp, ts...)
	tmp = append(tmp, nonce...)

	h := sha256.Sum256(tmp)
	ns := hex.EncodeToString(h[:])
	sb := ns[16:48]
	r.Header.Set(Wsign, sb)
}
