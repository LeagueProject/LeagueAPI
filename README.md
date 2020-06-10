
# LeagueAPI
# Build insturctions:
`git clone https://github.com/LeagueProject/LeagueAPI` \
`cd LeagueAPI` \
`go get github.com/julienschmidt/httprouter` \
`go build` \
`./LeagueAPI`  

# Api usage :
`get 34.67.7.77:8555/read/user?id=1005  -> returneaza un user dupa ID sub forma json` 
```
post 34.67.7.77:8555/add/user -> in body-ul din request se pune un fisier de tip json care descrie structura user
Ex : in body
 {
    "UID": 1005,
    "IEmail": "mihai.indreias@poli.ro",
    "PMail": "mihai.indreias@gmail.com",
    "Username": "Mihai Indreias",
    "Password" : "e10adc3949ba59abbe56e057f20f883e"  
    "YearOfStudy": 1,
    "College": "Politehnica Bucuresti",
    "University": "Politehnica",
    "Major": "Info",
    "Serie": "123"
}

->returneaza un json cu un userID si un cod

```

```
get 34.67.7.77:8555/login?user=cnmsr&pass=e10adc3949ba59abbe56e057f20f883e
->returneaza un json cu un sessionID si un cod


```
