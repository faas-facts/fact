from time import sleep

from factclient.fact import Fact
from factclient.io.TCPSender import TCPSender

def handler(event, context):
    io = TCPSender()
    Fact.boot({"inlcudeEnvironment": False, "io": io, "send_on_update": True})
    Fact.start(context, event)
    sleep(1)
    Fact.update(context, "test updating", {1: "experiment_python"})
    sleep(1)
    Fact.done(context, "test done", ["no more args"])
    return None
