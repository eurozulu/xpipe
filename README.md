#xpipe - extended pipe
####Overview

---
>.
Connects console streams (stdin, stdout) to network sockets (dial, listen) and vice versa.  
primarily used for diagnostics and testing but it can be used for a wide variety of functions based around networking.  


> **concept**  
> Creates a pipeline of byte streams, passing the output of one stream, into the input of the next,
> until usually arriving at the stdout.  
>
> 
---

Three pipe types:  
* std console
* netAddr (outbound)
* Cross netAddr (inbound)  


Start with an input stream (stdIn)
End with an output stream (stdOut)

-> X X... =>

remote request reposnse to stdout
xpipe remotehost:8090

stdin to remote socket
xpipe - remotehost:8090

remote socket inbound to stdout
xpipe ->:8090 -

proxy in :5555 to remotehost:8080
xpipe ->:5555 remotehost:8080

response server listens on 8080, posts connections top 8090 to process,  
result from 8090 posted to tmp port 1234, which pssed back result.  

xpipe ->:1234 ->:8080 :8090 :1234
xpipe ->:1234 :8080 | workit.sh |xpipe - :1234
wait on

xpipe connects 'pipes' of data together much like the regular operating system pipes, which connect `stdin` and `stdout`.  
In addition to the regular stdin and stdout streams, xpipe supports more complex streams such as network sockets and command line processes.    
Using a single command line, chaining the arguments together as 'named streams', complex piped processes can be created.  
Insprired by an old favorite tool, "netcat", `nc`,  

####Concept
`xpipe` accepts one or more arguments, each argument representing a 'named stream'  
Each additional argument is treated as another named stream, connecting its input to the previous arguments output.  
The final argument/named stream has its ouput piped into the stdout.  
Each named stream accepts an input stream and provides an output.  


Stream types are:  

std console  
Provides data from stdin and accepts data into stdout.  

network address  
Provides data from a network dial/request and accepts data a

Stream types are identified by the format of each argument.
* `ssss:nnnn`  outbound / remote network address
* `:nnnn`      inbound / listening port
* `-`          standard input/output streams (stdin  stdout)
* `[<cmd line> ... ]  Command line instruction to execute `


accepts one or more arguments, each treated as a 'named stream', to form a 'pipe' stream.   
Reading the data from the first stream it 'pipes' the resulting bytes onto the next, named stream in the pipe.  
This next 'stream' receives the resulting bytes and outputs a stream of its own.  
This continues along the pipe until the final named stream is encountered, which has its output piped to the standard console out.  

Named streams are items from which a stream of bytes can either be read or written to or both.  
Such items are:  



 raw network sockets to OS streams and vice versa.  

####What it does

At the most basic level, xpipe performs simple one way streaming by connecting the standard console streams (stdin, stdout)
to a network socket, allowing reading and writing of the respective console streams to the socket.  

-> write stream  
=> read stream  

stdin=> ->netaddr=> ->stdout    Reads stdin and writes it to a network address.  The response from the network address is written to stdout.  

stdin=> ->netaddr=> ->stdout | stdin=> ->netaddr=> ->stdout   As above, showing how streams are 'piped' from one place to another, using the standard console streams.  

netaddr=> ->stdout   Reads a network address and dumps the result to stdout.  

xpipe supports two way streaming functionality by listening on a local port for new connections, piping the inbound stream to
'stdout'.  
stdin=> ->localnetaddr=> ->stdout
As inbound connections are made on the local port, both stdin and stdout are conntected to the network socket.  
The network connection is read and piped into stdout whilst stdin is read and piped back into the open connection, as the 'response'.  

they are treated, initially as read stream,
having their content streamed to stdout.  Simultaneously, std

xpipe can also perform two way streaming (request/response) by listen on a local port for inbound traffic and streaming
the stdout back to the
####Usage
`xpipe  [FLAGS]... STREAM [STREAM]...`  

`STREAM` represents a stream.  Valid streams are:  
* Network addresses e.g. `hostname:port`  
* Executable commands  e.g. `${grep mytext}`
* Standard console streams e.g. `-`  

The first STREAM argument is mandatory, from which the initial source of data is read.
Any following STREAMs receive the resulting data from the previous stream and create a resulting stream of data themselves.  
This continues until the final STREAM's output is written to the standrd console.  

`FLAGS`  represent the various options to control how the pipe functions.  


Stream types  
Network addresses are in the form of a hostname:port  

####Basic Examples  

**Inbound streams**  
To listen on local port `8080` and dump the inbound streams to the console:  
`xpipe `  

To listen on local port `5555` and dump the inbound stream to a file:  
`xpipe :5555 >> myfile.txt`  


**Outbound streams**  
To stream the standard input console to a remote network socket `otherhost:8090`:    
`xpipe - otherhost:8090`  

To stream a file to a remote network socket `otherhost:8090`:    
`cat myfile.txt | xpipe - otherhost:8090`  

To listen on local port `5555` and 'proxy' the stream to a remote socket `otherhost:8090`:    
`xpipe :5555 otherhost:8090`  


####More Examples  
To test a REST endpoint `http://www.spoofer.org/v1/endpointtest?id=1234` from a pre-defined request file, 
first generate a HTTP request stream in the form of a simple text file, using `curl`.
Start `xpipe` to listen and dump the result to a file:  
`xpipe --timeout 10s :5555 >> myrequest.txt`  
Now generate the request using 'curl', replacing the desired host (and socket) with our local listening point:  
`curl http://localhost:5555/v1/endpointtest?id=1234`  

When the connection closes, the 'myrequest.txt' will contain the http request curl generated and posted to the listening socket.
Open the file and replace the hostname (and socket) to the desired url, `www.spoofer.org`  

Now we can use this file to test the end point by streaming the contents to our desired location with:  
`cat myrequest.txt | xpipe - http://www.spoofer.org/v1/testendpoint`  
The result will be streamed to the standard out.  


To create a simple mock server service with a pre-defined response, first generate a response file from an existing service, in the form of a simple text file, using curl.  
`curl -i --raw http://www.spoofer.org/v1/endpointtest?id=1234 >> myResponse.txt`  

`cat myResonse | xpipe :5555 | `
# xpipe
