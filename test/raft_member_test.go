// Copyright 2018 The tiglabs raft Authors.
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

package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/tiglabs/raft/proto"
)

func TestMemberWithNoLease(t *testing.T) {
	f, w := getLogFile("changemember_nolease.log")
	defer func() {
		w.Flush()
		f.Close()
	}()

	servers := initTestServer(peers, false, true)
	fmt.Println("waiting electing leader....")
	leadServer := waitElect(servers, w)
	printStatus(servers, w)

	// test add node
	w.WriteString(fmt.Sprintf("[%s] Add new node \r\n", time.Now().Format(format_time)))
	leader, term := leadServer.raft.LeaderTerm(1)
	newServer := createRaftServer(4, leader, term, peers, false, true)
	// add node
	resolver.addNode(4, 0)
	fmt.Println("starting add node")
	leadServer.sm.AddNode(proto.Peer{ID: 4})
	fmt.Println("added node")
	time.Sleep(time.Second)
	servers = append(servers, newServer)
	printStatus(servers, w)

	fmt.Println("starting put data")
	if err := leadServer.sm.Put("test2", "test2_val"); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	if vget, err := newServer.sm.Get("test2"); err != nil || vget != "test2_val" {
		t.Fatal("new add node not get the data")
	}
	fmt.Println("success put data")

	// test remove node
	w.WriteString(fmt.Sprintf("[%s] Remove node \r\n", time.Now().Format(format_time)))
	fmt.Println("starting remove node")
	leadServer.sm.RemoveNode(proto.Peer{ID: 4})
	fmt.Println("removed node")
	fmt.Println("starting put data")
	if err := leadServer.sm.Put("test3", "test3_val"); err != nil {
		t.Fatal(err)
	}
	fmt.Println("success put data")
	newServers := make([]*testServer, 0)
	for _, s := range servers {
		if s.nodeID == newServer.nodeID {
			s.raft.Stop()
		} else {
			newServers = append(newServers, s)
		}
	}
	servers = newServers
	time.Sleep(100 * time.Millisecond)
	newServer = createRaftServer(4, 0, 10, append(peers, proto.Peer{ID: 4}), false, false)
	servers = append(servers, newServer)
	time.Sleep(10 * time.Second)
	printStatus(servers, w)
	resolver.delNode(4)

	for _, s := range servers {
		s.raft.Stop()
	}
	time.Sleep(100 * time.Millisecond)
}

func TestMemberWithLease(t *testing.T) {
	f, w := getLogFile("changemember_lease.log")
	defer func() {
		w.Flush()
		f.Close()
	}()

	servers := initTestServer(peers, true, true)
	fmt.Println("waiting electing leader....")
	leadServer := waitElect(servers, w)
	printStatus(servers, w)

	// test add node
	w.WriteString(fmt.Sprintf("[%s] Add new node \r\n", time.Now().Format(format_time)))
	leader, term := leadServer.raft.LeaderTerm(1)
	newServer := createRaftServer(4, leader, term, peers, true, true)
	// add node
	resolver.addNode(4, 0)
	fmt.Println("starting add node")
	leadServer.sm.AddNode(proto.Peer{ID: 4})
	fmt.Println("added node")
	time.Sleep(time.Second)
	servers = append(servers, newServer)
	printStatus(servers, w)

	fmt.Println("starting put data")
	if err := leadServer.sm.Put("test2", "test2_val"); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	if vget, err := newServer.sm.Get("test2"); err != nil || vget != "test2_val" {
		t.Fatal("new add node not get the data")
	}
	fmt.Println("success put data")

	// test remove node
	w.WriteString(fmt.Sprintf("[%s] Remove node \r\n", time.Now().Format(format_time)))
	fmt.Println("starting remove node")
	leadServer.sm.RemoveNode(proto.Peer{ID: 4})
	fmt.Println("removed node")
	fmt.Println("starting put data")
	if err := leadServer.sm.Put("test3", "test3_val"); err != nil {
		t.Fatal(err)
	}
	fmt.Println("success put data")
	newServers := make([]*testServer, 0)
	for _, s := range servers {
		if s.nodeID == newServer.nodeID {
			s.raft.Stop()
		} else {
			newServers = append(newServers, s)
		}
	}
	servers = newServers
	time.Sleep(100 * time.Millisecond)
	newServer = createRaftServer(4, 0, 10, append(peers, proto.Peer{ID: 4}), false, false)
	servers = append(servers, newServer)
	time.Sleep(10 * time.Second)
	printStatus(servers, w)
	resolver.delNode(4)

	for _, s := range servers {
		s.raft.Stop()
	}

	time.Sleep(100 * time.Millisecond)
}

