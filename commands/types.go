package commands

type repositoryT struct {
	remoteURL string
	name      string
}

type projectT struct {
	name       string
	repository []repositoryT
}

// Sorting

type byNameProjectT []projectT

func (a byNameProjectT) Len() int           { return len(a) }
func (a byNameProjectT) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byNameProjectT) Less(i, j int) bool { return a[i].name < a[j].name }

type byNameRepositoryT []repositoryT

func (a byNameRepositoryT) Len() int           { return len(a) }
func (a byNameRepositoryT) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byNameRepositoryT) Less(i, j int) bool { return a[i].name < a[j].name }
