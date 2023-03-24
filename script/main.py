import xmltodict
import sys

class bcolors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKCYAN = '\033[96m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'

def readFile(filename) -> str:
  with open(filename) as f:
    return f.read()
  
def parseXML(xml) -> dict:
  return xmltodict.parse(xml)

def findTagXPath(xml: dict, tag: str, currentPath: str):
  list = []
  for k, v in xml.items():
    #print(k)
    if(k == tag):
      list.append(tuple([v,currentPath + "/" + k]))
    elif(type(v) == dict):
      list += findTagXPath(v, tag, currentPath + "/" + k)
    elif(type(v) == list):
      for i in v:
        list += findTagXPath(i, tag, currentPath + "/" + k)
  return list

def printOutput(list):
  for elem in list:
    print(bcolors.HEADER,"XPATH is :", elem[1],bcolors.ENDC)
    print(elem[0])

def retrieveConf(path, tag):
  xml = readFile(path)
  config = parseXML(xml)
  result = findTagXPath(config, tag, '')
  print(printOutput(result))

if __name__ == '__main__':
  # get args
  if(len(sys.argv) < 3):
    print("Usage: python3 main.py <path to xml file>")
    exit(1)
  else:
    path = sys.argv[1]
    tag = sys.argv[2]
    retrieveConf(path,tag)