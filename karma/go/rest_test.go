// Copyright 2015 Alex Browne.  All rights reserved.
// Use of this source code is governed by the MIT
// license, which can be found in the LICENSE file.

package main

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/go-humble/rest"
	"github.com/rusco/qunit"
)

type Todo struct {
	Id          int
	Title       string
	IsCompleted bool
}

func (t Todo) ModelId() string {
	return strconv.Itoa(t.Id)
}

func (t Todo) RootURL() string {
	return "http://localhost:3000/todos"
}

const (
	statusUnprocessableEntity = 422
)

func main() {
	// contentTypes is an array of all ContentTypes that we want to test for.
	// Note thate the test server must be capable of handling each type.
	contentTypes := []rest.ContentType{rest.ContentURLEncoded, rest.ContentJSON}
	// For each content type, we want to run all the tests and wait for the
	// tests to finish before continuing to the next type.
	for _, contentType := range contentTypes {
		client := rest.NewClient()
		client.ContentType = contentType
		wg := sync.WaitGroup{}
		// Need to update this if we add more tests.
		wg.Add(6)

		qunit.Test("ReadAll "+string(contentType), func(assert qunit.QUnitAssert) {
			qunit.Expect(2)
			done := assert.Async()
			go func() {
				expectedTodos := []*Todo{
					{
						Id:          0,
						Title:       "Todo 0",
						IsCompleted: false,
					},
					{
						Id:          1,
						Title:       "Todo 1",
						IsCompleted: false,
					},
					{
						Id:          2,
						Title:       "Todo 2",
						IsCompleted: true,
					},
				}
				gotTodos := []*Todo{}
				err := client.ReadAll(&gotTodos)
				assert.Ok(err == nil, fmt.Sprintf("client.ReadAll returned an error: %T %s", err, err))
				assert.Ok(reflect.DeepEqual(gotTodos, expectedTodos), fmt.Sprintf("Expected: %v, Got: %v", expectedTodos, gotTodos))
				done()
				wg.Done()
			}()
		})

		qunit.Test("Read "+string(contentType), func(assert qunit.QUnitAssert) {
			qunit.Expect(2)
			done := assert.Async()
			go func() {
				expectedTodo := &Todo{
					Id:          2,
					Title:       "Todo 2",
					IsCompleted: true,
				}
				gotTodo := &Todo{}
				err := client.Read("2", gotTodo)
				assert.Ok(err == nil, fmt.Sprintf("client.Read returned an error: %v", err))
				assert.Ok(reflect.DeepEqual(gotTodo, expectedTodo), fmt.Sprintf("Expected: %v, Got: %v", expectedTodo, gotTodo))
				done()
				wg.Done()
			}()
		})

		qunit.Test("Create "+string(contentType), func(assert qunit.QUnitAssert) {
			qunit.Expect(4)
			done := assert.Async()
			go func() {
				newTodo := &Todo{
					Title:       "Test",
					IsCompleted: true,
				}
				err := client.Create(newTodo)
				assert.Ok(err == nil, fmt.Sprintf("client.Create returned an error: %v", err))
				assert.Equal(newTodo.Id, 3, "newTodo.Id was not set correctly.")
				assert.Equal(newTodo.Title, "Test", "newTodo.Title was incorrect.")
				assert.Equal(newTodo.IsCompleted, true, "newTodo.IsCompleted was incorrect.")
				done()
				wg.Done()
			}()
		})

		qunit.Test("Update "+string(contentType), func(assert qunit.QUnitAssert) {
			qunit.Expect(4)
			done := assert.Async()
			go func() {
				updatedTodo := &Todo{
					Id:          1,
					Title:       "Updated Title",
					IsCompleted: true,
				}
				err := client.Update(updatedTodo)
				assert.Ok(err == nil, fmt.Sprintf("client.Update returned an error: %v", err))
				assert.Equal(updatedTodo.Id, 1, "updatedTodo.Id was incorrect.")
				assert.Equal(updatedTodo.Title, "Updated Title", "updatedTodo.Title was incorrect.")
				assert.Equal(updatedTodo.IsCompleted, true, "updatedTodo.IsCompleted was incorrect.")
				done()
				wg.Done()
			}()
		})

		qunit.Test("Delete "+string(contentType), func(assert qunit.QUnitAssert) {
			qunit.Expect(1)
			done := assert.Async()
			go func() {
				deletedTodo := &Todo{
					Id: 1,
				}
				err := client.Delete(deletedTodo)
				assert.Ok(err == nil, fmt.Sprintf("client.Update returned an error: %v", err))
				done()
				wg.Done()
			}()
		})

		qunit.Test("Non-2xx Response "+string(contentType), func(assert qunit.QUnitAssert) {
			qunit.Expect(5)
			done := assert.Async()
			go func() {
				err := client.Read("9999", &Todo{})
				assert.NotEqual(err, nil, "Expected an error from client.Read with invalid id, but got none.")
				httpErr, ok := err.(rest.HTTPError)
				assert.Equal(ok, true, fmt.Sprintf("Expected error of type rest.HTTPError but got %T", err))
				assert.Equal(httpErr.URL, Todo{}.RootURL()+"/9999", "httpErr.URL was incorrect")
				assert.Equal(httpErr.StatusCode, statusUnprocessableEntity, "httpErr.StatusCode was incorrect")
				assert.NotEqual(httpErr.Body, "", "httpErr.Body was empty")
				done()
				wg.Done()
			}()
		})

		// Wait for all the tests to finish before continuing to the next content type.
		wg.Wait()
	}
}
