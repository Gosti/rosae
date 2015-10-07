package main

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/mrgosti/rosa"
	"os"
	"os/user"
	"sort"
	"strings"
)

type Command func() error

var _PrivateKey *rsa.PrivateKey
var me *user.User

func generate() error {
	if len(os.Args) < 3 {
		return errors.New("Not enought parameter to generate a key\n Rosae generate usage => rosae generate random text")
	}
	_, _, err := rosa.Generate(strings.Join(os.Args[2:], " "), true)
	return err
}

func friends() error {
	var fd []string
	for _, v := range rosa.FriendList {
		fd = append(fd, v.Name)
	}
	sort.Strings(fd)
	for _, v := range fd {
		fmt.Println(v)
	}
	return nil
}

func friend() error {
	if len(os.Args) < 3 {
		return errors.New("Not enought parameter to find a friend\n Rosae friend usage => rosae friend friend_name")
	}
	him := rosa.SeekByName(os.Args[2])
	if him == nil {
		return errors.New("Friend not found")
	}
	fmt.Println(rosa.StringifyPublicKey(him.PublicKey))
	return nil
}

func add() error {
	if len(os.Args) < 4 {
		return errors.New("Not enought parameter to add a friend\n Rosae friend usage => rosae add friend_name encoded_key")
	}
	key, err := rosa.UnStringifyPublicKey(os.Args[3])
	if err != nil {
		return err
	}
	him := &rosa.Friend{os.Args[2], key}
	return him.Registrer(me.HomeDir + "/.rosa/friend_list")
}

func deleteFriend() error {
	if len(os.Args) < 3 {
		return errors.New("Not enought parameter to delete a friend\n Rosae delete usage => rosae delete friend_name")
	}
	him := rosa.SeekByName(os.Args[2])
	if him == nil {
		return errors.New("Friend not found\n")
	}
	return him.Delete(me.HomeDir + "/.rosa/friend_list")
}

func public() error {
	fmt.Println(rosa.StringifyPublicKey(&_PrivateKey.PublicKey))
	return nil
}

func help() error {
	fmt.Println("Rosae usage :")
	fmt.Println("\trosae generate some word ... : generate a key based on the word you use as parameters")
	fmt.Println("\trosae friends : give the name of all your friends")
	fmt.Println("\trosae friend friend_name : give encoded public key of a given friend")
	fmt.Println("\trosae public : give your public key")
	fmt.Println("\trosae add friend_name encoded_key: add a friend to your friends list")
	fmt.Println("\trosae delete friend_name: delete a friend from your friends list")
	fmt.Println("\trosae help : duh.")
	return nil
}

func checkErr(err error) {
	if err != nil {
		fmt.Printf("Rosae error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	var err error

	me, err = user.Current()
	checkErr(err)

	_PrivateKey, err = rosa.LoadPrivateKey(me.HomeDir + "/.rosa/key.priv")
	checkErr(err)

	err = rosa.LoadFriends(me.HomeDir + "/.rosa/friend_list")
	checkErr(err)

	var cmdList map[string]Command

	cmdList = make(map[string]Command)

	cmdList["friends"] = friends
	cmdList["friend"] = friend
	cmdList["add"] = add
	cmdList["delete"] = deleteFriend
	cmdList["generate"] = generate
	cmdList["public"] = public
	cmdList["help"] = help

	if len(os.Args) < 2 {
		checkErr(cmdList["help"]())
		os.Exit(0)
	}

	checkErr(cmdList[os.Args[1]]())
}
