# ggcr

simple cli application for reading messages in Goodgame chat

### installation (for linux)

	git clone https://github.com/Stasenko-Konstantin/ggcr 
	cd ggcr && ./build.sh
	
### usage

first of all you need get your channel id from GG. 
it can be done from share button in streamer dashboard.
there some html or smth else like: 

	<iframe frameborder="0" allowfullscreen width="800" height="450" src="https://goodgame.ru/player?214528"></iframe>
you need some number like this: 214528.
next add GGCR_ID env variable with this value, (re-)start terminal and simply run

	ggcr 
if you did all right then youll see smth like this:

	Goodgame: welcome!
and then it starts read all your chat history from start.
it load new messages in every minute.
