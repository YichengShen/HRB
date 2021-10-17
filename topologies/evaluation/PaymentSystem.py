from mininet.topo import Topo


class PaymentSystem(Topo):

    def build(self, e=10, bw=10000000, delay='50ms', jitter='25ms', loss=0):
        core = self.addSwitch('s1')

        for i in range(1, e+1):
            edge = self.addSwitch('s%s' % str(i+1))
            self.addLink(core, edge,
                         bw=bw,
                         delay=delay,
                         jitter=jitter,
                         loss=loss
                         )
            host =self.addHost('h%s' % i)
            self.addLink(edge, host,
                         bw=bw,
                         delay=delay,
                         jitter=jitter,
                         loss=loss
                         )