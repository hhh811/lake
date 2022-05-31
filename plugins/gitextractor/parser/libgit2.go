/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package parser

import (
	"context"
	"fmt"

	git "github.com/libgit2/git2go/v33"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitextractor/models"
)

const (
	BRANCH = "BRANCH"
	TAG    = "TAG"
)

type LibGit2 struct {
	store      models.Store
	logger     core.Logger
	ctx        context.Context     // for canceling
	subTaskCtx core.SubTaskContext // for updating progress
}

func NewLibGit2(store models.Store, subTaskCtx core.SubTaskContext) *LibGit2 {
	return &LibGit2{store: store,
		logger:     subTaskCtx.GetLogger(),
		ctx:        subTaskCtx.GetContext(),
		subTaskCtx: subTaskCtx}
}

func (l *LibGit2) LocalRepo(repoPath, repoId string) error {
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return err
	}
	return l.run(repo, repoId)
}

func (l *LibGit2) run(repo *git.Repository, repoId string) error {
	defer l.store.Close()
	l.subTaskCtx.SetProgress(0, -1)
	err := l.collectTags(repo, repoId)
	if err != nil {
		return err
	}
	err = l.collectBranches(repo, repoId)
	if err != nil {
		return err
	}
	return l.collectCommits(repo, repoId)
}

func (l *LibGit2) collectTags(repo *git.Repository, repoId string) error {
	return repo.Tags.Foreach(func(name string, id *git.Oid) error {
		select {
		case <-l.ctx.Done():
			return l.ctx.Err()
		default:
			break
		}
		var err1 error
		var tag *git.Tag
		var tagCommit string
		tag, _ = repo.LookupTag(id)
		if tag != nil {
			tagCommit = tag.TargetId().String()
		} else {
			tagCommit = id.String()
		}
		l.logger.Info("tagCommit", tagCommit)
		if tagCommit != "" {
			ref := &code.Ref{
				DomainEntity: domainlayer.DomainEntity{Id: fmt.Sprintf("%s:%s", repoId, name)},
				RepoId:       repoId,
				Name:         name,
				CommitSha:    tagCommit,
				RefType:      TAG,
			}
			err1 = l.store.Refs(ref)
			if err1 != nil {
				return err1
			}
			l.subTaskCtx.IncProgress(1)
		}
		return nil
	})
}

func (l *LibGit2) collectBranches(repo *git.Repository, repoId string) error {
	var repoInter *git.BranchIterator
	repoInter, err := repo.NewBranchIterator(git.BranchAll)
	if err != nil {
		return err
	}
	return repoInter.ForEach(func(branch *git.Branch, branchType git.BranchType) error {
		select {
		case <-l.ctx.Done():
			return l.ctx.Err()
		default:
			break
		}
		if branch.IsBranch() {
			name, err1 := branch.Name()
			if err1 != nil {
				return err1
			}
			var sha string
			if oid := branch.Target(); oid != nil {
				sha = oid.String()
			}
			ref := &code.Ref{
				DomainEntity: domainlayer.DomainEntity{Id: fmt.Sprintf("%s:%s", repoId, name)},
				RepoId:       repoId,
				Name:         name,
				CommitSha:    sha,
				RefType:      BRANCH,
			}
			ref.IsDefault, _ = branch.IsHead()
			err1 = l.store.Refs(ref)
			if err1 != nil {
				return err1
			}
			l.subTaskCtx.IncProgress(1)
			return nil
		}
		return nil
	})
}

