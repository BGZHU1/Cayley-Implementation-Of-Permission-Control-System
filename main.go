package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/kv/bolt"
	"github.com/cayleygraph/cayley/schema"
	"github.com/cayleygraph/quad"
	"github.com/cayleygraph/quad/voc"

	// Import RDF vocabulary definitions to be able to expand IRIs like rdf:label.
	_ "github.com/cayleygraph/quad/voc/core"
)

//define the rdf:types of objects
//AccessType : Private Public
type AccessType struct {
	rdfType struct{} `quad:"@type > ex:AccessType"`
	ID      quad.IRI `json:"@id"`
}

//define the role - user/admin
type Role struct {
	rdfType   struct{}   `quad:"@type > ex:Role"`
	ID        quad.IRI   `json:"@id"`
	HasAction []quad.IRI `json:"ex:hasAction"` // field name (predicate) may be written as json field name
}

//define the Action - read/write
type Action struct {
	rdfType struct{} `quad:"@type > ex:Action"`
	ID      quad.IRI `json:"@id"`
}

//define Document
type Document struct {
	rdfType            struct{} `quad:"@type > ex:Document"`
	ID                 quad.IRI `json:"@id"` //name of the document
	HasAuthorizedAgent quad.IRI `json:"ex:hasAuthorizedAgent"`
	Creator            quad.IRI `json:"ex:creator"`
	HasAccessType      quad.IRI `json:"ex:hasAccessType"`
}

//define Agent
type Agent struct {
	rdfType                       struct{}   `quad:"@type > ex:Agent"`
	ID                            quad.IRI   `json:"@id"` //name of the document
	HasRole                       quad.IRI   `json:"ex:hasRole"`
	HasAuthorizedActionOnResource []quad.IRI `json:"ex:hasAuthorizedActionOnResource"`
}

//this is the subclass of action

