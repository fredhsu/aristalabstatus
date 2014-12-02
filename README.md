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

Ideas:
For example, when you provision something, the very next thing you do is verify that the step was correct. You might ping an interface, or look at a counter, or view some protocol state. In these moments, the UX is less about the window that shows this to you and more about the workflow that captures validation as a natural next step after provisioning.

In another example, perhaps there is an alarm that indicates an interface is down or some loss threshold has been reached. The very next thing you will do is begin troubleshooting. Again, the wrapper is not so important as the steps it takes to actually capture meaningful information, correlate that information, and take appropriate action.

The point here is that the UX is determined more by the workflow than by the UI. In fact, the best UI is dependent more on a user’s context (where they are performing adjacent task
- See more at: http://www.plexxi.com/2014/11/networkings-ux-victims/#sthash.hBdNrZtA.dpuf


Add in Git version control of config
Add in Gerrit code review of config
Use config templates
