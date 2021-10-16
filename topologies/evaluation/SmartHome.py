from mininet.topo import Topo
from mininet.util import irange


class SmartHome(Topo):
    """
    Based on Mininet's SingleSwitchTopo.
    Added limitations on bandwidth, delay, jitter, and loss.
    """

    def build(self, k=40, bw=50000000, delay='10ms', jitter='3ms', loss=2):
        "k: number of hosts"
        self.k = k
        switch = self.addSwitch('s1')
        for h in irange(1, k):
            host = self.addHost('h%s' % h)
            self.addLink(host, switch,
                         bw=bw,
                         delay=delay,
                         jitter=jitter,
                         loss=loss)