package domain

import "github.com/alexedwards/argon2id"

func Matches(plainTextPassword, hashedPassword string) (bool, error) {
	return argon2id.ComparePasswordAndHash(plainTextPassword, hashedPassword)
}

func Hash(plainTextPassword string) (string, error) {
	params := &argon2id.Params{
		Memory:      12 * 1024,
		Iterations:  3,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	}

	hash, err := argon2id.CreateHash(plainTextPassword, params)
	if err != nil {
		return "", err
	}

	return hash, nil
}
