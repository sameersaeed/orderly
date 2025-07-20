# Orderly

A full stack e-commerce app created using React + Golang. Customers can create and submit orders through the frontend, which are then sent to a private Discord channel where they can be processed by site administrators.


# Setup steps

To get the client up and running, you can navigate to `~/client` (where `~` is the project's root directory) and run: `npm start`

To run the server, navigate to `~/server` in a separate terminal window and run: `go run main.go`


# Overview

Users can register an account and submit orders through the client

<img width="785" alt="Screenshot 2024-06-21 at 7 25 39 PM" src="https://github.com/sameersaeed/order_submission_tool/assets/70821949/f9615fc5-a292-4e0c-9d8f-0d216d8678b6">

The the server will then push the order to a Discord channel for the site administrators to view
The <img width="714" alt="Screenshot 2024-06-21 at 7 25 51 PM" src="https://github.com/sameersaeed/order_submission_tool/assets/70821949/c9f56c93-8c7b-4c49-a2f1-5b2f6fbe1dce">

The user can also then view their past orders through the client as well, and submit a ticket with the site administrators using their order ID for any issues they may have with their order
<img width="833" alt="Screenshot 2024-06-21 at 7 50 12 PM" src="https://github.com/sameersaeed/order_submission_tool/assets/70821949/a38a48d7-914b-405c-a81d-b290308064f2">

Site administrators also have the option to create, delete, and edit existing items within the store


### **A Site administrator's view of the shop page:**
<img width="769" alt="Screenshot 2024-06-21 at 8 03 02 PM" src="https://github.com/sameersaeed/order_submission_tool/assets/70821949/2010e52c-ee27-43a7-b301-6c6d3efa99ea">



### **Adding a new item:**
<img width="739" alt="Screenshot 2024-06-21 at 8 03 26 PM" src="https://github.com/sameersaeed/order_submission_tool/assets/70821949/d9c738c7-1a93-4199-8852-db9f493231c7">


### **Editing an existing item:**
<img width="764" alt="Screenshot 2024-06-21 at 8 03 52 PM" src="https://github.com/sameersaeed/order_submission_tool/assets/70821949/dc618ada-055a-4597-8326-1d5ae4489c01">
