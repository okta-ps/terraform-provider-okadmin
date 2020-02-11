
import requests
import sys, getopt
import re
from bs4 import BeautifulSoup
import json

def main(argv):
    baseUrl = ''            #This is Okta org tenant for example "https://dev-497881.oktapreview.com"
    OktaDomain = ''
    Org2OrgAppId = ''   #This is Org2Org application ID.
    APIToken = ''           #This is the spoke tenant API token
    username = ''          #This is the Username of your OKTA tenant
    password = ''
    answer = ''
    CreateUsers=False
    UpdateUserAttributes=False
    DeactivateUsers=False
    SyncPassword=False
    
    opts, args = getopt.getopt(argv,"b:d:o:t:u:p:a:c:n:e:s:",["baseUrl=","OktaDomain=","Org2OrgAppId=","APIToken=","username=","password=","answer=","CreateUsers=","UpdateUserAttributes=","DeactivateUsers=","SyncPassword="])
   
    for opt, arg in opts:
       if opt in ('-b', '--baseurl'):
           baseUrl = arg
       elif opt in ('-d', '--OktaDomain'):
           OktaDomain = arg
       elif opt in ('-o', '--Org2OrgAppId'):
           Org2OrgAppId = arg
       elif opt in ('-t', '--APIToken'):
           APIToken = arg
       elif opt in ('-u', '--username'):
           username = arg
       elif opt in ('-p', '--password'):
           password = arg
       elif opt in ('-a', '--answer'):
           answer = arg
       elif opt in ('-c', '--CreateUsers'):
           CreateUsers = arg
       elif opt in ('-n', '--UpdateUserAttributes'):
           UpdateUserAttributes = arg
       elif opt in ('-e', '--DeactivateUsers'):
           DeactivateUsers = arg
       elif opt in ('-s', '--SyncPassword'):
           SyncPassword = arg
    #Create the session
    S=requests.Session()
    
    baseUrl1= 'https://'+baseUrl+'.'+OktaDomain
    adminBaseUrl= 'https://'+baseUrl+'-admin.'+OktaDomain
   
    #Grab the adminXsrfToken
    adminXsrfToken = xsrf(S,baseUrl1,adminBaseUrl,username,password,answer)
    targetId = trgId(S,adminBaseUrl,Org2OrgAppId)
    sourceId = srcId(S,adminBaseUrl,targetId)
    
    #This is the method to Enable Integration
    print(Enable_Integration(S,adminBaseUrl,adminXsrfToken,Org2OrgAppId,APIToken))

    print(Update_OKTAOrg2Org_ToApp(S,adminBaseUrl,adminXsrfToken,Org2OrgAppId,CreateUsers,UpdateUserAttributes,DeactivateUsers,SyncPassword,sourceId,targetId))
    
    
#Get XSRF token
def xsrf(S,baseUrl,adminBaseUrl,username,password,answer):
    #Call the authn api
    body = {'username': username, 'password':password}
    auth = S.post(baseUrl+'/api/v1/authn', json=body)
    json_response=auth.json()
        
    status = json_response["status"]
    
    if(status=="MFA_REQUIRED"):
        stateToken = json_response["stateToken"]
        resp_dict = json.loads(auth.text)
        data=(resp_dict['_embedded']['factors'][0]) 
        factorId = (data['id'])
        header = {"accept": "application/json","Content-Type" : "application/json","user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36"}
        body={"answer":answer,"stateToken":stateToken}    
        response=S.post(baseUrl+'/api/v1/authn/factors/'+factorId+'/verify?rememberDevice=false',json=body,headers=header)   
        json_response=response.json()
        
    sessionToken = json_response['sessionToken']

    #Call the api with sessionToken
    response=S.post(baseUrl+'/login/sessionCookieRedirect?token='+sessionToken+'&redirectUrl=/')

    #Get admin login token
    response=S.get(baseUrl+'/home/admin-entry')
    Token=(re.compile('"token":\["(.*)"\]').search(response.text)).group(1)

    #Use token to sign in to Okta Admin app, and get admin xsrfToken
    body = {'token': Token}
    header = {"Content-Type" : "application/x-www-form-urlencoded"}
    response=S.post(adminBaseUrl+'/admin/sso/request',data=body,headers=header)
    soup = BeautifulSoup(response.text,'html.parser')
    adminXsrfToken=soup.find(id='_xsrfToken').text
    return(adminXsrfToken)

#Get OrgId
def orgId(S,baseUrl):
    response=S.get(baseUrl+'/.well-known/okta-organization')
    orgId = response.json()
    return orgId['id']

#Get targetId
def trgId(S,adminBaseUrl,Org2OrgAppId):
    response=S.get(adminBaseUrl+'/api/v1/apps/'+Org2OrgAppId+'/user/types/effective?expand=schema%2CappLogo%2Capp')
    tId = response.json()
    targetId=tId['id']
    return targetId

#Get sourceId    
def srcId(S,adminBaseUrl,targetId):    
    response=S.get(adminBaseUrl+'/api/internal/v1/mappings?target='+targetId)
    sId=response.json()
    srcsId=sId[0]
    sourceId=srcsId['sourceId']
    return sourceId


def Enable_Integration(S,adminBaseUrl,adminXsrfToken,Org2OrgAppId,APIToken):                       # Method to Enable Integration
    body={'_xsrfToken':adminXsrfToken,
      'enabled':True,
      '_enabled':'on',
      'token': APIToken,
      '_preferUsernameOverEmail':'on',
      'importGroups':True,
      '_importGroups':'on',
      '_profileMaster':'on',
      '_pushNewAccount':'on',
      '_pushProfile':'on',
      '_pushDeactivation':'on',
      '_pushPassword':'on',
      'syncUniquePassword':True,
      '_cycleSyncedPassword':'on'}
    response=S.post(adminBaseUrl+"/admin/app/okta_org2org/instance/"+Org2OrgAppId+"/settings/user-mgmt", data=body)
    return(response.status_code)
    
