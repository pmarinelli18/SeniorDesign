CREATE DATABASE BattleShip;
USE BattleShip;

Create Table Accounts(
    userName VARCHAR(50), 
    password VARCHAR(50),
    PRIMARY KEY(userName)
);

Create Table BoatState(
    ipAddress VARCHAR(50),
    cannonActive BOOL DEFAULT false, 
    mountedMGActive BOOL DEFAULT false, 
    torpedoActive BOOL DEFAULT false, 
    radarActive BOOL DEFAULT false, 
    shipHealthStatus DOUBLE DEFAULT 1.0,
    PRIMARY KEY(ipAddress)
);

INSERT INTO Accounts(userName, password) 
    Values
    ('userName 1', 'password'),
    ('userName 2', 'password2');
    
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
