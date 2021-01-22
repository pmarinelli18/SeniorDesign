CREATE DATABASE BattleShip;
USE BattleShip;

Create Table Accounts(
	userName VARCHAR(50), 
    password VARCHAR(50),
    PRIMARY KEY(userName)
);

INSERT INTO Accounts(userName, password) 
	Values
	('userName 1', 'password'),
    ('userName 2', 'password2');