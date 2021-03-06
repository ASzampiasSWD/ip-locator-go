# ip-locator-go
This go code generates a Google Earth file from a text file of IP Addresses. The code also generates a geo-stats.csv file with information about the IPs. 
IP information includes latitude, longitude, city, region, country, ISP, and if the IP belongs to a VPN provider. 

## Steps
Install the Go Library by twpayne. Navigate to the ip-locator-go folder and issue this command:
<pre>
go get github.com/twpayne/go-kml
</pre>

Run the script

<pre>
go run main.go
- OR - 
go build
./ip-locator-go
</pre>

## Optional VPN and ISP Detection

Steps:
<pre>
1. Register for an account on the website <a href="https://iphub.info/">https://iphub.info/</a>.
2. Generate a <b>free</b> API key. This will let you make 1000 requests a day.
3. Make a file called iphub-info.txt within the ip-locator-go directory.
4. Paste the API key in the file.
</pre>

## How-to Video
[![How To Video](https://img.youtube.com/vi/mDdzq0UE_d8/0.jpg)](https://www.youtube.com/watch?v=mDdzq0UE_d8)
