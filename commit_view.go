package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
)

type ViewIndex struct {
	activeIndex    uint
	viewStartIndex uint
}

type CommitView struct {
	repoData     RepoData
	activeBranch *Oid
	active       bool
	viewIndex    map[*Oid]*ViewIndex
}

func NewCommitView(repoData RepoData) *CommitView {
	return &CommitView{
		repoData:  repoData,
		viewIndex: make(map[*Oid]*ViewIndex),
	}
}

func (commitView *CommitView) Initialise() (err error) {
	log.Info("Initialising CommitView")
	return
}

func (commitView *CommitView) Render(win RenderWindow) (err error) {
	log.Debug("Rendering CommitView")

	var viewIndex *ViewIndex
	var ok bool
	if viewIndex, ok = commitView.viewIndex[commitView.activeBranch]; !ok {
		return fmt.Errorf("No ViewIndex exists for oid %v", commitView.activeBranch)
	}

	commits := commitView.repoData.Commits(commitView.activeBranch)
	rows := win.Rows() - 2
	rowDiff := viewIndex.activeIndex - viewIndex.viewStartIndex

	if rowDiff < 0 {
		viewIndex.viewStartIndex = viewIndex.activeIndex
	} else if rowDiff >= rows {
		viewIndex.viewStartIndex += (rowDiff - rows) + 1
	}

	commitIndex := viewIndex.viewStartIndex

	for rowIndex := uint(0); rowIndex < rows && commitIndex < uint(len(commits)); rowIndex++ {
		commit := commits[commitIndex]
		author := commit.commit.Author()

		if err = win.SetRow(rowIndex+1, " %v %s %s", author.When, author.Name, commit.commit.Summary()); err != nil {
			break
		}

		commitIndex++
	}

	if err = win.SetSelectedRow((viewIndex.activeIndex-viewIndex.viewStartIndex)+1, commitView.active); err != nil {
		return
	}

	win.DrawBorder()

	return err
}

func (commitView *CommitView) OnRefSelect(oid *Oid) (err error) {
	log.Debugf("CommitView loading commits for selected oid %v", oid)

	if _, ok := commitView.viewIndex[oid]; ok {
		return
	}

	if err = commitView.repoData.LoadCommits(oid); err != nil {
		return
	}

	commitView.activeBranch = oid
	commitView.viewIndex[oid] = &ViewIndex{}
	return
}

func (commitView *CommitView) Handle(keyPressEvent KeyPressEvent, channels HandlerChannels) (err error) {
	log.Debugf("CommitView handling key %v", keyPressEvent)
	return
}

func (commitView *CommitView) OnActiveChange(active bool) {
	log.Debugf("CommitView active %v", active)
	commitView.active = active
}
