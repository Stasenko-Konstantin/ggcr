# ggcr

a simple cli application for reading messages in the Goodgame chat

### installation (for linux)

    git clone https://github.com/Stasenko-Konstantin/ggcr 
    cd ggcr && ./build.sh

### usage

first of all, you need to get your channel id from GG.
it can be done by clicking the share button on the streamer dashboard.
there is some HTML code, smth like:

    <iframe frameborder="0" allowfullscreen width="800" height="450" src="https://goodgame.ru/player?214528"></iframe>

you need to extract the number from it, for example: `214528`.
next, add the `GGCR_ID` environment variable with this value, (re)run your terminal, and simply run

    ggcr 

if you set everything up correctly, you'll see something like this:

    Goodgame: welcome!

and then it starts reading your chat history from the beginning.
it loads new messages every minute.
