package commands

import (
	"fmt"
	"time"

	"github.com/Petroviiic/GoGit/internal/core"
)

func RunMerge(repo *core.Repository, theirsBranch string) error {
	oursBranch, isDetached := repo.GetCurrentBranch()

	if isDetached {
		return fmt.Errorf("merge: unable to merge, HEAD is detached at %s ", oursBranch)
	}
	oursCommitHash := repo.GetBranchCommit(oursBranch)
	if oursCommitHash == "" {
		return fmt.Errorf("merge: no commits on branch %s", oursBranch)
	}

	theirsCommitHash := repo.GetBranchCommit(theirsBranch)
	if theirsCommitHash == "" {
		return fmt.Errorf("merge: %s - not something we can merge", theirsBranch)
	}

	if oursCommitHash == theirsCommitHash {
		fmt.Println("Already up to date.")
		return nil
	}

	baseCommitHash, err := repo.GetMergeBase(oursCommitHash, theirsCommitHash)

	if err != nil {
		return err
	}

	obj1, err := repo.LoadObject(oursCommitHash)
	if err != nil {
		return err
	}
	oursCommit := obj1.(*core.Commit)

	obj2, err := repo.LoadObject(theirsCommitHash)
	if err != nil {
		return err
	}
	theirsCommit := obj2.(*core.Commit)

	oursFiles := map[string]core.IndexEntry{}
	if err := core.GetFilesFromTreeHash(oursCommit.TreeHash, repo, "", oursFiles); err != nil {
		return err
	}

	theirsFiles := map[string]core.IndexEntry{}
	if err := core.GetFilesFromTreeHash(theirsCommit.TreeHash, repo, "", theirsFiles); err != nil {
		return err
	}

	filesToRemove := []string{}
	newHash := ""
	newFiles := make(map[string]core.IndexEntry)

	if oursCommitHash == baseCommitHash {
		//Fast-Forward
		//ours behind theirs
		if err := fastForwardMerge(oursFiles, theirsFiles, &filesToRemove); err != nil {
			return err
		}
		newHash = theirsCommitHash
		newFiles = theirsFiles

	} else {
		//Three-Way-Merge
		//take tree from ours and theirs latest commits
		//take tree from base commit

		obj3, err := repo.LoadObject(baseCommitHash)
		if err != nil {
			return err
		}
		baseCommit := obj3.(*core.Commit)

		baseFiles := map[string]core.IndexEntry{}
		if err := core.GetFilesFromTreeHash(baseCommit.TreeHash, repo, "", baseFiles); err != nil {
			return err
		}

		if err := threeWayMerge(baseFiles, oursFiles, theirsFiles, &filesToRemove, &newFiles); err != nil {
			return err
		}

		newHash, err = generateMergeCommit(repo, oursCommitHash, theirsCommitHash, oursBranch, theirsBranch, newFiles)
		if err != nil {
			return err
		}
	}

	if err := RemoveOldFiles(filesToRemove, repo); err != nil {
		return err
	}
	if err := RestoreWorkingDirectoryFiles(newFiles, oursFiles, "", repo); err != nil {
		return err
	}

	err = repo.SetBranchCommit(oursBranch, newHash)
	if err != nil {
		return err
	}

	if err := repo.SaveIndex(newFiles); err != nil {
		return err
	}

	return nil
}

func fastForwardMerge(oursFiles, theirsFiles map[string]core.IndexEntry, filesToRemove *[]string) error {
	for path, entry := range oursFiles {
		if val, ok := theirsFiles[path]; ok {
			if entry.Hash != val.Hash {
				*filesToRemove = append(*filesToRemove, path)
			}
		} else {
			*filesToRemove = append(*filesToRemove, path)
		}
	}
	return nil
}

func threeWayMerge(baseFiles, oursFiles, theirsFiles map[string]core.IndexEntry, filesToRemove *[]string, mergedFiles *map[string]core.IndexEntry) error {
	//merge these 3 into one big map
	allPaths := make(map[string]bool)

	for p := range baseFiles {
		allPaths[p] = true
	}
	for p := range oursFiles {
		allPaths[p] = true
	}
	for p := range theirsFiles {
		allPaths[p] = true
	}

	for path := range allPaths {
		baseEntry, inBase := baseFiles[path]
		oursEntry, inOurs := oursFiles[path]
		theirsEntry, inTheirs := theirsFiles[path]

		//new files
		if (inOurs != inTheirs) && !inBase { //inOurs xor inTheirs
			if inOurs {
				(*mergedFiles)[path] = oursEntry
			} else if inTheirs {
				(*mergedFiles)[path] = theirsEntry
			}
			continue
		}

		//deleting
		if inBase {
			if !inOurs && inTheirs && theirsEntry.Hash == baseEntry.Hash {
				*filesToRemove = append(*filesToRemove, path)
				continue
			}
			if inOurs && !inTheirs && oursEntry.Hash == baseEntry.Hash {
				*filesToRemove = append(*filesToRemove, path)
				continue
			}
		}

		//one side modified
		if inBase && inOurs && inTheirs {
			if baseEntry.Hash == oursEntry.Hash && baseEntry.Hash != theirsEntry.Hash {
				(*mergedFiles)[path] = theirsEntry
				continue
			}

			if baseEntry.Hash != oursEntry.Hash && baseEntry.Hash == theirsEntry.Hash {
				(*mergedFiles)[path] = oursEntry
				continue
			}
		}

		//nothing changed
		if inBase && inOurs && inTheirs && oursEntry.Hash == theirsEntry.Hash {
			(*mergedFiles)[path] = oursEntry
			continue
		}

		//TODO : dodaj git diff, tj da prikaze sta je konkretno conflictovano
		//TODO : obrisi komentare i procisti ovo

		if inBase && baseEntry.Hash != oursEntry.Hash && baseEntry.Hash != theirsEntry.Hash {
			return fmt.Errorf("merge conflict at path %s \n", path)
		}
		return fmt.Errorf("merge conflict at path %s \n", path)
	}
	return nil
}

func generateMergeCommit(repo *core.Repository, oursCommitHash, theirsCommitHash, oursBranch, theirsBranch string, mergedFiles map[string]core.IndexEntry) (string, error) {

	hierarchyRoot := core.CreateFolderHierarchy(mergedFiles)

	treeHash, err := core.CreateTreeStructure(hierarchyRoot, repo)

	if err != nil {
		return "", err
	}

	parentHashes := []string{oursCommitHash, theirsCommitHash}
	mergeCommit := core.NewCommit(
		"", //author and commiter could be extracted from global config file, which i dont have
		"",
		fmt.Sprintf("Merge branch '%s' into '%s'", theirsBranch, oursBranch),
		treeHash,
		parentHashes,
		time.Now().UTC(),
	)

	commitHash, err := mergeCommit.StoreObject(repo)
	if err != nil {
		return "", err
	}
	return commitHash, err
}
