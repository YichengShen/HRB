from mininet.topo import Topo
from mininet.link import TCLink
from mininet.util import irange


BW_LIMITED = 5

class SmartHome(Topo):
    """
    Based on Mininet's SingleSwitchTopo.
    Added limitations on bandwidth, delay, jitter, and loss.
    """

    def build(self, k=40, bw=5, delay='10ms', jitter='3ms', loss=2):
        "k: number of hosts"
        self.k = k
        switch = self.addSwitch('s1')
        for h in irange(1, k):
            host = self.addHost('h%s' % h)
            # Set 1/3 hosts to have highly limited bandwidth
            if h % 3 == 0:
                link_bw = BW_LIMITED
            else:
                link_bw = bw
            self.addLink(host, switch,
                         cls=TCLink,
                         bw=link_bw,
                         delay=delay,
                         jitter=jitter,
                         loss=loss)