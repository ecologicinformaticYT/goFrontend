goFrontend is a fronted server for HTML websites (made with HTML/JavaScript/CSS & some images).
It doesn't support any other language.

To host a site on your PC with goFrontend, put it's files into goFrontend_<branch>/<your_OS>/site/user (the user folder is the root of the website),
then start the server using .exe/.bat (Windows) or .sh (Linux/Debian). You can delete the files that are already there, as they are just examples.

To see a page type <your-domain-name-or-localhost>/<path-to-page> (don't forget to precise the file extension) in your browser's address bar, then click "Enter".
To request the index.html page, just type <you-domain-name-or-localhost> and click "Enter".

Requirements : see in requirements.txt
License : MIT license (see it in LICENSE.txt)

Improvements (compared to V1) :
- use of conventional port (443) for HTTPS
- added docstrings
- added a development sandbox, so you can have 1 official version of your frontend code and 1 test version, running simultaneously on the same machine

How to use goFrontend :
- make sure you have golang >= 1.22.5 installed on your device
- put your website into the www folder (you can delete everything that's already there, except __devmode__)
- if you want to run goFrontend, without the sandbox, double-click "start_goFrontend.bat"
- if you want to run goFrontend with the sandbox, double-click "start_goFrontend_with_dev_sandbox.bat"
- if you want to run only the sandbox, double-click "start_dev_sandbox.bat"

The sandbox is located into : www/__devmode__/__devsandbox__  
The frontend code that the sandbox returns is located into www/__devmode__/__devsandbox__/www

Ports :
    hosted frontend code :
    - HTTP : 80
    - HTTPS : 443

    sandbox :
    - HTTP : 65080
    - HTTPS : 65443

WARNINGS :
1) This software is provided without any warranty
2) You are solely responsible for every damage caused by an improper use of that software 

©2024 Arthur CALON, goFrontend (V2), MIT License