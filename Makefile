run:
	go run sprite/main.go -o /tmp/core.css -p ~/work/rmn/www/gui/sass -d ~/work/rmn/www/gui/im/sass ~/work/rmn/www/gui/sass/rmn/_core.scss
home:
	go run sprite/main.go -gen ~/work/rmn/www/gui/build/im -o ~/work/rmn/www/gui/build/css/tests/RMN-15500/home.css -p ~/work/rmn/www/gui/sass -d ~/work/rmn/www/gui/im/sass ~/work/rmn/www/gui/sass/tests/RMN-15500/home.scss
profile:
	go run sprite/main.go --cpuprofile=sprite.prof -o /tmp/home.css -p ~/work/rmn/www/gui/sass -d ~/work/rmn/www/gui/im/sass ~/work/rmn/www/gui/sass/tests/RMN-15500/home.scss
