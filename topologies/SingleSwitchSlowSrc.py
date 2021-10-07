#!/usr/bin/python

from mininet.topo import Topo
from mininet.util import irange


class SingleSwitchSlowSrc(Topo):
    """
    Based on Mininet's SingleSwitchTopo.
    Added option to limit the source node's bandwidth.
    """

    def build(self, k=2, bw_src=0.05):
        "k: number of hosts"
        self.k = k
        switch = self.addSwitch('s1')
        for h in irange(1, k):
            host = self.addHost('h%s' % h)
            if h == 1:
                self.addLink(host, switch, bw=bw_src)
            self.addLink(host, switch)