type AuthorizedActionOnResource struct {
	rdfType             struct{}      `quad:"@type > ex:AuthorizedActionOnResource"`
	ID                  quad.IRI      `json:"@id"` //name of the action, which is a subclass of action
	HasResource         quad.IRI      `json:"ex:hasResource"`
	HasActionOnResource []quad.IRI    `json:"ex:hasActionOnResource"`
	SuperClasses        []interface{} `json:"rdfs:subClassOf"`
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	voc.RegisterPrefix("ex:", "http://coordy.org/")

	sch := schema.NewConfig()
	// Override a function to generate IDs. Can be changed to generate UUIDs, for example.
	sch.GenerateID = func(_ interface{}) quad.Value {
		return quad.BNode(fmt.Sprintf("node%d", rand.Intn(1000)))
	}

	// File for your new BoltDB. Use path to regular file and not temporary in the real world
	tmpdir, err := ioutil.TempDir("", "example")
	checkErr(err)

	defer os.RemoveAll(tmpdir) // clean up

	// Initialize the database
	err = graph.InitQuadStore("bolt", tmpdir, nil)
	checkErr(err)

	// Open and use the database
	store, err := cayley.NewGraph("bolt", tmpdir, nil)
	checkErr(err)
	defer store.Close()
	qw := graph.NewWriter(store)

	/**start here created relationships**/

	// create accessType
	accessTypePrivate := AccessType{
		ID: quad.IRI("ex:Private"),
		//quad.IRI("ex:Private")
	}
	fmt.Printf("saving: %+v\n", accessTypePrivate)
	_, err = sch.WriteAsQuads(qw, accessTypePrivate)
	checkErr(err)
	err = qw.Close()
	checkErr(err)

	accessTypePublic := AccessType{
		ID: quad.IRI("ex:Public"),
	}
	fmt.Printf("saving: %+v\n", accessTypePublic)
	_, err = sch.WriteAsQuads(qw, accessTypePublic)
	checkErr(err)
	err = qw.Close()
	checkErr(err)

	userRole := Role{
		ID:        quad.IRI("ex:User"),
		HasAction: []quad.IRI{quad.IRI("ex:Read")},
	}
	fmt.Printf("saving: %+v\n", userRole)
	_, err = sch.WriteAsQuads(qw, userRole)
	checkErr(err)
	err = qw.Close()
	checkErr(err)

	adminRole := Role{
		ID:        quad.IRI("ex:Admin"),
		HasAction: []quad.IRI{quad.IRI("ex.Read"), quad.IRI("ex.Write")},
		//has action convert to array, atomatically do it twice : quad.Value ---interface to all value
	}
	fmt.Printf("saving: %+v\n", adminRole)
	_, err = sch.WriteAsQuads(qw, adminRole)
	checkErr(err)
	err = qw.Close()
	checkErr(err)

	//create action - read write
	readAction := Action{
		ID: quad.IRI("ex:Read"),
	}
	fmt.Printf("saving: %+v\n", readAction)
	_, err = sch.WriteAsQuads(qw, readAction)
	checkErr(err)
	err = qw.Close()
	checkErr(err)

	writeAction := Action{
		ID: quad.IRI("ex:Write"),
	}
	fmt.Printf("saving: %+v\n", writeAction)
	_, err = sch.WriteAsQuads(qw, writeAction)
	checkErr(err)
	err = qw.Close()
	checkErr(err)

	//create document
	pdf1 := Document{
		ID:                 quad.IRI("ex:Bijie.pdf"),
		HasAuthorizedAgent: quad.IRI("ex:Bijie"),
		Creator:            quad.IRI("ex:Bijie"),
		HasAccessType:      quad.IRI("ex:Public"),
	}
	fmt.Printf("saving: %+v\n", pdf1)
	_, err = sch.WriteAsQuads(qw, pdf1)
	checkErr(err)
	err = qw.Close()
	checkErr(err)

	pdf2 := Document{
		ID:                 quad.IRI("ex:Privatebook.pdf"),
		HasAuthorizedAgent: quad.IRI("ex:Bijie"),
		Creator:            quad.IRI("ex:Bijie"),
		HasAccessType:      quad.IRI("ex:Private"),
	}
	fmt.Printf("saving: %+v\n", pdf2)
	_, err = sch.WriteAsQuads(qw, pdf2)
	checkErr(err)
	err = qw.Close()
	checkErr(err)

	//create agent
	//agent bijie needs to both read and write

	bijie := Agent{
		ID:                            quad.IRI("ex:Bijie"),
		HasRole:                       quad.IRI("ex:Admin"),
		HasAuthorizedActionOnResource: []quad.IRI{quad.IRI("ex:Read"), quad.IRI("ex:Write")},
	}
	fmt.Printf("saving: %+v\n", bijie)
	_, err = sch.WriteAsQuads(qw, bijie)
	checkErr(err)
	err = qw.Close()
	checkErr(err)

	//create intermidiate object AuthorizedActionOnResource
	//this is the subclass of action

	authorizedAction := AuthorizedActionOnResource{
		ID:                  quad.IRI("ex:Read"), //subclass of Action
		HasResource:         quad.IRI("ex:Bijie.pdf"),
		HasActionOnResource: []quad.IRI{quad.IRI("ex:Read"), quad.IRI("ex:Write")}, //the value in Action
		SuperClasses:        []interface{}{quad.IRI("ex:Action")},
	} //subject: ID object:ex:Read predicate: ex:hasActionOnResource

	fmt.Printf("saving: %+v\n", authorizedAction)
	_, err = sch.WriteAsQuads(qw, authorizedAction)
	checkErr(err)
	err = qw.Close()
	checkErr(err)

	//print out all quads
	// Print quads
	fmt.Println("\n############################")

	fmt.Println("\nquads:")
	ctx := context.TODO()
	it := store.QuadsAllIterator().Iterate()
	defer it.Close()
	for it.Next(ctx) {
		q := store.Quad(it.Result())
		fmt.Println(q.NQuad())
	}

	fmt.Println("\n###########the path serach result###############")
	//define relationships -- search example
	p1 := cayley.StartPath(store, quad.IRI("ex:Bijie.pdf")).
		Out(quad.IRI("ex:creator")).Out(quad.IRI("ex:hasRole")).Out(quad.IRI("ex:hasAction"))

	err = p1.Iterate(nil).EachValue(nil, func(value quad.Value) {
		nativeValue := quad.NativeOf(value) // this converts RDF values to normal Go types
		fmt.Println(nativeValue)
		//fmt.Println(value)
	})

	if err != nil {
		log.Fatalln(err)
	}

	p2 := path.StartPath(store, quad.IRI("ex:Admin")).Follow(path.StartMorphism(quad.IRI("ex:creator"))).Out(quad.IRI("ex:hasRole"))
	err2 := p2.Iterate(nil).EachValue(nil, func(value quad.Value) {
		nativeValue := quad.NativeOf(value) // this converts RDF values to normal Go types
		fmt.Println(nativeValue)
		//fmt.Println(value)
	})

	if err2 != nil {
		log.Fatalln(err)
	}

}
