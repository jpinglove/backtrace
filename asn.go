package main

import (
	"fmt"
	"github.com/fatih/color"
	"net"
	"strings"
)

type Result struct {
	i int
	s string
}

var (
	ips = []string{"123.185.103.7", "182.204.47.29", "202.96.209.133", "210.22.97.1","211.136.112.200"}
	names = []string{"大连电信", "大连手机电信", "上海电信", "上海联通", "上海移动"}
	m = map[string]string{"AS4134": "电信163  [普通线路]", "AS4809": "电信CN2  [优质线路]", "AS4837": "联通4837 [普通线路]",
		"AS9929": "联通9929 [优质线路]", "AS58807": "移动CMIN2[优质线路]", "AS9808": "移动CMI  [普通线路]", "AS58453": "移动CMI  [普通线路]"}
)

func trace(ch chan Result, i int) {

	hops, err := Trace(net.ParseIP(ips[i]))
	if err != nil {
		s := fmt.Sprintf("%v %-15s %v", names[i], ips[i], err)
		ch <- Result{i, s}
		return
	}

	for _, h := range hops {
		for _, n := range h.Nodes {
			asn := ipAsn(n.IP.String())
			as := m[asn]
			var c func(a ...interface{}) string
			switch asn {
			case "":
				continue
			case "AS9929":
				c = color.New(color.FgHiYellow).Add(color.Bold).SprintFunc()
			case "AS4809":
				c = color.New(color.FgHiMagenta).Add(color.Bold).SprintFunc()
			case "AS58807":
				c = color.New(color.FgHiBlue).Add(color.Bold).SprintFunc()
			default:
				c = color.New(color.FgWhite).Add(color.Bold).SprintFunc()
			}

			s := fmt.Sprintf("%v %-15s %-23s", names[i], ips[i], c(as))
			ch <- Result{i, s}
			return
		}
	}
	c := color.New(color.FgRed).Add(color.Bold).SprintFunc()
	s := fmt.Sprintf("%v %-15s %v", names[i], ips[i], c("测试超时"))
	ch <- Result{i, s}
}

func ipAsn(ip string) string {

	switch {
	case strings.HasPrefix(ip, "59.43"):
		return "AS4809"
	case strings.HasPrefix(ip, "202.97"):
		return "AS4134"
	case strings.HasPrefix(ip, "218.105") || strings.HasPrefix(ip, "210.51"):
		return "AS9929"
	case strings.HasPrefix(ip, "219.158"):
		return "AS4837"
	case strings.HasPrefix(ip, "223.120.19") || strings.HasPrefix(ip, "223.120.17") || strings.HasPrefix(ip, "223.120.16"):
		return "AS58807"
	case strings.HasPrefix(ip, "223.118") || strings.HasPrefix(ip, "223.119") || strings.HasPrefix(ip, "223.120") || strings.HasPrefix(ip, "223.121"):
		return "AS58453"
	default:
		return ""
	}
}
