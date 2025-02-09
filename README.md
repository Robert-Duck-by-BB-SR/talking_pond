# What is talking pond?
It's a simple messenger. That's it. You can send messages and receive them. 
With vim motions. In terminal.

# How do I use it?
First you need a server [TPS](https://github.com/Robert-Duck-by-BB-SR/tps).  
To access that server you need to get an IP of the server and a key that is generated on the server for your user.
Ask your server administrator to generate one for you by providing your desired username.
Next run the client
```bash
./talking_pond
```
In the login screen first insert ip of your server and key in the next field.

# Movements

## Normal Mode
`ctrl-w` Enter Window Mode  
`:` Enter Command Mode  
`h/j/k/l` move to the previous/next item in selected window  
`Enter` Open chat/press a button  
`I` (shift i) jump to Insert Mode  
`i` (in the inputtable item) enter Insert Mode  
`ctrl-u/ctrl-d` if content of the component is vertically scrollable  
`q` to exit modal window  
## Command Mode  
`:q` quit application (shocker)  
`:new` open a modal window to create a new conversation  
`ESC/ctrl-c` enter normal mode  
## Window Mode  
`h/j/k/l` move to the previous/next window  
`ESC/ctrl-c/ENTER` enter normal mode  
## Insert Mode  
just type the text wtf do you expect to see here?  
`ESC/ctrl-c` enter normal mode  


# License
Any company or individual shall be publicly shamed for breaking current license.
