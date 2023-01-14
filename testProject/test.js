const exampleSocket = new WebSocket("ws:localhost:8080/ws");

exampleSocket.onopen = (event) => {
    console.log("connection opened")

}

exampleSocket.onmessage = (event) => {
    if (event.data == "reload") {
        exampleSocket.send("reload successful");
        window.location.reload();
    }
}