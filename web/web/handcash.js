// function handcashReceive(handle){
//   var req = new XMLHttpRequest();
//   req.onreadystatechange = function() {
//     console.log(req)
//       if (req.readyState === 4) {
//           var response = req.responseText;
//           var json = JSON.parse(response);
//           alert('Handcash API:\nTo send money to ' + handle + ' use the Cash address: ' + json['cashAddr']);
//       }
//   };
//
//   req.open('GET', 'http://api.handcash.io/api/receivingAddress/' + handle);
//   req.send(null);
//
// }


import {AuthorizationRequest,
    Cashport,
    GrantedAuthorization,
    PaymentRequestFactory,
    PersonalInfoPermission,
    SignTransactionRequestBuilder
} from 'cashport-sdk';

function cashport(){
  const appId = 'L77MZzEO72ZZSrRg58ysiGvveqFe51rK9lMDXKILD6YJf4lNibacSUx0xr979duv';
  console.log(appId)
  let cashport: Cashport = new Cashport();

  let cashport = new Cashport();
  let permissions = [PersonalInfoPermission.HANDLE, PersonalInfoPermission.FIRST_NAME, PersonalInfoPermission.LAST_NAME, PersonalInfoPermission.EMAIL];
  let authRequest = new AuthorizationRequest(permissions, appId);
  cashport.loadAuthorizationRequest('cashport-container', authRequest, {
      onDeny: () => {
          // Authorization not granted :(
          console.log("notOK")
      },
      onSuccess: (authorization) => {
          console.log("OK")
          let authToken = authorization.authToken;
          let expirationTimestamp = authorization.expirationTimestamp;
          let channelId = authorization.channelId;
          let spendLimitInSatoshis = authorization.spendLimitInSatoshis;

          let personalInfo = authorization.personalInfo;
          let handle = personalInfo.handle;
          let firstName = personalInfo.firstName;
          let lastName = personalInfo.lastName;
          let email = personalInfo.email;
      }
  });
}
//
// console.log("tst")
// cashport()
// handcashReceive("simeon")
