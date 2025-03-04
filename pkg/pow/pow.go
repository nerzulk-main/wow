package pow

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
)

func Solve(ctx context.Context, challenge []byte, difficulty int) (string, error) {
	prefix := strings.Repeat("0", difficulty)
	for i := 0; i < math.MaxInt64; i++ {
		if i%10000 == 0 {
			if err := ctx.Err(); err != nil {
				return "", err
			}
		}
		solution := fmt.Sprintf("%x", i)
		sum := sha256.Sum256(append(challenge, solution...))
		hashHex := hex.EncodeToString(sum[:])
		if strings.HasPrefix(hashHex, prefix) {
			return solution, nil
		}
	}

	return "", fmt.Errorf("solution not found")
}

func ValidateChallenge(complexity int, data []byte, result string) bool {
	sum := sha256.Sum256(append(data, result...))
	hashHex := hex.EncodeToString(sum[:])
	return strings.HasPrefix(hashHex, strings.Repeat("0", complexity))
}
