# deployment
Scripts, utilities, and examples for deploying CoreDNS.


## MacOS
The default settings will proxy all requests to hostnames not found in your host file to Google's DNS-over-HTTPS.

To install:
  - Run `brew tap "coredns/deployment" "https://github.com/coredns/deployment"`
  - Run `brew install coredns`
  - Run `sudo brew services start coredns`
  - test with `dig google.com @127.0.0.1` and you should see  `SERVER: 127.0.0.1#53(127.0.0.1)`

Using CoreDNS as your default resolver:
 - Open Network Preferences
 - Select your interface i.e Wi-Fi
 - Click `Advanced`
 - Select the `DNS` tab
 - Click the `+` below the `DNS Servers` list box
 - Type `127.0.0.1` and hit enter
 - Click `OK`
 - Click `Apply`
