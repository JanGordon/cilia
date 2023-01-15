const devSocket = new WebSocket("ws:localhost:8080/ws");

devSocket.onopen = (event) => {
    console.log("connection opened")

}

devSocket.onmessage = async (event) => {
    if (event.data == "reload") {
        devSocket.send("reload successful");
        window.location.reload();
    } else if (event.data == "reloadhtml") {
        devSocket.send("reload successful");
        fetch(window.location.pathname)
        .then(function (response) {
            return response.text()
        }).then(function (data) {
            document.body.innerHTML = data
            console.log("reloaded html", data)
        })
    }
}