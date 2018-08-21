package balloon

import (
	"github.com/iknite/molletree/balloon/history"
	"github.com/iknite/molletree/balloon/hyper"
)

type Balloon struct {
	history *history.Tree
	hyper   *hyper.Tree
	version uint64
}

func NewBalloon() *Balloon {
	return &Balloon{
		history.NewTree(),
		hyper.NewTree(),
		0,
	}
}

func (b *Balloon) Add(message string) (historyCommitment, hyperCommitment, digest []byte) {
	// historyChan := make(chan []byte)
	// hyperChan := make(chan []byte)
	//
	// digest = b.history.hash.Do(encstring.ToBytes(message))
	//
	// go func() { historyChan <- b.history.Add(digest) }()
	// go func() { hyperChan <- b.hyper.Add(digest, b.version) }()
	//
	// historyCommitment = <-historyChan
	// hyperCommitment = <-hyperChan
	//
	return
}
