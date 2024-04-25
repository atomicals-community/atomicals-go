package merkle

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/atomicals-core/atomicals/witness"
)

type MerkleNode struct {
	Right []byte
	Left  []byte
}

func CheckValidateProof(expected_root_hash, target_hash []byte, proof []witness.Proof) bool {
	formatted_proof := make([]*MerkleNode, 0)
	for _, item := range proof {
		if item.P {
			formatted_proof = append(formatted_proof, &MerkleNode{
				Right: item.D,
			})
		} else {
			formatted_proof = append(formatted_proof, &MerkleNode{
				Left: item.D,
			})
		}
	}
	return ValidateProof(formatted_proof, target_hash, expected_root_hash)
}

func ValidateProof(proof []*MerkleNode, targetHash, expectedRootHash []byte) bool {
	computedHash := targetHash
	for _, node := range proof {
		if node.Left != nil && node.Right != nil {
			return false // Both left and right hashes are provided, which is invalid
		}
		if node.Left != nil {
			computedHash = ComputeParentHash(node.Left, computedHash)
		} else {
			computedHash = ComputeParentHash(computedHash, node.Right)
		}
	}
	return hex.EncodeToString(computedHash) == hex.EncodeToString(expectedRootHash)
}

func ComputeParentHash(leftHash, rightHash []byte) []byte {
	// Concatenate the left and right hashes and hash the result using SHA-256
	combinedHash := append(leftHash, rightHash...)
	hash := sha256.New()
	hash.Write([]byte(combinedHash))
	return hash.Sum(nil)
}
