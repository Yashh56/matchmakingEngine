@echo off
call go run ./cmd -player_id=01 -mmr=1920 -region=asia -ping=30  
call go run ./cmd -player_id=02 -mmr=2020 -region=asia -ping=20  
call go run ./cmd -player_id=03 -mmr=1950 -region=asia -ping=30  
call go run ./cmd -player_id=04 -mmr=2050 -region=asia -ping=50  
