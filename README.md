# alfred-go-drive-search
An Alfred Google Drive workflow powered by Go

# Prerequisites 
In order to use this workflow you _**MUST**_ create your own Google OAuth2 app and generate a `client_id` & `client_secret`. Step-by-step instructions are available [here](https://learn2torials.com/a/google-oauth-setup). Once completed, copy and paste the `client_id` & `client_secret` into their respective workflow environment variables:

<img width="1013" alt="Screenshot 2022-11-03 at 11 05 12" src="https://user-images.githubusercontent.com/11773454/199705597-a8d58da2-4e01-48a9-9ce0-2b340bb21058.png">

# Download & installation
Grab the the latest version from the [releases page](https://github.com/coheff/alfred-go-drive-search/releases/tag/v0.1.0). Double click workflow file to import into Alfred.

If running on macOS Catalina or later, you _**MUST**_ add Alfred to the list of security exceptions for running unsigned software. Step-by-step instructions are available on the awgo wiki [here](https://github.com/deanishe/awgo/wiki/Catalina).

You can also start your application once ("with terminal") and except to open it. This way, Alfred will not be able to run just anything. To achieve it, you can install the workflow and then right click it to open the folder in Finder. Then right-click the 'alfred-go-drive-search' and open with Terminal. Then agree to open it. It will run and complain, not being started by Alfred. After that, the workflow will work (until an update of the executable).

_The first time you run the workflow you will be prompted to go through Googles's 3-legged OAuth2 flow. Once this has been completed a token will be stored in Keychain. The default expiry on tokens is 7 days (I believe); therefore you will be prompted once a week to refresh this token by re-enrolling via 3-legged OAuth2 flow:_

<img width="745" alt="Screenshot 2022-11-03 at 11 22 05" src="https://user-images.githubusercontent.com/11773454/199708495-a0d9d820-bd88-48c5-bf88-7b8583e0fbfb.png">

# Usage
- Trigger a name search using the keyword `gd` followed by a search query.
- Trigger a fulltext (document body) search by using the flag `-f=` followed by a search query e.g. `-f=burgers` will search for all documents contains the word "burgers". For longer fulltext searches you can wrap your query in quotes e.g. `-f="best burgers in town"`.
- You can combine name & fulltext searches e.g. `menus 2022 -f=burgers` will return a document whose name contains "menus 2022" and body contains "burgers".

# License
Distributed under the MIT License. See [LICENSE](https://github.com/coheff/alfred-hunt-tor/blob/main/LICENSE) for more information.
