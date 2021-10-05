#!/usr/bin/python

from mininet.topo import Topo
from mininet.util import dumpNodeConnections


def indexGen(prefix):
    i = 1
    while True:
        yield '%s%s' % (prefix, i)
        i += 1


class CustomFatTree(Topo):

    def build(self, e=5, n=1):
        self.aggSwitchList = []
        self.edgeSwitchList = []
        self.hostList = []
        self.sgen = indexGen("s")
        self.hgen = indexGen("h")

        core = self.addSwitch(next(self.sgen))

        for _ in range(e):
            agg = self.addSwitch(next(self.sgen))
            edge = self.addSwitch(next(self.sgen))
            self.aggSwitchList.append(agg)
            self.edgeSwitchList.append(edge)
            self.addLink(agg, edge)

        for agg in self.aggSwitchList:
            self.addLink(core, agg)

        for edge in self.edgeSwitchList:
            for _ in range(n):
                host = self.addHost(next(self.hgen))
                self.hostList.append(host)
                self.addLink(host, edge)


topos = { 'customfattree' : ( lambda : CustomFatTree()) }