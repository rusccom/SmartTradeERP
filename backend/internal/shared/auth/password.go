package auth

import (
    "crypto/rand"
    "crypto/subtle"
    "encoding/base64"
    "fmt"
    "strconv"
    "strings"

    "golang.org/x/crypto/argon2"
)

type hashParams struct {
    memory      uint32
    iterations  uint32
    parallelism uint8
}

func HashPassword(password string) (string, error) {
    params := defaultParams()
    salt, err := randomBytes(16)
    if err != nil {
        return "", err
    }
    hash := deriveKey(password, salt, params, 32)
    return encodeHash(params, salt, hash), nil
}

func VerifyPassword(password, encoded string) bool {
    params, salt, expected, err := decodeHash(encoded)
    if err != nil {
        return false
    }
    actual := deriveKey(password, salt, params, uint32(len(expected)))
    return subtle.ConstantTimeCompare(actual, expected) == 1
}

func defaultParams() hashParams {
    return hashParams{memory: 64 * 1024, iterations: 1, parallelism: 4}
}

func deriveKey(password string, salt []byte, params hashParams, keyLength uint32) []byte {
    raw := []byte(password)
    return argon2.IDKey(raw, salt, params.iterations, params.memory, params.parallelism, keyLength)
}

func encodeHash(params hashParams, salt, hash []byte) string {
    saltB64 := base64.RawStdEncoding.EncodeToString(salt)
    hashB64 := base64.RawStdEncoding.EncodeToString(hash)
    return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
        params.memory,
        params.iterations,
        params.parallelism,
        saltB64,
        hashB64,
    )
}

func decodeHash(encoded string) (hashParams, []byte, []byte, error) {
    parts := strings.Split(encoded, "$")
    if len(parts) != 6 {
        return hashParams{}, nil, nil, fmt.Errorf("invalid hash format")
    }
    params, err := parseParams(parts[3])
    if err != nil {
        return hashParams{}, nil, nil, err
    }
    salt, err := base64.RawStdEncoding.DecodeString(parts[4])
    if err != nil {
        return hashParams{}, nil, nil, err
    }
    hash, err := base64.RawStdEncoding.DecodeString(parts[5])
    if err != nil {
        return hashParams{}, nil, nil, err
    }
    return params, salt, hash, nil
}

func parseParams(raw string) (hashParams, error) {
    params := hashParams{}
    for _, field := range strings.Split(raw, ",") {
        kv := strings.SplitN(field, "=", 2)
        if len(kv) != 2 {
            return hashParams{}, fmt.Errorf("invalid argon2 params")
        }
        if err := assignParam(&params, kv[0], kv[1]); err != nil {
            return hashParams{}, err
        }
    }
    return params, nil
}

func assignParam(params *hashParams, key, value string) error {
    parsed, err := strconv.Atoi(value)
    if err != nil {
        return err
    }
    switch key {
    case "m":
        params.memory = uint32(parsed)
    case "t":
        params.iterations = uint32(parsed)
    case "p":
        params.parallelism = uint8(parsed)
    default:
        return fmt.Errorf("unknown argon2 param")
    }
    return nil
}

func randomBytes(size int) ([]byte, error) {
    bytes := make([]byte, size)
    _, err := rand.Read(bytes)
    if err != nil {
        return nil, err
    }
    return bytes, nil
}
