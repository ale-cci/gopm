package chunk

type ChunkIterator struct {
	Files    []Iterator
	currFile int
}

func (ci *ChunkIterator) getCurrentFile() Iterator {
	ci.currFile = ci.currFile % len(ci.Files)

	if ci.currFile < 0 {
		ci.currFile += len(ci.Files)
	}

	return ci.Files[ci.currFile]
}

func (ci *ChunkIterator) Current() string {

	file := ci.getCurrentFile()
	return file.Current()
}

func (ci *ChunkIterator) Next() bool {
	file := ci.getCurrentFile()
	if file.Next() {
		ci.currFile++
	}
	return false
}

func (ci *ChunkIterator) Prev() bool {
	file := ci.getCurrentFile()

	if file.Prev() {
		ci.currFile -= 1
	}
	return false
}
