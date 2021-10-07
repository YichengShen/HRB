#!/usr/bin/python

from mininet.topo import Topo


class CustomTree(Topo):
    """
    Based on Mininet's TreeTopo.
    Allow an extra parameter, totalHosts, to put a limit on the total number of nodes.
    """

    def build(self, depth=1, fanout=2, totalHosts=2):
        # Numbering:  h1..N, s1..M
        self.hostNum = 1
        self.switchNum = 1
        # Build topology
        self.addTree(depth, fanout, totalHosts)

    def addTree(self, depth, fanout, totalHosts):
        """Add a subtree starting with node n.
           returns: last node added"""
        node = None
        isSwitch = depth > 0
        if isSwitch:
            node = self.addSwitch('s%s' % self.switchNum)
            self.switchNum += 1
            for _ in range(fanout):
                child = self.addTree(depth - 1, fanout, totalHosts)
                if child:
                    self.addLink(node, child)
        else:
            if (self.hostNum-1) < totalHosts:
                node = self.addHost('h%s' % self.hostNum)
                self.hostNum += 1
        return node


topos = { 'customtree' : ( lambda : CustomTree()) }