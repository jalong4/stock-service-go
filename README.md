# Stock Service
Powering the StockViewer App

## About This Service
This application provides a comprehensive set of tools for managing your financial portfolios, tracking stock performance, and more.

## API Routes


* **POST** `/users/login` - Authenticate user and provide tokens

* **POST** `/users/register` - Register a new user

* **GET** `/users/` - Retrieve all users (Requires Auth)

* **GET** `/users/id/:_id` - Retrieve a user by their ID (Requires Auth)

* **DELETE** `/users/id/:_id` - Delete a user by their ID (Requires Auth)

* **PUT** `/users/id/:_id` - Update a user by their ID (Requires Auth)

* **GET** `/holdings/` - Retrieve all holdings (Requires Auth)

* **POST** `/holdings/` - Add a new holding (Requires Auth)

* **GET** `/holdings/id/:_id` - Retrieve a holding by its ID (Requires Auth)

* **DELETE** `/holdings/id/:_id` - Delete a holding by its ID (Requires Auth)

* **PUT** `/holdings/id/:_id` - Update a holding by its ID (Requires Auth)

* **GET** `/holdings/ticker/:ticker` - Retrieve holdings by ticker (Requires Auth)

* **GET** `/holdings/account/:account` - Retrieve holdings by account (Requires Auth)


<br><br>
© 2024 Long Software Inc. All rights reserved.