func TestMemberWithLeaseRemoveBeforeAdd(t *testing.T) {
	f, w := getLogFile("changemember_lease.log")
	defer func() {
		w.Flush()
		f.Close()
	}()

	servers := initTestServer(peers, true, true)
	fmt.Println("waiting electing leader....")
	leadServer := waitElect(servers, w)
	printStatus(servers, w)

	// test add node
	w.WriteString(fmt.Sprintf("[%s] Add new node \r\n", time.Now().Format(format_time)))
	leader, term := leadServer.raft.LeaderTerm(1)
	newServer := createRaftServer(4, leader, term, peers, true, true)
	// add node
	resolver.addNode(4, 0)
	fmt.Println("starting add node")
	leadServer.sm.AddNode(proto.Peer{ID: 4})
	fmt.Println("added node")
	time.Sleep(time.Second)
	servers = append(servers, newServer)
	printStatus(servers, w)

	fmt.Println("starting put data")
	if err := leadServer.sm.Put("test2", "test2_val"); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	if vget, err := newServer.sm.Get("test2"); err != nil || vget != "test2_val" {
		t.Fatal("new add node not get the data")
	}
	fmt.Println("success put data")

	// test remove node
	w.WriteString(fmt.Sprintf("[%s] Remove node \r\n", time.Now().Format(format_time)))
	fmt.Println("starting remove node")
	leadServer.sm.RemoveNode(proto.Peer{ID: 4})
	fmt.Println("removed node")
	fmt.Println("starting put data")
	if err := leadServer.sm.Put("test3", "test3_val"); err != nil {
		t.Fatal(err)
	}
	fmt.Println("success put data")
	newServers := make([]*testServer, 0)
	for _, s := range servers {
		if s.nodeID == newServer.nodeID {
			s.raft.Stop()
		} else {
			newServers = append(newServers, s)
		}
	}
	servers = newServers
	time.Sleep(100 * time.Millisecond)
	newServer = createRaftServer(4, 0, 10, append(peers, proto.Peer{ID: 4}), false, false)
	servers = append(servers, newServer)
	time.Sleep(10 * time.Second)
	printStatus(servers, w)
	resolver.delNode(4)

	for _, s := range servers {
		s.raft.Stop()
	}

	time.Sleep(100 * time.Millisecond)
}

func TestResetPeersWithLease(t *testing.T) {
	f, w := getLogFile("reset_peers_lease.log")
	defer func() {
		w.Flush()
		f.Close()
	}()
	peers = []proto.Peer{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7}}
	servers := initTestServer(peers, true, true)
	fmt.Println("waiting electing leader....")
	leadServer := waitElect(servers, w)
	printStatus(servers, w)

	var newServers []*testServer
	for index, s := range servers {
		if index > 3 {
			newServers = append(newServers, s)
			continue
		}
		s.raft.Stop()
	}

	newPeers := []proto.Peer{{ID: 5}, {ID: 6}, {ID: 7}}
	if len(newServers) > 0 {
		for _, ns := range newServers {
			go func(ns *testServer) {
				if err := ns.sm.ResetPeers(newPeers); err != nil {
					fmt.Printf("reset peer test error, err[%v]", err)
				}
			}(ns)
		}

	} else {
		fmt.Printf("no new peers")
		return
	}

	leadServer = waitElect(newServers, w)

	// test add node
	if _, err := w.WriteString(fmt.Sprintf("[%s] Add new node after reset peers \r\n", time.Now().Format(format_time))); err != nil {
		fmt.Printf("write string error, err[%v]", err)
	}
	leader, term := leadServer.raft.LeaderTerm(1)
	newServer := createRaftServer(4, leader, term, newPeers, true, true)
	// add node
	resolver.addNode(4, 0)
	fmt.Println("starting add node")
	leadServer.sm.AddNode(proto.Peer{ID: 4})
	fmt.Println("added node")
	time.Sleep(time.Second)
	newServers = append(newServers, newServer)
	printStatus(newServers, w)

	for _, s := range newServers {
		s.raft.Stop()
	}
	time.Sleep(100 * time.Millisecond)
}
