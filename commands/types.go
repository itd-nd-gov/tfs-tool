package commands

type repository struct {
	remoteURL string
	name      string
}

type repositories []repository

func (a repositories) Len() int      { return len(a) }
func (a repositories) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type repositoriesByName struct{ repositories }

func (a repositoriesByName) Less(i, j int) bool {
	return a.repositories[i].name < a.repositories[j].name
}

//

type project struct {
	name string
	repositories
}

type projects []project

func (a projects) Len() int      { return len(a) }
func (a projects) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type projectsByName struct{ projects }

func (a projectsByName) Less(i, j int) bool { return a.projects[i].name < a.projects[j].name }
