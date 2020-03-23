package mc_test

import (
	"context"
	"fmt"
	"testing"

	"git.in.zhihu.com/knowledge-market/education/pkg/util/tools/mc"
)

func init() {
	mc.Setup(mc.LocalStore)
}

func getName(id int) (string, error) {
	return fmt.Sprintf("%d", id), nil
}

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

const (
	firstName  = "first name"
	lastName   = "last name"
	firstName2 = "first name2"
	lastName2  = "last name2"
)

func getPerson() (Person, error) {
	return Person{
		FirstName: firstName,
		LastName:  lastName,
		Age:       30,
	}, fmt.Errorf("defined error")
}

func TestCache1(t *testing.T) {
	var name string
	ret, _ := mc.Cache("", "v1", 600, getName, 1)
	mc.Unmarshal(ret[0], &name)
	if name != "1" {
		t.Fatalf("name expected: %s, actual: %s", "1", name)
	}
}

func TestCache2(t *testing.T) {
	var person Person
	ret1, err := mc.Cache("", "v1", 500, getPerson)

	mc.Unmarshal(ret1[0], &person)

	if person.FirstName != firstName {
		t.Fatalf("first name expected: %s, actual: %s", firstName, person.FirstName)
	}

	if err == nil || err.Error() != "defined error" {
		t.Fatalf("error expected: %s, actual: %s", "defined error", err.Error())
	}
}

func TestCache3(t *testing.T) {
	var getPerson = func(context.Context) (Person, error) {
		var p Person
		ret, err := mc.Cache("", "v1", 0,
			func(_ context.Context) (Person, error) {
				return Person{
					FirstName: firstName2,
					LastName:  lastName2,
					Age:       30,
				}, nil
			},
			context.TODO())
		mc.Unmarshal(ret[0], &p)
		return p, err
	}

	person, err := getPerson(context.TODO())

	if person.LastName != lastName2 || err != nil {
		t.Fatalf("last name expected: %s, actual: %s", lastName2, person.LastName)
	}

	// 测试缓存取数据
	person2, err := getPerson(context.TODO())

	if person2.LastName != lastName2 || err != nil {
		t.Fatalf("last name expected: %s, actual: %s", lastName2, person.LastName)
	}
}

type Group struct {
	Num     int
	Persons []*Person
}

func GetGroup(ctx context.Context, name string) Group {
	var g Group
	ret, _ := mc.Cache("", "v1", 0,
		func(_ context.Context, _ string) Group {
			persons := []*Person{
				&Person{
					Age: 10,
				},
				&Person{
					Age: 20,
				},
			}
			return Group{
				Num:     2,
				Persons: persons,
			}
		}, ctx, name)

	mc.Unmarshal(ret[0], &g)
	return g
}

func TestCache4(t *testing.T) {
	ctx := context.TODO()
	group := GetGroup(ctx, "vip")

	if group.Num != 2 {
		t.Fatalf("group num expected: %d, actual: %d", 2, group.Num)
	}

	if group.Persons[0].Age != 10 {
		t.Fatalf("person 1 age expected: %d, actual: %d", 10, group.Persons[0].Age)
	}

	if group.Persons[1].Age != 20 {
		t.Fatalf("person 2 age expected: %d, actual: %d", 20, group.Persons[1].Age)
	}

	group = GetGroup(ctx, "vip")

	if group.Num != 2 {
		t.Fatalf("group num expected: %d, actual: %d", 2, group.Num)
	}

	if group.Persons[0].Age != 10 {
		t.Fatalf("person 1 age expected: %d, actual: %d", 10, group.Persons[0].Age)
	}

	if group.Persons[1].Age != 20 {
		t.Fatalf("person 2 age expected: %d, actual: %d", 20, group.Persons[1].Age)
	}
}
