arista-labstatus
================

Web application to monitor the status of Arista EOS switches

# getSwitchConfigs
- Grabs switch configs and writes them to files
## TODO:
* Use nested channels to do the write -- pipelining
    * One channel to initiate
    * another channel to do the file writing
* Allow CLI input of what directory to write to and what file to read in
* Seperate to a different directory/util

