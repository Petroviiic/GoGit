package commands

import (
	"fmt"

	"github.com/Petroviiic/GoGit/internal/core"
)

func RunMerge(repo *core.Repository, theirsBranch string) error {
	oursBranch := repo.GetCurrentBranch()
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
	fmt.Println(baseCommitHash)

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

	//Fast-Forward
	//ours behind theirs
	//refs/heads/ourbranch.setbranchcommit(theirs)
	if oursCommitHash == baseCommitHash {
		if err := repo.SetBranchCommit(oursBranch, theirsCommitHash); err != nil {
			return err
		}

		filesToRemove := []string{}
		for path, entry := range oursFiles {
			if val, ok := theirsFiles[path]; ok {
				if entry.Hash != val.Hash {
					filesToRemove = append(filesToRemove, path)
				}
			} else {
				filesToRemove = append(filesToRemove, path)
			}
		}

		if err := RemoveOldFiles(filesToRemove, repo); err != nil {
			return err
		}

		if err := RestoreWorkingDirectoryFiles(theirsFiles, oursFiles, "", repo); err != nil {
			return err
		}

		if err := repo.SaveIndex(theirsFiles); err != nil {
			return err
		}

		return nil
	}

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

	mergedFiles := make(map[string]core.IndexEntry)
	filesToRemove := []string{}

	for path := range allPaths {
		baseEntry, inBase := baseFiles[path]
		oursEntry, inOurs := oursFiles[path]
		theirsEntry, inTheirs := theirsFiles[path]

		if (inOurs != inTheirs) && !inBase { //inOurs xor inTheirs
			//new file
			if inOurs {
				mergedFiles[path] = oursEntry
			} else if inTheirs {
				mergedFiles[path] = theirsEntry
			}
			continue

		}
		if inBase && inOurs && !inTheirs && baseEntry.Hash == oursEntry.Hash {
			//deleted in theirs
			filesToRemove = append(filesToRemove, path)
			continue
		}
		if inBase && !inOurs && inTheirs && baseEntry.Hash == theirsEntry.Hash {
			//deleted in ours
			mergedFiles[path] = theirsEntry
			continue
		}

		if inBase && inOurs && inTheirs && baseEntry.Hash == oursEntry.Hash && baseEntry.Hash != theirsEntry.Hash {
			mergedFiles[path] = theirsEntry
			continue

		}
		if inBase && inOurs && inTheirs && baseEntry.Hash != oursEntry.Hash && baseEntry.Hash == theirsEntry.Hash {
			mergedFiles[path] = oursEntry
			continue

		}

		if inBase && inOurs && inTheirs && oursEntry.Hash == theirsEntry.Hash {
			mergedFiles[path] = oursEntry
			continue
		}

		if baseEntry.Hash != oursEntry.Hash && baseEntry.Hash != theirsEntry.Hash {
			return fmt.Errorf("merge conflict at path %s \n", path)
		}
	}
	//loop through each file; maybe add ok files to a list or dict
	//	if unique in ours/theirs and not present in base => new, ok
	//	if present in base and ours and file.hash(base) == file.hash(ours) and not present in theirs => deleted, ok/skip
	//	if file.hash(base) == file.hash(ours) && file.hash(base) != file.hash(theirs); changed only in ours; ok
	//	if file.hash(base) == file.hash(theirs) && file.hash(base) != file.hash(ours); changed only in theirs; ok

	//	if file.hash(ours) == file.hash(theirs); ok
	//		else
	// 			if file.hash(base) != file.hash(ours) && file.hash(base) != file.hash(theirs); conflict

	//	if !conflict;
	// 		make merge commit;
	//		recreate new directory,
	// 		and tree and add ours and theirscommits as parents

	// obj, err := repo.LoadObject("35ff57b34823f848b33bf6544828fb150b811663")
	// theirsCommit := obj.(*core.Commit)
	// fmt.Println(theirsCommit)
	return nil
}

func SyncWorkingDirectory(oursFiles, theirsFiles, allFiles map[string]string) {
	// for path, hash := range allFiles{

	// }
}
