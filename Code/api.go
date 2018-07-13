package main

import (
	"net/http"
	"bytes"
	"strings"
	"fmt"
	"io"
	"SinglyLinkedList"
)

type BitSignApi struct {
	users *SinglyLinkedList.SinglyLinkedList
	documents *SinglyLinkedList.SinglyLinkedList
	signature *SinglyLinkedList.SinglyLinkedList
}

type user struct {
	username string
	password string
	publicKey string
	privateKey string
	email	 string
	name	 string
	lastName string
}

type document struct {
	hashCode string
	value string
	signaturePointer *SinglyLinkedList.Node
}

type signature struct {
	userPointer *SinglyLinkedList.Node
	signatureHash string
}

func newBitSignApi () BitSignApi {
	var v = BitSignApi{}
	v.users = SinglyLinkedList.New()
	v.documents = SinglyLinkedList.New()
	v.signature = SinglyLinkedList.New()
	return v
}

type Node struct {
	Value interface{}

	next *Node
}

func validate(api BitSignApi, username, password string) *SinglyLinkedList.Node {
	 for api.users.Head != nil {
	 	val := (api.users.Head.Value).(user)
		 if username == val.name && password == val.password {
			 return api.users.Head
		 }
	 }
	return nil
}


func (api BitSignApi) listDocuments(response http.ResponseWriter, request *http.Request) {
	user := validate(api, "yo", "man")

}
func (api BitSignApi) uploadDocuments(response http.ResponseWriter, request *http.Request) {
	var Buf bytes.Buffer
	documentToUplaod, header, err := request.FormFile("document")
	if err != nil {
		panic(err)
	}
	defer documentToUplaod.Close()
	name := strings.Split(header.Filename, ".")
	fmt.Printf("File name %s\n", name[0])
	io.Copy(&Buf, documentToUplaod)
	contents := Buf.String()
	fmt.Println(contents)
	Buf.Reset()
	//Todo upload the file to the blockchain
	return
}
func (api BitSignApi) showDocument(response http.ResponseWriter, request *http.Request) {

}
func (api BitSignApi) addSignee(response http.ResponseWriter, request *http.Request) {

}
func (api BitSignApi) removeSignee(response http.ResponseWriter, request *http.Request) {

}
func (api BitSignApi) signDocument(response http.ResponseWriter, request *http.Request) {

}