CREATE DATABASE BattleShip;
USE BattleShip;

Create Table Accounts(
    userName VARCHAR(50), 
    password VARCHAR(50),
    PRIMARY KEY(userName)
);

INSERT INTO Accounts(userName, password) 
    Values
    ('userName1', 'password'),
    ('userName2', 'password2'),
    ('u', 'p');
    
Create Table BoatState(
    IpAddress VARCHAR(50),
    UserName VARCHAR(50),
    navigationPosition INTEGER DEFAULT 0, 
    RadarState VARCHAR(15) DEFAULT 'Enabled', 
    ShipHealth INTEGER DEFAULT 100, 
    NumberOfCannons INTEGER DEFAULT 3, 
    TorpedoState VARCHAR(15) DEFAULT 'Standby',
    TorpedoDamage INTEGER DEFAULT 0,
    FinishedMiniGame BOOL DEFAULT false,
    GameActive BOOL DEFAULT true,
    PRIMARY KEY(ipAddress)
);

SET SQL_SAFE_UPDATES = 0;