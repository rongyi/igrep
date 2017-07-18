package igrep

import (
	"strings"
)

var (
	DropStats []string = []string{"n_bytes", "cookie", "idle_age", "hard_age"}
	ARPResponder string = "NXM_OF_ARP_OP[],"
)

// cookie=0xb87550d1bb38f685, duration=451145.578s, table=20, n_packets=0, n_bytes=0, idle_age=65534, hard_age=65534, priority=2,dl_vlan=20,dl_dst=fa:16:3e:6e:b2:b1 actions=strip_vlan,load:0x62 ->NXM_NX_TUN_ID[],output:9
func FlowAction(flow string) string {
	index := strings.Index(flow, "actions")
	if index == -1 {
		return flow
	}
	return flow[index:]
}

func FlowBeforeAction(flow string) string {
	index := strings.Index(flow, " actions")
	if index == -1 {
		return flow
	}
	return flow[:index]
}

func FlowPropDict(flow string) map[string]string {
	ret := make(map[string]string)
	kvs := strings.Split(flow, ",")
	for _, kv := range kvs {
		stats := strings.Split(strings.TrimSpace(kv), "=")
		if len(stats) > 1 {
			ret[stats[0]] = stats[1]
		} else {
			ret[stats[0]] = ""
		}
	}

	return ret
}

func FlowPropFilter(dict map[string]string) {
	for _, f := range DropStats {
		delete(dict, f)
	}
}

func FlowDropStats(flow string) string {
	ret := flow
	for _, f := range DropStats {
		fIndex := strings.Index(ret, f)
		if fIndex == -1 {
			continue
		}
		before := ret[:fIndex]
		after := ret[fIndex:]
		afterCommaIndex := strings.Index(after, ",")
		if afterCommaIndex != -1 {
			ret = before + ret[afterCommaIndex + 2 + len(before):]
		} else {
			ret = before
		}
	}
	return ret
}

func ActionARPClean(action string) string {
	if strings.Index(action, "NXM_OF_ARP_OP") == -1 {
		return action
	}
	return "actions=proxyarp..." + action[strings.Index(action, "mod_dl_src"):]
}

func ActionLearnClean(action string) string {
	if strings.Index(action, "learn") == -1 {
		return action
	}
	return "actions=(learn and put to table 20)," + action[strings.Index(action, "),") + len("),"):]
}
