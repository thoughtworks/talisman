package talismanrc

// Mapping of language scopes to names of files that should be ignored anywhere in a repository
var knownScopes = map[string][]string{
	"node":      {"pnpm-lock.yaml", "yarn.lock", "package-lock.json"},
	"go":        {"makefile", "go.mod", "go.sum", "Gopkg.toml", "Gopkg.lock", "glide.yaml", "glide.lock"},
	"images":    {"*.jpeg", "*.jpg", "*.png", "*.tiff", "*.bmp"},
	"bazel":     {"*.bzl"},
	"terraform": {".terraform.lock.hcl"},
	"php":       {"composer.lock"},
	"python":    {"poetry.lock", "Pipfile.lock", "requirements.txt"},
}
