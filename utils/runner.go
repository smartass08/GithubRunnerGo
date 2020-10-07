package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)
var db DB
var channels map[string]chan bool


func MakeRunner(repo string, token string)error {
	var packedrunner string
	d, err := os.Open("runner-dependencies")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer d.Close()
	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ".gz" {
				packedrunner = file.Name()
				break
			}
		}
	}
	prpath, err := filepath.Abs("runner-dependencies/" + packedrunner)
	if err != nil {
		fmt.Println(err)
		return err
	}
	dir, err := os.Stat("runners/" + repo)
	if os.IsExist(err) {
		return nil
	}
	if os.IsNotExist(err) {
		errDir := os.MkdirAll("runners/"+repo, 0755)
		if errDir != nil {
			fmt.Println(err)
			return err
		}
	}
	postpath, err := filepath.Abs("runners/" + repo)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = os.Link(prpath, postpath+"/"+packedrunner)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = os.Chdir(postpath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	cmd := exec.Command("tar", "xzf", fmt.Sprintf("./%v", packedrunner))
	stdin, err := cmd.StdinPipe()
	defer stdin.Close()
	if err != nil {
		fmt.Println(err)
	}

	wg.Add(1)
	go func() {
		_, _ = io.WriteString(stdin, fmt.Sprintf("SelfHostedGo"))
	}()
	wg.Wait()

	wg.Add(1)
	go func() {
		_, _ = io.WriteString(stdin, fmt.Sprintf("\n"))
	}()
	wg.Done()

	wg.Add(1)
	go func() {
		_, _ = io.WriteString(stdin, fmt.Sprintf("\n"))
	}()
	wg.Done()
	return nil
}


func StartRunner(closeCh chan bool, repo string, token string) {
	defer wg.Done()
	_, ok := channels[repo]
	if ok{
		delete(channels, repo)
	}
	db.Access(GetDbUrl())
	db.Insert(repo, token, true)
	channels[repo] = closeCh
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
