package hash

import "golang.org/x/crypto/bcrypt"

func Generate(s string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		// I'm panicking here because I don't think there's any way for this to
		// error that shouldn't immediately cause a page. I've never seen it happen
		// and I think it's just for interface satisfaction, so I feel safe here.
		// Plus, it's not a recoverable error. The user can't fix a broken hashing
		// algorithm.
		panic(err)
	}
	return string(hash)
}
