/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package zookeeper

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestShowChildren(t *testing.T) {
	methodsMap, err := NewZookeeperMetadataReport(
		"dubbo-cli",
		[]string{"127.0.0.1:2181"},
	).ShowChildren()
	assert.NoError(t, err)
	for k, v := range methodsMap {
		t.Logf("interface: %s\n", k)
		t.Logf("methods: %v\n", v)
	}
}
