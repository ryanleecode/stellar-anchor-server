package mock

import (
	"github.com/drdgvhbh/gitserver/internal/git"
	"github.com/stretchr/testify/mock"
)

type CommitIter struct {
	mock.Mock
}

func (c *CommitIter) Next() (git.Commit, error) {
	args := c.Called()

	return nil, args.Error(1)
}

func (c *CommitIter) ForEach(fn func(git.Commit) error) error {
	args := c.Called(fn)

	return args.Error(0)
}

func (c *CommitIter) Close() {
}

type Reference struct {
	mock.Mock
}

func (r *Reference) Name() git.ReferenceName {
	args := r.Called()

	return git.ReferenceName(args.String(0))
}

func (r *Reference) Hash() git.Hash {
	args := r.Called()

	return args.Get(0).(git.Hash)
}

type Repository struct {
	mock.Mock
}

func (r *Repository) Head() (git.Reference, error) {
	args := r.Called()

	return &Reference{}, args.Error(1)
}

func (r *Repository) Log(options *git.LogOptions) (git.CommitIter, error) {
	args := r.Called(options)

	return &CommitIter{}, args.Error(1)
}

func (r *Repository) References() (git.ReferenceIter, error) {
	args := r.Called()

	return args.Get(0).(git.ReferenceIter), args.Error(1)
}

func (r *Repository) Reference(name git.ReferenceName) (git.Reference, error) {
	args := r.Called()

	return args.Get(0).(git.Reference), args.Error(1)
}

func (r *Repository) FindCommit(hash git.Hash) (git.Commit, error) {
	args := r.Called()

	return args.Get(0).(git.Commit), args.Error(1)
}

func (r *Repository) ReferenceMap() (map[string]git.References, error) {
	args := r.Called()

	return args.Get(0).(map[string]git.References), args.Error(1)
}

func (r *Repository) Diff(hash git.Hash) (git.Changes, error) {
	args := r.Called()

	return args.Get(0).(git.Changes), args.Error(1)
}

type Reader struct {
	mock.Mock
}

func (m *Reader) Open(path string) (git.Repository, error) {
	args := m.Called(path)

	return &Repository{}, args.Error(1)
}
