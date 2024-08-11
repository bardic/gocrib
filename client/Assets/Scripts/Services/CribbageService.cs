using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using RestSharp;
using RestSharp.Authenticators;
using System.Threading;


namespace com.openbracket.cribbage.services
{
    public class CribbageService
    {
        public async void Init(CancellationToken ct)
        {
            var options = new RestClientOptions("https://api.twitter.com/1.1")
            {
                Authenticator = new HttpBasicAuthenticator("username", "password")
            };
            var client = new RestClient(options);
            var request = new RestRequest("statuses/home_timeline.json");
            // The cancellation token comes from the caller. You can still make a call without it.
            var response = await client.GetAsync(request, ct);
        }
    }
}