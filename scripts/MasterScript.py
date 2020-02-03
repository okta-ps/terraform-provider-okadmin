# Script for OKTAOrg2Org 
import requests
import re

baseUrl ='https://dev-497881.oktapreview.com'
adminBaseUrl = "https://dev-497881-admin.oktapreview.com"
#Variables that needs to be set and updated for OKTAOrg2Org 
CreateUsers = 'True'
UpdateUserAttributes = 'False'
DeactivateUsers = 'False'
SyncPassword = 'False'
SyncUniquePassword = 'True'

# For batch mode, set username and password
username='alka_maurya@tecnics.com'
password='Zaq12wsx'

#Get XSRF token

#Sign in to Okta. Keep $session (cookies), and get user xsrfToken
#Create the session
S=requests.Session()

#Call the authn api
body = {'username': username, 'password':password}
auth = S.post(baseUrl+'/api/v1/authn', json=body)
json_response=auth.json()
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
adminXsrfToken=(re.compile('<span.* id="_xsrfToken">(.*)</span>').search(response.text)).group(1)

#Perform the OKTAOrg2Org 
body={'_xsrfToken':adminXsrfToken,'enabled':True,'_enabled':'on','_preferUsernameOverEmail':'on','_profileMaster':'on','pushNewAccount':CreateUsers,'_pushNewAccount':'on','pushProfile':UpdateUserAttributes,'_pushProfile':'on','pushDeactivation':DeactivateUsers,'_pushDeactivation':'on','pushPassword':SyncPassword,'_pushPassword':'on','syncUniquePassword':SyncUniquePassword,'_cycleSyncedPassword':'on'}
response=S.post(adminBaseUrl+"/admin/app/okta_org2org/instance/0oap5lq94h5xh4UHO0h7/settings/user-mgmt", data=body)
print(response.status_code)
print(response.text)
    
