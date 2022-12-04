## This file is a part of Bashly standard library
##
## Usage:
## Use any of the functions below to color or format a portion of a string.
##
##   echo "before $(red this is red) after"
##   echo "before $(green_bold this is green_bold) after"
##
## Color output will be disabled if `NO_COLOR` environment variable is set
## in compliance with https://no-color.org/
##
print_in_color() {
  local color="$1"
  shift
  if [[ -z ${NO_COLOR+x} ]]; then
    printf "$color%b\e[0m" "$*";
  else
    printf "%b" "$*";
  fi
}

red() { print_in_color "\e[31m" "$*"; }
green() { print_in_color "\e[32m" "$*"; }
yellow() { print_in_color "\e[33m" "$*"; }
blue() { print_in_color "\e[34m" "$*"; }
magenta() { print_in_color "\e[35m" "$*"; }
cyan() { print_in_color "\e[36m" "$*"; }
white() { print_in_color "\e[37m" "$*"; }
gray() { print_in_color "\e[0;37m" "$*"; }
bold() { print_in_color "\e[1m" "$*"; }
underlined() { print_in_color "\e[4m" "$*"; }
red_bold() { print_in_color "\e[1;31m" "$*"; }
green_bold() { print_in_color "\e[1;32m" "$*"; }
yellow_bold() { print_in_color "\e[1;33m" "$*"; }
blue_bold() { print_in_color "\e[1;34m" "$*"; }
magenta_bold() { print_in_color "\e[1;35m" "$*"; }
cyan_bold() { print_in_color "\e[1;36m" "$*"; }
white_bold() { print_in_color "\e[1;37m" "$*"; }
red_underlined() { print_in_color "\e[4;31m" "$*"; }
green_underlined() { print_in_color "\e[4;32m" "$*"; }
yellow_underlined() { print_in_color "\e[4;33m" "$*"; }
blue_underlined() { print_in_color "\e[4;34m" "$*"; }
magenta_underlined() { print_in_color "\e[4;35m" "$*"; }
cyan_underlined() { print_in_color "\e[4;36m" "$*"; }




println_in_color() {
  local color="$1"
  shift
  if [[ -z ${NO_COLOR+x} ]]; then
    printf "$color%b\e[0m\n" "$*";
  else
    printf "%b\n" "$*";
  fi
}

red_ln() { println_in_color "\e[31m" "$*"; }
green_ln() { println_in_color "\e[32m" "$*"; }
yellow_ln() { println_in_color "\e[33m" "$*"; }
blue_ln() { println_in_color "\e[34m" "$*"; }
magenta_ln() { println_in_color "\e[35m" "$*"; }
cyan_ln() { println_in_color "\e[36m" "$*"; }
white_ln() { println_in_color "\e[37m" "$*"; }
gray_ln() { println_in_color "\e[0;37m" "$*"; }
bold_ln() { println_in_color "\e[1m" "$*"; }
underlined_ln() { println_in_color "\e[4m" "$*"; }
red_bold_ln() { println_in_color "\e[1;31m" "$*"; }
green_bold_ln() { println_in_color "\e[1;32m" "$*"; }
yellow_bold_ln() { println_in_color "\e[1;33m" "$*"; }
blue_bold_ln() { println_in_color "\e[1;34m" "$*"; }
magenta_bold_ln() { println_in_color "\e[1;35m" "$*"; }
cyan_bold_ln() { println_in_color "\e[1;36m" "$*"; }
white_bold_ln() { println_in_color "\e[1;37m" "$*"; }
red_underlined_ln() { println_in_color "\e[4;31m" "$*"; }
green_underlined_ln() { println_in_color "\e[4;32m" "$*"; }
yellow_underlined_ln() { println_in_color "\e[4;33m" "$*"; }
blue_underlined_ln() { println_in_color "\e[4;34m" "$*"; }
magenta_underlined_ln() { println_in_color "\e[4;35m" "$*"; }
cyan_underlined_ln() { println_in_color "\e[4;36m" "$*"; }