def Update_OKTAOrg2Org_ToApp(S,adminBaseUrl,adminXsrfToken,Org2OrgAppId,CreateUsers,UpdateUserAttributes,DeactivateUsers,SyncPassword,sourceId,targetId):                 # Method to Updte OktaOrg2OrgToApp
    body={'_xsrfToken':adminXsrfToken,
      'enabled':True,
      '_enabled':'on',
      '_preferUsernameOverEmail':'on',
      '_profileMaster':'on',
      'pushNewAccount':CreateUsers,
      '_pushNewAccount':'on',
      'pushProfile':UpdateUserAttributes,
      '_pushProfile':'on',
      'pushDeactivation':DeactivateUsers,
      '_pushDeactivation':'on',
      'pushPassword':SyncPassword,
      '_pushPassword':'on',
      'syncUniquePassword':True,
      '_cycleSyncedPassword':'on'}
    response=S.post(adminBaseUrl+"/admin/app/okta_org2org/instance/"+Org2OrgAppId+"/settings/user-mgmt", data=body)
   
    #Set OKTA Org2Org Initial State Attribute Mapping
    headers={"Accept" : "application/json", "Content-Type" : "application/json","User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.66 Safari/537.36","X-Okta-XsrfToken" : adminXsrfToken}
    request_payload={"sourceId":sourceId,
                     "targetId":targetId,
                     "propertyMappings":[{"targetField":"email","sourceExpression":"user.email","pushStatus":"DONT_PUSH"},
                                         {"targetField":"firstName","sourceExpression":"user.firstName","pushStatus":"DONT_PUSH"},
                                         {"targetField":"lastName","sourceExpression":"user.lastName","pushStatus":"DONT_PUSH"},
                                         {"targetField":"middleName","sourceExpression":"user.middleName","pushStatus":"DONT_PUSH"},
                                         {"targetField":"honorificPrefix","sourceExpression":"user.honorificPrefix","pushStatus":"DONT_PUSH"},
                                         {"targetField":"honorificSuffix","sourceExpression":"user.honorificSuffix","pushStatus":"DONT_PUSH"},
                                         {"targetField":"title","sourceExpression":"user.title","pushStatus":"DONT_PUSH"},
                                         {"targetField":"displayName","sourceExpression":"user.displayName","pushStatus":"DONT_PUSH"},
                                         {"targetField":"nickName","sourceExpression":"user.nickName","pushStatus":"DONT_PUSH"},
                                         {"targetField":"profileUrl","sourceExpression":"user.profileUrl","pushStatus":"DONT_PUSH"},
                                         {"targetField":"primaryPhone","sourceExpression":"user.primaryPhone","pushStatus":"DONT_PUSH"},
                                         {"targetField":"streetAddress","sourceExpression":"user.streetAddress","pushStatus":"DONT_PUSH"},
                                         {"targetField":"city","sourceExpression":"user.city","pushStatus":"DONT_PUSH"},
                                         {"targetField":"mobilePhone","sourceExpression":"user.mobilePhone","pushStatus":"DONT_PUSH"},
                                         {"targetField":"secondEmail","sourceExpression":"user.secondEmail","pushStatus":"DONT_PUSH"},
                                         {"targetField":"state","sourceExpression":"user.state","pushStatus":"DONT_PUSH"},
                                         {"targetField":"zipCode","sourceExpression":"user.zipCode","pushStatus":"DONT_PUSH"},
                                         {"targetField":"countryCode","sourceExpression":"user.countryCode","pushStatus":"DONT_PUSH"},
                                         {"targetField":"postalAddress","sourceExpression":"user.postalAddress","pushStatus":"DONT_PUSH"},
                                         {"targetField":"preferredLanguage","sourceExpression":"user.preferredLanguage","pushStatus":"DONT_PUSH"},
                                         {"targetField":"locale","sourceExpression":"user.locale","pushStatus":"DONT_PUSH"},
                                         {"targetField":"timezone","sourceExpression":"user.timezone","pushStatus":"DONT_PUSH"},
                                         {"targetField":"userType","sourceExpression":"user.userType","pushStatus":"DONT_PUSH"},
                                         {"targetField":"employeeNumber","sourceExpression":"user.employeeNumber","pushStatus":"DONT_PUSH"},
                                         {"targetField":"costCenter","sourceExpression":"user.costCenter","pushStatus":"DONT_PUSH"},
                                         {"targetField":"organization","sourceExpression":"user.organization","pushStatus":"DONT_PUSH"},
                                         {"targetField":"division","sourceExpression":"user.division","pushStatus":"DONT_PUSH"},
                                         {"targetField":"department","sourceExpression":"user.department","pushStatus":"DONT_PUSH"},
                                         {"targetField":"managerId","sourceExpression":"user.managerId","pushStatus":"DONT_PUSH"},
                                         {"targetField":"manager","sourceExpression":"user.manager","pushStatus":"DONT_PUSH"},
                                         {"targetField":"initialStatus","sourceExpression":"\"active_with_pass\"","pushStatus":"DONT_PUSH"}]}
    response=S.put(adminBaseUrl+"/api/internal/v1/mappings", json=request_payload,headers=headers)
    
    return(response.status_code)       

if __name__== "__main__":
        main(sys.argv[1:])