from urllib3 import PoolManager
from urllib3.util import make_headers
from constants import NS_USERNAME, NS_PASSWORD
from xml.etree import ElementTree
from datetime import datetime

http = PoolManager()
ns = {'soap': 'http://www.w3.org/2003/05/soap-envelope'}


class NsApiError(Exception):
    pass


def ns_api(path, query):
    """https://www.ns.nl/reisinformatie/ns-api"""
    url = "http://webservices.ns.nl/" + path
    request = http.request(
        'GET',
        url,
        fields=query,
        headers=make_headers(basic_auth=NS_USERNAME + ':' + NS_PASSWORD)
    )
    if request.status != 200:
        faultstring = request.status
        try:
            tree = ElementTree.fromstring(request.data)
            faultstring = tree.find(
                './soap:Body/soap:Fault/faultstring', ns
            ).text
        except Exception:
            pass

        raise NsApiError(faultstring)

    tree = ElementTree.fromstring(request.data)
    if tree.tag == 'error':
        raise NsApiError(tree.find('message').text)

    return tree


class Departure():
    """Vertrektijd: item van de /ns-api-avt api (actuele vertrektijden)"""

    def __init__(self, node):
        self.id = int(node.find('RitNummer').text)
        self.track = int(node.find('VertrekSpoor').text)
        self.time = datetime.strptime(
            node.find('VertrekTijd').text,
            '%Y-%m-%dT%H:%M:%S%z'
        )
        if node.find('VertrekVertraging'):
            self.delay = node.find('VertrekVertragingTekst').text
        self.destination = node.find('EindBestemming').text
        if node.find('Opmerkingen'):
            self.remarks = []
            for remark in node.findall('./Opmerkingen/Opmering'):
                self.remarks.append(remark.text)

    def __str__(self):
        sb = []
        for key in self.__dict__:
            sb.append("{key}='{value}'".format(
                key=key, value=self.__dict__[key]))

        return '{' + ', '.join(sb) + '}'

    def __repr__(self):
        return self.__str__()


def departures():
    """Actuele vertrektijden vanuit Krommenie - Assendelft naar Amsterdam"""
    nodes = ns_api("ns-api-avt", {"station": "KMA"})
    items = []
    for node in nodes:
        item = Departure(node)
        if (item.track == 2):  # Richting Amsterdam
            items.append(item)

        if len(items) == 2:
            break

    return items
