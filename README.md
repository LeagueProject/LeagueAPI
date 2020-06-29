
# Build instructions:
`git clone https://github.com/LeagueProject/LeagueAPI` \
`cd LeagueAPI` \
`go get github.com/julienschmidt/httprouter` \
`go get github.com/lib/pq` \
`go build` \
`./LeagueAPI`  

# Api usage :

# Read
```
get $host:$port/read/user?id=3060246767619880621  -> returneaza un user dupa ID sub forma json
    EX :
    {
    "UID": 3060246767619880621,
    "IEmail": "mihai.indreias@gmail.com",
    "PMail": "mihai.indreias@gmail.com",
    "Username": "g0g05arui",
    "Password": "",
    "YearOfStudy": 2,
    "College": "Politehnica",
    "University": "Poli Buc",
    "Major": "CS",
    "Serie": "S123",
    "FirstName": "Mihai",
    "LastName": "Indreias",
    "Following": [
        3060246767619880621,
        5090235015690038407
    ],
    "Followers": [
        3060246767619880621
    ]
}
```
```
get $host:$port/read/message?sid=6021198439185793553&uid=10517390909092599&mid=2023716118525109935 
-> returneaza un mesaj dupa uID(user) & sID(session) & mID(message) sub forma json
Ex:
{
    "ID":123345678,
    "AuthorID":10517390909092599,
    "Text":"Test",
    "Date":"06.15.2020:14:58.",
    "Receiver":10517390909092599,
    "TypeOfReceiver":"person"
}
```
# Add
```
post $host:$port/add/user
    + In body-ul requestului :
    {
    "IEmail":"mihai.indreias@gmail.com",
    "PMail":"mihai.indreias@gmail.com",
    "Username":"g0g05arui",
    "Password":"Test",
    "YearOfStudy":2,
    "College":"Politehnica",
    "University":"Poli Buc",
    "Serie":"S123",
    "Major":"CS",
    "FirstName":"Mihai",
    "LastName":"Indreias"
}
Returneaza un HTTPResonse ( vezi data_types.go ) cu raspunsul 
```
```
post $host:$port/add/message?id=123
    -> id = id-ul sessiunii
    + In body-ul requestului:
    {
    "ID":123345678,
    "AuthorID":10517390909092599,
    "Text":"Test",
    "Date":"06.15.2020:14:58.",
    "Receiver":10517390909092599,
    "TypeOfReceiver":"person"
    }
Returneaza un HTTPResonse ( vezi data_types.go ) cu raspunsul 
```
# Activate
```
get $host:$port/activate?id=6030937711187684612
->Activeaza userul cu id-ul $id
Returneaza un HTTPResonse ( vezi data_types.go ) cu raspunsul 
```
```
post $host:$port/login
+ In body-ul requestului:

{
    "Username":"g0g05arui",
    "Password":"test"
}
Returneaza un HTTPResonse ( vezi data_types.go ) cu raspunsul 
    ->{"Response":["0","Not verified"],"Code":404}
    ->{"Response":["3060246767619880621","5010926770203421665"],"Code":200}   ( Response[0]=uID & Response[1]=sID)
    ->{"Response":["0"],"Code":400}
```

# Check session

```
post $host:$port:/check/session?sid=12345&uid=10005
Returneaza un HTTPResonse ( vezi data_types.go ) cu raspunsul 
```
# FollowStatus

```
post $host:$port:followStatus/follow?from=3060246767619880621&to=5090235015690038407&sid=7003377197631441573
->Da follow userului $TO de la $FROM care Session Id-ul $SID
Returneaza un HTTPResonse ( vezi data_types.go ) cu raspunsul 
```
```
post $host:$port:followStatus/unfollow?from=3060246767619880621&to=5090235015690038407&sid=7003377197631441573
->Da unfollow userului $TO de la $FROM care Session Id-ul $SID
Returneaza un HTTPResonse ( vezi data_types.go ) cu raspunsul 
```