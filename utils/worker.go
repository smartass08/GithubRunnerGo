package utils

import (
	"fmt"
	"os"
	"os/exec"

)
var db DB
var channels map[string]chan bool

func StartRunner(closeCh chan bool, repo string, token string) {
	defer wg.Done()
	_, ok := channels[repo]
	if ok{
		delete(channels, repo)
	}
	db.Access(GetDbUrl())
	db.Insert(repo, token, true)
	channels[repo] = closeCh
	_ = os.Chdir("actions-runner")
	exec.Command("./config.sh", "--url", fmt.Sprintf("%v", repo), "--token", fmt.Sprintf("%v", token))
	for {
		stop := <-closeCh
		if stop{
			delete(channels, repo)
			db.Delete(repo, token, true)
			return
		}
	}
}

func StopRunner(repo string)  {
	_,ok := channels[repo]
	if ok{
		stop := channels[repo]
		stop <- true
		delete(channels, repo)
	}
}

func GetAllChannels()map[string]chan bool{
	return channels
}

func CheckAll(){
	for _,v := range ALL{
		if v.Running{
			ch := make(chan bool)
			ch <- false
			StartRunner(ch, v.Repo, v.Token)
		}
	}
}