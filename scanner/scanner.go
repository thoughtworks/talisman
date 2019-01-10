package scanner

import (
	"os/exec"
	"strings"
	"talisman/git_repo"
)

func GetAdditions(blob_details string) []git_repo.Addition {
	var additions []git_repo.Addition
	blob_array := strings.Split(blob_details, "\n")
	for _, blob := range blob_array {
		object_details := strings.Split(blob, " ")
		object_hash := object_details[0]
		data := get_data(object_hash)
		file_path := object_details[1]
		new_addition := git_repo.NewAddition(file_path, data)
		additions = append(additions, new_addition)

	}
	return additions
}

func get_data(object_hash string) []byte {
	out, _ := exec.Command("git", "cat-file", "-p", object_hash).CombinedOutput()
	return out
}
