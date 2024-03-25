// Copyright (c) 2023 Palantir Technologies. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cycles

import (
	"testing"
)

// testGraph is the following graph:
// 1<--3<--6<->7
// |  ^^   ^   ^
// | / |   |   | /|
// V/  |   |   |/ |
// 2<--4<->5<--8<-|
var testGraph = newGraph[int](8).
	addNode(1).
	addNode(2).
	addNode(3).
	addNode(4).
	addNode(5).
	addNode(6).
	addNode(7).
	addNode(8).
	addEdgesByID(1, 2).
	addEdgesByID(2, 3).
	addEdgesByID(3, 1).
	addEdgesByID(4, 2, 3, 5).
	addEdgesByID(5, 4, 6).
	addEdgesByID(6, 3, 7).
	addEdgesByID(7, 6).
	addEdgesByID(8, 5, 7, 8)

func TestReverseGraph(t *testing.T) {
	expected := newGraph[int](8).
		addNode(1).
		addNode(2).
		addNode(3).
		addNode(4).
		addNode(5).
		addNode(6).
		addNode(7).
		addNode(8).
		addEdgesByID(1, 3).
		addEdgesByID(2, 1, 4).
		addEdgesByID(3, 2, 4, 6).
		addEdgesByID(4, 5).
		addEdgesByID(5, 4, 8).
		addEdgesByID(6, 5, 7).
		addEdgesByID(7, 6, 8).
		addEdgesByID(8, 8)
	actual := reverseGraph(testGraph)
	assertGraphsAreEqual(t, expected, actual)
}
