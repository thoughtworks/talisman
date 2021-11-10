package gitrepo

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Reader func(string) ([]byte, error)

func NewCommittedRepoFileReader(wd string) Reader {
	return makeReader(wd, GIT_HEAD_PREFIX)
}

func NewBatchGitObjectReader(root string) Reader {
	repo := &GitRepo{root}
	cmd := repo.makeRepoCommand("git", "cat-file", "--batch=%(objectsize)")
	inputPipe, err := cmd.StdinPipe()
	if err != nil {
		logrus.Fatalf("error creating stdin pipe for batch git file reader subprocess: %v", err)
	}
	outputPipe, err := cmd.StdoutPipe()
	if err != nil {
		logrus.Fatalf("error creating stdout pipe for batch git file reader subprocess: %v", err)
	}
	batchReader := BatchGitObjectReader{
		repo:         repo,
		cmd:          cmd,
		inputWriter:  bufio.NewWriter(inputPipe),
		outputReader: bufio.NewReader(outputPipe),
	}
	err = batchReader.start()
	if err != nil {
		logrus.Fatalf("error starting batch command: %v", err)
	}
	return batchReader.makeReader()
}

func NewRepoFileReader(wd string) Reader {
	return makeReader(wd, GIT_STAGED_PREFIX)
}

func makeReader(repoRoot, prefix string) Reader {
	return func(path string) ([]byte, error) {
		return GitRepo{repoRoot}.readRepoFile(path, prefix)
	}
}

type BatchGitObjectReader struct {
	repo         *GitRepo
	cmd          *exec.Cmd
	inputWriter  *bufio.Writer
	outputReader *bufio.Reader
}

func (bgor *BatchGitObjectReader) start() error {
	return bgor.cmd.Start()
}

type gitCatFileReadResult struct {
	contents []byte
	err      error
}

func (bgor *BatchGitObjectReader) makeReader() Reader {
	pathChan := make(chan (string))
	resultsChan := make(chan (gitCatFileReadResult))
	go bgor.doCatFile(pathChan, resultsChan)
	return func(path string) ([]byte, error) {
		pathChan <- path
		result := <-resultsChan
		return result.contents, result.err
	}
}

func (bgor *BatchGitObjectReader) doCatFile(pathChan chan (string), resultsChan chan (gitCatFileReadResult)) {
	for {
		path := <-pathChan
		//Write file-path expression to process input
		bgor.inputWriter.Write([]byte(":" + path + "\n"))
		bgor.inputWriter.Flush()

		//Read line containing filesize from process output
		filesizeBytes, err := bgor.outputReader.ReadBytes('\n')
		if err != nil {
			logrus.Errorf("error reading filesize: %v", err)
			resultsChan <- gitCatFileReadResult{[]byte{}, err}
			continue
		}
		filesize, err := strconv.Atoi(string(filesizeBytes[:len(filesizeBytes)-1]))
		if err != nil {
			logrus.Errorf("error parsing filesize: %v", err)
			resultsChan <- gitCatFileReadResult{[]byte{}, err}
			continue
		}
		logrus.Debugf("Git Batch Reader: FilePath: %v, Size:%v", path, filesize)

		//Read file contents upto filesize bytes from process output
		fileBytes := make([]byte, filesize)
		n, err := io.ReadFull(bgor.outputReader, fileBytes)
		if n != filesize || err != nil {
			logrus.Errorf("error reading exactly %v bytes of %v: %v (read %v bytes)", filesize, path, err, n)
			resultsChan <- gitCatFileReadResult{[]byte{}, err}
			continue
		}

		//Read and discard trailing newline (not doing this will cause errors going forward)
		b, err := bgor.outputReader.ReadByte()
		if err != nil || b != '\n' {
			resultsChan <- gitCatFileReadResult{[]byte{}, fmt.Errorf("error discarding trailing newline : %v trailing byte: %v", err, b)}
			continue
		}

		resultsChan <- gitCatFileReadResult{fileBytes, nil}
	}
}
