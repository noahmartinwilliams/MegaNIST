package main 

import "io/ioutil"
import "sync"

func FindFiles(dirname string, doPanic bool) chan string {
	retc := make(chan string, 1024)
	go func() {
		defer close(retc)
		var wg sync.WaitGroup
		wg.Add(1)
		FindFilesIntern(&wg, retc, dirname, doPanic)
		wg.Wait()
	} ()
	return retc
}

func FindFilesIntern(wg *sync.WaitGroup, retc chan string, dirname string, doPanic bool) {
	go func() {
		defer wg.Done()
		files, err := ioutil.ReadDir(dirname)
		if err != nil && doPanic {
			panic(err)
		} else if err != nil && !doPanic {
			return 
		}

		for _, file := range files {
			if !file.IsDir() {
				retc <- dirname + "/" + file.Name()
			} else {
				wg.Add(1)
				FindFilesIntern(wg, retc, dirname + "/" + file.Name(), doPanic)
			}
		}
	} ()
}
