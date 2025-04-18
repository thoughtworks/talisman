package gitrepo

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"

	"github.com/sirupsen/logrus"
)

type ReadFunc func(string) ([]byte, error)
type BatchReader interface {
	Start() error
	Read(string) ([]byte, error)
	Shutdown() error
}

type BatchGitObjectReader struct {
	repo         *GitRepo
	cmd          *exec.Cmd
	inputWriter  *bufio.Writer
	outputReader *bufio.Reader
	read         ReadFunc
}

func (bgor *BatchGitObjectReader) Start() error {
	return bgor.cmd.Start()
}

func (bgor *BatchGitObjectReader) Shutdown() error {
	return bgor.cmd.Process.Kill()
}

func (bgor *BatchGitObjectReader) Read(expr string) ([]byte, error) {
	return bgor.read(expr)
}

func newBatchGitObjectReader(root string) *BatchGitObjectReader {
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
	return &batchReader
}

func NewBatchGitHeadPathReader(root string) BatchReader {
	bgor := newBatchGitObjectReader(root)
	bgor.read = bgor.makePathReader(GIT_HEAD_PREFIX)
	return bgor
}

func NewBatchGitStagedPathReader(root string) BatchReader {
	bgor := newBatchGitObjectReader(root)
	bgor.read = bgor.makePathReader(GIT_STAGED_PREFIX)
	return bgor
}

func NewBatchGitObjectHashReader(root string) BatchReader {
	bgor := newBatchGitObjectReader(root)
	bgor.read = bgor.makeObjectHashReader()
	return bgor
}

type gitCatFileReadResult struct {
	contents []byte
	err      error
}

func (bgor *BatchGitObjectReader) makePathReader(prefix string) ReadFunc {
	pathChan := make(chan ([]byte))
	resultsChan := make(chan (gitCatFileReadResult))
	go bgor.doCatFile(pathChan, resultsChan)
	return func(path string) ([]byte, error) {
		pathChan <- []byte(prefix + ":" + path)
		result := <-resultsChan
		return result.contents, result.err
	}
}

func (bgor *BatchGitObjectReader) makeObjectHashReader() ReadFunc {
	objectHashChan := make(chan ([]byte))
	resultsChan := make(chan (gitCatFileReadResult))
	go bgor.doCatFile(objectHashChan, resultsChan)
	return func(objectHash string) ([]byte, error) {
		objectHashChan <- []byte(objectHash)
		result := <-resultsChan
		return result.contents, result.err
	}
}

func (bgor *BatchGitObjectReader) doCatFile(gitExpressionChan chan ([]byte), resultsChan chan (gitCatFileReadResult)) {
	for {
		gitExpression := <-gitExpressionChan
		gitExpression = append(gitExpression, '\n')
		//Write file-path expression to process input
		bgor.inputWriter.Write(gitExpression)
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
		logrus.Debugf("Git Batch Reader: FilePath: %v, Size:%v", gitExpression, filesize)

		//Read file contents upto filesize bytes from process output
		fileBytes := make([]byte, filesize)
		n, err := io.ReadFull(bgor.outputReader, fileBytes)
		if n != filesize || err != nil {
			logrus.Errorf("error reading exactly %v bytes of %v: %v (read %v bytes)", filesize, gitExpression, err, n)
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
