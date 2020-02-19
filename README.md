#xpipe - extended pipe
####Overview
xpipe is a general networking hackers tool, primarilly to probe, test and diagnose network resources, but has a wide variety of uses.  
Inspired by an established tool, [netcat](https://en.wikipedia.org/wiki/Netcat), this does a very similar thing, but in a different way.  

Like netcat, it connects console streams (stdin, stdout) to network sockets (dial, listen) and vice versa.  It differs from netcat by having the concept of a pipeline.
Rather than acting on a single network location, it creates a pipeline of locations, with each steaming output into the input of the next.  This allows for
multiple, simultaneous connections on a single pipeline, including multiple inbound connections in a single pipeline.  

Using this principle the pipeline can open out and listen in on multiple ports to any other host simultaneously, feeding streams from multiple locations as a single pipe.  
Connections are not limited to other hosts.  The pipeline can even listen to itself!!  Establishing connections to itself, the pipe can 'back feed' or loop the stream
feeding the output of the pipe, back into itself, to be passed back as the next response to another inbound request.  
In addition to looping a single pipe, multiple pipeline can communicate with one another such as in an OS type command line pipe:  
`xpipe @5555 @6666 | grep "name\: " | lookupname.sh | xpipe :5555`  

Here a regular cmdline "... | grep | lookupname.sh" is enclosed in two pipelines.  The first accepts two connections. "clients" connect on the last inbound socket (6666)
and their request stream is piped into the cmdline, grep.  The output of the shell script is then piped into another pipeline, which sends the result back to the first pipeline
over socket 5555.  On arrival, the first pipeline pipes the response into the original, open connection.

This allows the pipeline to react to inbound connections, processing or selecting the response and feeding that response back into the pipe, 
and so, back to the requester/client, making it possible to create a workable web service from a single command line!  

xpipe is designed to be run as a command line tool or within shell scripts.
It operates at a low level in the network stack, where few of the well known application layer protocols exist or fancy things like HTTP :-).  
Known to the OSI bods as layer 4, Down here we have raw bytes moving about between numbered points.
  
\disclaimer\  
xpipe is designed to play with.  Much of what is does can be achieved with existing tools and a bit of scripting, often more easily.  
`netcat` and `curl` alone can do most of it and `nmap` is most certainly a whole lot more powerful at scanning and probing networks than this tool will ever be.    
What it is designed for, is make some of these tasks easily repeatable and scriptable.  To present the moving of data in a logical manner that doesn't take 3 years of 
university to understand.  
  
  
Hope you can have fun with it.  

####Use cases
xpipe is designed as a testing and diagnostic tool.  It can scan for open ports, ip addresses.  probe open sockets and test responses etc.  
These kind of things most 'ops' and development teams do all the time when maintaining infrastructure or developing network based services.  
xpipe was originally written to serve unit testing for development of web based services.  
Unit testing often requires 'mock' services to provide realistic responses to code being tested, which can prove complex to set up and maintain.  
Using xpipe as a mock server, pre-defined responses can be placed into text files and served based on predictable requests.  
Inversely, xpipe is useful for testing development of server side code.  Again using pre-formed requests in text files, these can be piped to the test 
server, repeatedly and reliably, to test for expected results.  All these tasks can usually be accomplised with simple xpipe commands and some predefined text files.  

Although designed for development work, it might prove useful to anybody that works on network resources.  Having the ability to merge scripts and network request/responces
allows for some powerful and funky tools to be crafted with basic shell scriping.

####Usage
`xpipe [OPTIONS]... SOURCE [TARGET]...`

SOURCE and any optional TARGETS are chained together, with the stream from SOURCE being passed to the first TARGET,  
and the stream from that TARGET, passed to the next TARGET and so on.  The final TARGETs output is written to stdout.  

SOURCE  
Represents where the initial stream of data will come from.  
It the only required argument.  


TARGET    
The location the stream from SOURCE or the previous TARGET is sent.  
Any number of TARGETs can be stated.  Each target is fed the output of the TARGET before it to form the pipeline.  


Both SOURCE and TARGET are stated as one of the following:  
* `-`   A single dash represents the console, stdin and stdout.  
* `<netaddr>`  create an outbound (dial) connection and stream the pipe through it.  
* `@<netaddr>` create an inbound (listening) port and wait for inbound connections.  
* `+`   A plus sign will echo the pipe to stdout, without breaking the pipe.  

`-`  console pipe accepts upstream data and writes it to the stdout.  
Any input from the stdin is passed downstream, into the pipe.  
  
'+' Echo is simply a window into the pipe.  upstream data is passed directly to the downstream pipe.  
However it is also 'echoed' or copied to the stdout.  

