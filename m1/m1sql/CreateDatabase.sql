/* --------------------------------------------------------------------
// CreateDatabase.sql -- Creates the database for the part numbering sys. 
//
// Created 2018-09-20 DLB
// --------------------------------------------------------------------
*/


drop Database M1Data;

Create Database M1Data;
Use M1Data;

create table Accounts(
  Aid int,
  ShortName char(10),
  DName varchar(120),
  FName varchar(120),
  Notes varchar(1200),
  Active int
);

insert into Accounts(Aid, Active, ShortName, DName, FName, Notes) values(1, 0, "ML",    "ML Checking",   "Merrill Lynch Checking",     "");
insert into Accounts(Aid, Active, ShortName, DName, FName, Notes) values(2, 1, "FMB",   "FM Checking",   "Farmers and Merchants Bank", "");
insert into Accounts(Aid, Active, ShortName, DName, FName, Notes) values(3, 1, "BofA",  "BofA Checking", "Bank of America Checking",   "");
insert into Accounts(Aid, Active, ShortName, DName, FName, Notes) values(4, 1, "DVisa", "Dal's Visa",    "Dal's Visa at ML",           "");
insert into Accounts(Aid, Active, ShortName, DName, FName, Notes) values(5, 1, "CVisa", "Carol's Visa",  "Carols Visa at ML",          ""); 
insert into Accounts(Aid, Active, ShortName, DName, FName, Notes) values(6, 1, "HDept", "Home Dept",     "Home Dept",                  "");

create table AccountAlias(
  Aid int,
  Alias varchar(120),
  Notes varchar(1200)
);

insert into AccountAlias(Aid, Alias, Notes) values(1, "ML Check",   "");
insert into AccountAlias(Aid, Alias, Notes) values(2, "FRB Check",  "");
insert into AccountAlias(Aid, Alias, Notes) values(3, "BofA Check", "");
insert into AccountAlias(Aid, Alias, Notes) values(4, "Visa 3232",  ""); 
insert into AccountAlias(Aid, Alias, Notes) values(4, "Visa 3627",  ""); 
insert into AccountAlias(Aid, Alias, Notes) values(4, "Visa 6974",  ""); 
insert into AccountAlias(Aid, Alias, Notes) values(4, "Visa 7759",  ""); 
insert into AccountAlias(Aid, Alias, Notes) values(4, "Visa 9389",  ""); 
insert into AccountAlias(Aid, Alias, Notes) values(5, "Visa 1638",  ""); 
insert into AccountAlias(Aid, Alias, Notes) values(5, "Visa 2859",  ""); 
insert into AccountAlias(Aid, Alias, Notes) values(5, "Visa 3799",  ""); 
insert into AccountAlias(Aid, Alias, Notes) values(5, "Visa 4007",  ""); 
insert into AccountAlias(Aid, Alias, Notes) values(5, "Visa 4907",  ""); 
insert into AccountAlias(Aid, Alias, Notes) values(5, "Visa 9713",  "");
insert into AccountAlias(Aid, Alias, Notes) values(5, "Visa 9836",  "");

create table Categories(
  Cid char(32),
  Name varchar(120),
  Notes varchar(1200)
);

create table CatAlias(
  Cid char(32),
  Alias varchar(120),
  Location varchar(120),
  Notes varchar(1200)
);

create table Vendors(
  Vid  char(32),
  DName varchar(120),
  FName varchar(120),
  DefaultCat varchar(120),
  BusinessType varchar(120),
  PrimaryProduct varchar(120),
  Notes varchar(1200)
);

create table VendorAlias(
  Vid char(32),
  Alias varchar(120)
);

create table Transactions(
  Tid char(32),
  Amount int,                  /* In cents, NOT dollars */
  Description varchar(240),
  DatePosted date,
  DateSettled date,
  Account char(32),            /* Required */
  Vendor char(32),             /* 0  = unknown */
  MonthNum int,                /* Statement Month, 1-12 */
  BankInfo varchar(240),
  Location varchar(240),
  CheckNum varchar(32)
);

create table CatList(
  Tid char(32),
  Cid char(32), 
  Amount int
);

