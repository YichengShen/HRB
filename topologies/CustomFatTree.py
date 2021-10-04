#!/usr/bin/python

from mininet.topo import Topo
from mininet.util import dumpNodeConnections


def indexGen(prefix):
    i = 1
    while True:
        yield '%s%s' % (prefix, i)
        i += 1


class CustomFatTree(Topo):

    def build(self, c=3, a=3, e=5, n=1):
        self.coreSwitchList = []
        self.aggSwitchList = []
        self.edgeSwitchList = []
        self.hostList = []
        self.sgen = indexGen("s")
        self.hgen = indexGen("h")

        for i in range(c):
            self.coreSwitchList.append(self.addSwitch(next(self.sgen)))

        for i in range(a):
            self.aggSwitchList.append(self.addSwitch(next(self.sgen)))

        for i in range(e):
            self.edgeSwitchList.append(self.addSwitch(next(self.sgen)))

        for core in self.coreSwitchList:
            for agg in self.aggSwitchList:
                self.addLink(core, agg)

        for agg in self.aggSwitchList:
            for edge in self.edgeSwitchList:
                self.addLink(agg, edge)

        for edge in self.edgeSwitchList:
            for i in range(n):
                host = self.addHost(next(self.hgen))
                self.hostList.append(host)
                self.addLink(edge, host)


topos = { 'customfattree' : ( lambda : CustomFatTree()) }