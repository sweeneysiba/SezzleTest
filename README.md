# SezzleTest
#can import the postman collection
    `https://www.getpostman.com/collections/03d52b301da316efbd14`



Get Request
    http://127.0.0.1:8080/order/list
    http://127.0.0.1:8080/cart/1602752657037971475/complete     token required
    http://127.0.0.1:8080/item/list
    http://127.0.0.1:8080/user/list
    http://127.0.0.1:8080/cart/list
Post Request 
    
    http://127.0.0.1:8080/user/create
        {
            "name":"siba sankar nayak",
            "username":"sweetsiba",
            "password":"Sibasankar@1"
        }

    http://127.0.0.1:8080/user/login
        {
            "username":"sweetsiba",
            "password":"Sibasankar@1"
        }
    http://127.0.0.1:8080/item/create 
        {
            "name":"item 3"
        }
    http://127.0.0.1:8080/cart/add
    [
        {
            "id": 1602752670801156439,
            "name": "item 3"
        }
    ]