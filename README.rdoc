= Usage

 tobi@tobbook8 ~  $ tinny --help
 Usage of tinny:
   -compiler="make": compiler to execute. make by default.
   -executable="server": executable to launch. Must support -pXXXX for port assignment
   -port=3333: port

= Purpose

Tinny is a tiny web server thats used during development. When tinny recieves a request it will launch a compiler task ( such as make ), start an executable in the current directory and then forward the HTTP request to that server.

= Requirements 

You must be working on a http sever. The http server needs to understand the -pXXXX command line switch that specifies the server it's running on. 

= Example:

 tobi@tobbook8 ~ $ cd ~/Code/go/imagery/src
 tobi@tobbook8 ~/Code/go/imagery/src $ tinny -compiler=gb
 Listening on port 3333
   -> forwarding to executable server -p4357