`netaddr` Network address is a `<hostname:port>` combination.  
Hostname can be a host, an ip or ommited for localhost.  
A port is always required.  (No such thing as a default port down here in the basement, level 4 ;-) )  
  
  
**examples**
  + `remotehost:8007`
  + `localhost:7788`
  + `192.168.10.12:8080`
  + `:8090`   
  
`netaddr`  For outbound connections, the upstream pipeline is blocked until the outbound connection is established.
Once established, the pipeline is passed through it, upstream pipe is pushed to the remote and remote 'response' is fed to the downstream pipe.  
The connection remains open until remote closes or returns an EOF or local times out or closes pipe.  

`@netaddr`  For inbound connections, upstream data is blocked until an inbound connection is made.  
On connection, the pipeline passes "through" the connection, upstream pipeline data is pushed as 'response' data, to the remote  
and the remotes 'request' data is pushed downstream, to the rest of the pipe.  
  
  
####Simple Examples
`xpipe - remotehost:1234`  *one way chat*  
Source is stdin '-' being streamed to remotehost on port 1234.
This will push what is being typed on the console, to the remote host.  Any response from the remote host will be piped back to the local console out (stdout)  
  
`cat myfile.txt | xpipe remotehost:1234`  *copy file to remote*  
This copies the file 'myfile.txt' to the remotehost, listening on port 1234.  

`cat myfile.txt | xpipe - @:5555`  *blind request / response*  
Waits for inbound connections on port 5555.  Source is stdin '-', which is the myfile.txt stream, blocked until a connection arrives.  
The rest of the pipeline can still continue, provided another source of input is available, a new connection being one.    
A new, inbound connection, connects the pipeline so the blocked stream is piped back to that connection and any stream read from the remote is piped downstream, into the pipe.  
  
  
`xpipe @:5555 + remotehost:666`  *pipe window / dumb proxy*
This accepts connections on localhost:5555 and streams the inbound connection over to remotehost on port 1234.
The `+` in the middle copies the stream to the localhost stdout.  The response from remotehost is also streamed back to the stdout.  
Useful for local development sometimes, to re/mis direct packets to another service or just see what is being passed through the pipe.  
  

####More Examples
  
two way chat
Server side
`xpipe @5555 @6666 - :5555`  
Two inbound port listen, 5555 is a private port, 6666 the external one.   6666 feed into the 'console'.  
The console feeds into an outbound connection, looping back to the first inbound.

client
`xpipe - localhost:6666`  
Client connects their console to the 'external' port 6666 and they have two way chat.  
  
  
  
*One line web server*  
`xpipe @:1234 @:5555 | INPUT=`\`cat -\` `;  if [[ $INPUT == *"testpage"* ]]; then cat testdata/testpage.html; else bin/dothis.sh $INPUT; fi; | xpipe :1234`    
  
Demonstrates a looping pipe where two seperate pipelines enclose some standard bash scripting.  
  
The regular bashscript:  
`INPUT=`\`cat -\` `;  if [[ $INPUT == *"testpage"* ]]; then cat testdata/testpage.html; else bin/dothis.sh $INPUT; fi`
(With weird cat - because I can't escape back quotes in markdown :-/)  
cat - pipes the stdin, into the INPUT variable, checks if that variable contains our 'testpage' in the request and returns that page if it does.  
otherwise it executes a bash script `bin/dothis.sh`, passing the INPUT request as the command line.  
The output of either the testpage or script is the result of the script.  

The pipeline:  
The pipeline is formed around the bash script, connecting an inbound network socket as the script input and using a loopback pipe to capture the result of the script
and feed it back to the same network connection.   
    
Having this 'open' scripting process in the center of a pipeline is what gives this tool its extrodiary flexability.  
Inbound data can be processed, using simple grep, awk, sed like tools to generate live/active responses, which are feed back to the client.  
  
  
####Tips
When creating mock server responses or requests, its helpful to have a template to work from, rather that type the whole thing yourself.  
To grab examples, use xpipe to listen to real requests and copy them to text files.  
`xpipe @5555` | myrequest.txt  
Then using a browser, adjust your test url so host is localhost and port is 5555 `localhost:5555` and fire the request.  
Obviously, there will be no response, but you now have a text file 'myrequest.txt' containing what the browser just sent.  
Edit that file and re-adjust the urls and any query parameters / headers you want to change and it can be used to fire a 'real' looking request.  
`cat myeditedrequest.txt | xpipe realhost:80 | myresponse.txt`  
Now you have a text file 'myresponse.txt' of, hopefully, a realistic response.  You usually have to tinker with the request a bit.  
The result are files you can use to test / probe a remote service with scriptable requests:  
as with the example above, or mock a server responding when testing clients.  
`cat myresponse.txt | xpipe @:8080`  
