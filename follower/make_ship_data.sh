#!/bin/bash

grep -f selected_ships.txt data/AIS_2020_06_01.csv > data_2020_06_01
grep -f selected_ships.txt data/AIS_2020_06_02.csv > data_2020_06_02
grep -f selected_ships.txt data/AIS_2020_06_03.csv > data_2020_06_03
grep -f selected_ships.txt data/AIS_2020_06_04.csv > data_2020_06_04
grep -f selected_ships.txt data/AIS_2020_06_05.csv > data_2020_06_05
grep -f selected_ships.txt data/AIS_2020_06_06.csv > data_2020_06_06
grep -f selected_ships.txt data/AIS_2020_06_07.csv > data_2020_06_07

cat data_2020_06_0* | sort -t',' -k 2 | awk 'NR % 40 == 0' > ship_data.csv