func (l *LibGit2) collectCommits(repo *git.Repository, repoId string) error {
	opts, err := getDiffOpts()
	if err != nil {
		return err
	}
	odb, err := repo.Odb()
	if err != nil {
		return err
	}
	return odb.ForEach(func(id *git.Oid) error {
		select {
		case <-l.ctx.Done():
			return l.ctx.Err()
		default:
			break
		}
		commit, _ := repo.LookupCommit(id)
		if commit == nil {
			return nil
		}
		commitSha := commit.Id().String()
		l.logger.Debug("process commit: %s", commitSha)
		c := &code.Commit{
			Sha:     commitSha,
			Message: commit.Message(),
		}
		author := commit.Author()
		if author != nil {
			c.AuthorName = author.Name
			c.AuthorEmail = author.Email
			c.AuthorId = author.Email
			c.AuthoredDate = author.When
		}
		committer := commit.Committer()
		if committer != nil {
			c.CommitterName = committer.Name
			c.CommitterEmail = committer.Email
			c.CommitterId = committer.Email
			c.CommittedDate = committer.When
		}
		if err != l.storeParentCommits(commitSha, commit) {
			return err
		}
		if commit.ParentCount() > 0 {
			parent := commit.Parent(0)
			if parent != nil {
				var stats *git.DiffStats
				if stats, err = l.getDiffComparedToParent(c.Sha, commit, parent, repo, opts); err != nil {
					return err
				} else {
					c.Additions += stats.Insertions()
					c.Deletions += stats.Deletions()
				}
			}
		}
		err = l.store.Commits(c)
		if err != nil {
			return err
		}
		repoCommit := &code.RepoCommit{
			RepoId:    repoId,
			CommitSha: c.Sha,
		}
		err = l.store.RepoCommits(repoCommit)
		if err != nil {
			return err
		}
		l.subTaskCtx.IncProgress(1)
		return nil
	})
}

func (l *LibGit2) storeParentCommits(commitSha string, commit *git.Commit) error {
	var commitParents []*code.CommitParent
	for i := uint(0); i < commit.ParentCount(); i++ {
		parent := commit.Parent(i)
		if parent != nil {
			if parentId := parent.Id(); parentId != nil {
				commitParents = append(commitParents, &code.CommitParent{
					CommitSha:       commitSha,
					ParentCommitSha: parentId.String(),
				})
			}
		}
	}
	return l.store.CommitParents(commitParents)
}

func (l *LibGit2) getDiffComparedToParent(commitSha string, commit *git.Commit, parent *git.Commit, repo *git.Repository, opts *git.DiffOptions) (*git.DiffStats, error) {
	var err error
	var parentTree, tree *git.Tree
	parentTree, err = parent.Tree()
	if err != nil {
		return nil, err
	}
	tree, err = commit.Tree()
	if err != nil {
		return nil, err
	}
	var diff *git.Diff
	diff, err = repo.DiffTreeToTree(parentTree, tree, opts)
	if err != nil {
		return nil, err
	}
	var commitFile *code.CommitFile
	commitFile, err = l.generateCommitFileFromDiff(commitSha, diff)
	if err != nil {
		return nil, err
	}
	if commitFile != nil {
		err = l.store.CommitFiles(commitFile)
		if err != nil {
			l.logger.Error("CommitFiles error:", err)
		}
	}
	var stats *git.DiffStats
	stats, err = diff.Stats()
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (l *LibGit2) generateCommitFileFromDiff(commitSha string, diff *git.Diff) (*code.CommitFile, error) {
	var commitFile *code.CommitFile
	var err error
	err = diff.ForEach(func(file git.DiffDelta, progress float64) (
		git.DiffForEachHunkCallback, error) {
		if commitFile != nil {
			err = l.store.CommitFiles(commitFile)
			if err != nil {
				l.logger.Error("CommitFiles error:", err)
				return nil, err
			}
		}
		commitFile = new(code.CommitFile)
		commitFile.CommitSha = commitSha
		commitFile.FilePath = file.NewFile.Path
		return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
			return func(line git.DiffLine) error {
				if line.Origin == git.DiffLineAddition {
					commitFile.Additions += line.NumLines
				}
				if line.Origin == git.DiffLineDeletion {
					commitFile.Deletions += line.NumLines
				}
				return nil
			}, nil
		}, nil
	}, git.DiffDetailLines)
	return commitFile, err
}

func getDiffOpts() (*git.DiffOptions, error) {
	opts, err := git.DefaultDiffOptions()
	if err != nil {
		return nil, err
	}
	opts.NotifyCallback = func(diffSoFar *git.Diff, delta git.DiffDelta, matchedPathSpec string) error {
		return nil
	}
	return &opts, nil
}
