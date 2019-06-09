# OSS
Online Shopping System - HUST Database Course Project//



## Database DESIGN

Operator Roles: `Administrator`, `Merchant` & `Customer`.

Users Table:

`UserID` int auto-inc primary-key

`username` string

`password` string

`role` string: `admin`/`merchant`/`customer`.



User Information Table:

`UInfoID` int auto-inc primary-key

`UID` int foreign-key(User Table)

`Name` string

`Phone` string

`Address` string



Shop Information Table:

`ShopID` int auto-inc primary-key

`BelongsTo` int foreign-key(User Table)

`ShopName` string



Stock Information Table:

`StockId` int auto-inc primary-key

`SoldAt` int foreign-key(ShopInformationTable)

`Name` string

`Type` string

`Desc` string

`Price` float

`Picture` string, base64 encoded picture

`Promotion` float

`StockCount` int

`Enabled` boolean



Order Information Table

`OrderId` int auto-inc primary-key

`FromCustomer` int foreign-key(UserInformationTable)

`OrderedStock` int foreign-key(StockInformationTable)

`Amount` int

`InstantPrice` float

`PriceSum` float

`PaymentDetail` string

`DeliveryDetail` string

`Status` string: `issued`, `paid`, `delivered`, `finished`.



```mysql
create table users (
	uid int not null auto_increment,
    username text not null,
    password text not null,
    role text not null,
    primary key (uid)
);

create table user_info (
	uinfoid int not null auto_increment,
    uid int not null,
    name text not null,
    phone text not null,
    address text not null,
    primary key (uinfoid),
    foreign key (uid) references users (uid)
);

create table shop_info (
	shopid int not null auto_increment,
    belongs int not null,
    name text not null,
    primary key (shopid),
    foreign key (belongs) references users (uid)
);

create table stock_info (
	stockid int not null auto_increment,
    sold_at int not null,
    name text not null,
    type text not null,
    description text not null,
    price double not null,
    picture text not null,
    promotion double not null,
    stock_count int not null,
    enabled boolean not null,
    primary key (stockid),
    foreign key (sold_at) references shop_info (shopid)
);

create table order_info (
	orderid int not null auto_increment,
    cust int not null,
    stock int not null,
    amount int not null,
    inst_price double not null,
    price_sum double not null,
    payment_detail text not null,
    delivery_detail text not null,
    stat text not null,
    primary key (orderid),
    foreign key (cust) references users (uid),
    foreign key (stock) references stock_info (stockid)
);
```

