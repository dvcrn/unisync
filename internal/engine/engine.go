package engine

type Engine interface {
	Sync(pathA, pathB string)
	SyncAToB(pathA, pathB string)
	SyncBToA(pathA, pathB string)
}
