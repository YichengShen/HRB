from mininet.net import Mininet
from mininet.util import dumpNodeConnections
from mininet.log import setLogLevel

from mininet.topo import LinearTopo
from mininet.topolib import TreeTopo
from topologies.CustomTree import CustomTree
from topologies.CustomFatTree import CustomFatTree
from topologies.SingleSwitchSlowSrc import SingleSwitchSlowSrc
from topologies.evaluation.SmartHome import SmartHome

from os import system
from time import sleep
import yaml


def test_network(net):
    print("Waiting switch connections")
    net.waitConnected()

    print("Testing network connectivity")
    percent_loss = net.pingAll()
    assert percent_loss==0., "Network is not fully connected"

def line_count(fname):
    """Returns the number of lines in a file."""
    count = 0
    with open(fname) as f:
        for _ in f:
            count += 1
    return count


if __name__ == '__main__':
    cfg_net = yaml.load(open('config_net.yaml', 'r'), Loader=yaml.FullLoader)
    topoName = cfg_net['topology']
    print("Mininet Topology:", topoName)

    # Clean
    system("sudo mn --clean")
    system("sudo rm watch.txt")
    system("sudo bash ./shell/clear_data.sh")

    # Set mininet log level
    setLogLevel('info')

    # Create network
    if topoName == "EvalSmartHome":
        topo = SmartHome(k=cfg_net['smart_home']['k'],
                         bw=cfg_net['smart_home']['bw'],
                         delay=cfg_net['smart_home']['delay'],
                         jitter=cfg_net['smart_home']['jitter'],
                         loss=cfg_net['smart_home']['loss'])
    elif topoName == "LinearTopo":
        topo = LinearTopo(k=cfg_net['linear_topo']['k'],
                          n=cfg_net['linear_topo']['n'])
    elif topoName == "CustomTree":
        topo = CustomTree(depth=cfg_net['custom_tree']['depth'],
                          fanout=cfg_net['custom_tree']['fanout'],
                          totalHosts=cfg_net['custom_tree']['total_hosts'])
    elif topoName == "CustomFatTree":
        topo = CustomFatTree(e=cfg_net['custom_fat_tree']['e'],
                             n=cfg_net['custom_fat_tree']['n'])
    elif topoName == "SingleSwitchSlowSrc":
        topo = SingleSwitchSlowSrc(k=cfg_net['single_switch_slow_src']['k'],
                                   bw_src=cfg_net['single_switch_slow_src']['bw_src'])
    elif topoName == "TreeTopo":
        topo = TreeTopo(depth=cfg_net['tree_topo']['depth'],
                        fanout=cfg_net['tree_topo']['fanout'])
    net = Mininet(topo)

    # Start network
    net.start()

    test_network(net)

    hosts = topo.hosts(sort=True)
    hosts = [net.get(host) for host in hosts]
    for i, h in enumerate(hosts):
        h.cmdPrint("export PATH=$PATH:/usr/local/go/bin")
        h.cmdPrint("cd ~/go/src/HRB")
        h.cmdPrint(f"go run main.go >> watch.txt &")

    print("Waiting for enough finishing messages")
    while True:
        num_lines = line_count("watch.txt")
        print(num_lines)
        if num_lines == len(hosts):
            break
        sleep(1)

    net.stop()

    # Run the analysis script to calculate mean throughput
    system("python3 analysis.py")



