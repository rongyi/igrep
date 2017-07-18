package igrep

import (
	"testing"
)

func TestFlowAction(t *testing.T) {
	input := `cookie=0xb87550d1bb38f685, duration=451145.578s, table=20, n_packets=0, n_bytes=0, idle_age=65534, hard_age=65534, priority=2,dl_vlan=20,dl_dst=fa:16:3e:6e:b2:b1 actions=strip_vlan,load:0x62->NXM_NX_TUN_ID[],output:9`
	expect := `actions=strip_vlan,load:0x62->NXM_NX_TUN_ID[],output:9`
	action := FlowAction(input)
	if action != expect {
		t.Fatal("FlowAction fail")
	}
}

func TestFlowBefore(t *testing.T) {
	input := `cookie=0xb87550d1bb38f685, duration=451145.578s, table=20, n_packets=0, n_bytes=0, idle_age=65534, hard_age=65534, priority=2,dl_vlan=20,dl_dst=fa:16:3e:6e:b2:b1 actions=strip_vlan,load:0x62->NXM_NX_TUN_ID[],output:9`
	expect := `cookie=0xb87550d1bb38f685, duration=451145.578s, table=20, n_packets=0, n_bytes=0, idle_age=65534, hard_age=65534, priority=2,dl_vlan=20,dl_dst=fa:16:3e:6e:b2:b1`
	before := FlowBeforeAction(input)
	if before != expect {
		t.Fatal("FlowAction fail")
	}
}

func TestDropFlow(t *testing.T) {
	input := `cookie=0xb87550d1bb38f685, duration=451145.578s, table=20, n_packets=0, n_bytes=0, idle_age=65534, hard_age=65534, priority=2,dl_vlan=20,dl_dst=fa:16:3e:6e:b2:b1`
	expect := `duration=451145.578s, table=20, n_packets=0, priority=2,dl_vlan=20,dl_dst=fa:16:3e:6e:b2:b1`
	before := FlowDropStats(input)
	if before != expect {
		t.Fatal("FlowAction fail")
	}
}
